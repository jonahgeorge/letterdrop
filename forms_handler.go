package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/flosch/pongo2"
	"github.com/gorilla/mux"
	"github.com/jonahgeorge/letterdrop/models"
	repo "github.com/jonahgeorge/letterdrop/repositories"
)

func (app *Application) FormsIndexHandler(w http.ResponseWriter, r *http.Request, currentUser *models.User) {
	forms, _ := repo.NewFormsRepository(app.db).FindByUserId(currentUser.Id)

	app.Render(w, r, "forms/index", pongo2.Context{
		"forms": forms,
	})
}

func (app *Application) FormsNewHandler(w http.ResponseWriter, r *http.Request, currentUser *models.User) {
	app.Render(w, r, "forms/new", pongo2.Context{})
}

func (app *Application) FormsCreateHandler(w http.ResponseWriter, r *http.Request, currentUser *models.User) {
	session, _ := app.GetSession(r)

	recaptchaSecretKey := r.PostFormValue("recaptcha_secret_key")
	description := r.PostFormValue("description")

	form := new(models.Form)
	form.Name = r.PostFormValue("name")
	if len(description) > 0 {
		form.Description = &description
	} else {
		form.Description = nil
	}
	if len(recaptchaSecretKey) > 0 {
		form.RecaptchaSecretKey = &recaptchaSecretKey
	} else {
		form.RecaptchaSecretKey = nil
	}

	_, err := repo.NewFormsRepository(app.db).Create(currentUser.Id, form.Name, form.Description, form.RecaptchaSecretKey)
	if err != nil {
		session.AddFlash("An error occured while creating your form")
		session.Save(r, w)
		app.Render(w, r, "forms/new", pongo2.Context{"form": form})
		return
	}

	session.AddFlash("Successfully created form!")
	session.Save(r, w)

	http.Redirect(w, r, "/forms", 302)
}

func (app *Application) FormsShowHandler(w http.ResponseWriter, r *http.Request, currentUser *models.User) {
	session, _ := app.GetSession(r)
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	form, _ := repo.NewFormsRepository(app.db).FindById(id)
	submissions, _ := repo.NewSubmissionsRepository(app.db).FindByFormId(form.Id)
	if !currentUser.CanView(form) {
		session.AddFlash("You are not authorized to access this resource.")
		session.Save(r, w)
		http.Redirect(w, r, "/forms", 302)
		return
	}

	app.Render(w, r, "forms/show", pongo2.Context{
		"form":        form,
		"submissions": submissions,
	})
}

func (app *Application) FormsEditHandler(w http.ResponseWriter, r *http.Request, currentUser *models.User) {
	session, _ := app.GetSession(r)
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	form, _ := repo.NewFormsRepository(app.db).FindById(id)
	if !currentUser.CanUpdate(form) {
		session.AddFlash("You are not authorized to access this resource.")
		session.Save(r, w)
		http.Redirect(w, r, "/forms", 302)
		return
	}

	app.Render(w, r, "forms/edit", pongo2.Context{"form": form})
}

func (app *Application) FormsUpdateHandler(w http.ResponseWriter, r *http.Request, currentUser *models.User) {
	session, _ := app.GetSession(r)
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	form, err := repo.NewFormsRepository(app.db).FindById(id)
	if !currentUser.CanUpdate(form) {
		session.AddFlash("You are not authorized to access this resource.")
		session.Save(r, w)
		http.Redirect(w, r, "/forms", 302)
		return
	}

	form.Name = r.PostFormValue("name")
	description := r.PostFormValue("description")
	if len(description) > 0 {
		form.Description = &description
	} else {
		form.Description = nil
	}

	recaptchaSecretKey := r.PostFormValue("recaptcha_secret_key")
	if len(recaptchaSecretKey) > 0 {
		form.RecaptchaSecretKey = &recaptchaSecretKey
	} else {
		form.RecaptchaSecretKey = nil
	}

	_, err = repo.NewFormsRepository(app.db).Update(id, form.Name, form.Description, form.RecaptchaSecretKey)
	if err != nil {
		session.AddFlash("An error occured while updating this form")
		session.Save(r, w)
		app.Render(w, r, "forms/edit", pongo2.Context{"form": form})
		return
	}

	session.AddFlash("Successfully updated form!")
	session.Save(r, w)

	http.Redirect(w, r, fmt.Sprintf("/forms/%d", form.Id), 302)
}

func (app *Application) FormsDestroyHandler(w http.ResponseWriter, r *http.Request, currentUser *models.User) {
	session, _ := app.GetSession(r)
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	form, _ := repo.NewFormsRepository(app.db).FindById(id)
	if !currentUser.CanDelete(form) {
		session.AddFlash("You are not authorized to access this resource.")
		session.Save(r, w)
		http.Redirect(w, r, "/forms", 302)
		return
	}

	_, err := repo.NewFormsRepository(app.db).Delete(form.Id)
	if err != nil {
		session.AddFlash("An error occured while deleting this form")
		session.Save(r, w)
		app.Render(w, r, "forms/edit", pongo2.Context{"form": form})
		return
	}

	session.AddFlash("Successfully deleted form!")
	session.Save(r, w)

	http.Redirect(w, r, "/forms", 302)
}
