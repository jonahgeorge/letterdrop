package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/flosch/pongo2"
	"github.com/gorilla/mux"
)

func (app *Application) FormsIndexHandler(w http.ResponseWriter, r *http.Request, currentUser *User) {
	forms, _ := NewFormsRepository(app.db).FindByUserId(currentUser.id)

	app.Render(w, r, "forms/index", pongo2.Context{
		"forms": forms,
	})
}

func (app *Application) FormsNewHandler(w http.ResponseWriter, r *http.Request, currentUser *User) {
	app.Render(w, r, "forms/new", pongo2.Context{})
}

func (app *Application) FormsCreateHandler(w http.ResponseWriter, r *http.Request, currentUser *User) {
	session, _ := app.GetSession(r)

	recaptchaSecretKey := r.PostFormValue("recaptcha_secret_key")
	description := r.PostFormValue("description")

	form := new(Form)
	form.name = r.PostFormValue("name")
	if len(description) > 0 {
		form.description = &description
	} else {
		form.description = nil
	}
	if len(recaptchaSecretKey) > 0 {
		form.recaptchaSecretKey = &recaptchaSecretKey
	} else {
		form.recaptchaSecretKey = nil
	}

	_, err := NewFormsRepository(app.db).Create(currentUser.id, form.name, form.description, form.recaptchaSecretKey)
	if err != nil {
		session.AddFlash("An error occured while creating your form")
		session.Save(r, w)
		app.Render(w, r, "forms/new", pongo2.Context{
			"form": form,
		})
		return
	}

	session.AddFlash("Successfully created form!")
	session.Save(r, w)
	http.Redirect(w, r, "/forms", 302)
}

func (app *Application) FormsShowHandler(w http.ResponseWriter, r *http.Request, currentUser *User) {
	session, _ := app.GetSession(r)
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	form, _ := NewFormsRepository(app.db).FindById(id)
	submissions, _ := NewSubmissionsRepository(app.db).FindByFormId(form.id)
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

func (app *Application) FormsEditHandler(w http.ResponseWriter, r *http.Request, currentUser *User) {
	session, _ := app.GetSession(r)
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	form, _ := NewFormsRepository(app.db).FindById(id)
	if !currentUser.CanUpdate(form) {
		session.AddFlash("You are not authorized to access this resource.")
		session.Save(r, w)
		http.Redirect(w, r, "/forms", 302)
		return
	}

	app.Render(w, r, "forms/edit", pongo2.Context{
		"form": form,
	})
}

func (app *Application) FormsUpdateHandler(w http.ResponseWriter, r *http.Request, currentUser *User) {
	session, _ := app.GetSession(r)
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	form, err := NewFormsRepository(app.db).FindById(id)
	if !currentUser.CanUpdate(form) {
		session.AddFlash("You are not authorized to access this resource.")
		session.Save(r, w)
		http.Redirect(w, r, "/forms", 302)
		return
	}

	form.name = r.PostFormValue("name")
	description := r.PostFormValue("description")
	if len(description) > 0 {
		form.description = &description
	} else {
		form.description = nil
	}

	recaptchaSecretKey := r.PostFormValue("recaptcha_secret_key")
	if len(recaptchaSecretKey) > 0 {
		form.recaptchaSecretKey = &recaptchaSecretKey
	} else {
		form.recaptchaSecretKey = nil
	}

	_, err = NewFormsRepository(app.db).Update(id, form.name, form.description, form.recaptchaSecretKey)
	if err != nil {
		session.AddFlash("An error occured while updating this form")
		session.Save(r, w)
		app.Render(w, r, "forms/edit", pongo2.Context{
			"form": form,
		})
		return
	}

	session.AddFlash("Successfully updated form!")
	session.Save(r, w)
	http.Redirect(w, r, fmt.Sprintf("/forms/%d", form.id), 302)
}

func (app *Application) FormsDestroyHandler(w http.ResponseWriter, r *http.Request, currentUser *User) {
	session, _ := app.GetSession(r)
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	form, _ := NewFormsRepository(app.db).FindById(id)
	if !currentUser.CanDelete(form) {
		session.AddFlash("You are not authorized to access this resource.")
		session.Save(r, w)
		http.Redirect(w, r, "/forms", 302)
		return
	}

	_, err := NewFormsRepository(app.db).Delete(form.id)
	if err != nil {
		session.AddFlash("An error occured while deleting this form")
		session.Save(r, w)
		app.Render(w, r, "forms/edit", pongo2.Context{
			"form": form,
		})
		return
	}

	session.AddFlash("Successfully deleted form!")
	session.Save(r, w)
	http.Redirect(w, r, "/forms", 302)
}
