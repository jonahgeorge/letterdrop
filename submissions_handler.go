package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/flosch/pongo2"
	"github.com/gorilla/mux"
	"github.com/haisum/recaptcha"
	"github.com/jonahgeorge/letterdrop/mailers"
	"github.com/jonahgeorge/letterdrop/models"
	repo "github.com/jonahgeorge/letterdrop/repositories"
)

func (app *Application) SubmissionsCreateHandler(w http.ResponseWriter, r *http.Request) {
	uuid := mux.Vars(r)["uuid"]
	r.ParseForm()

	form, err := repo.NewFormsRepository(app.db).FindByUuid(uuid)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Process Recaptcha if enabled
	if form.RecaptchaSecretKey != nil {
		recaptchaClient := recaptcha.R{
			Secret: *form.RecaptchaSecretKey,
		}

		ok := recaptchaClient.Verify(*r)
		if !ok {
			http.Error(w, "Recaptcha Verification failed", 403)
			return
		}
	}

	r.Form.Del("g-recaptcha-response")
	json, err := json.MarshalIndent(r.Form, "", " ")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	_, err = repo.NewSubmissionsRepository(app.db).Create(form.Id, string(json))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	user, _ := repo.NewUsersRepository(app.db).FindById(form.UserId)

	mailers.SendSubmissionNotification(app.emailClient, user, form, json)

	http.Redirect(w, r, r.Referer(), 302)
}

func (app *Application) SubmissionsDestroyHandler(w http.ResponseWriter, r *http.Request, currentUser *models.User) {
	session, _ := app.GetSession(r)

	formId, _ := strconv.Atoi(mux.Vars(r)["formId"])
	submissionId, _ := strconv.Atoi(mux.Vars(r)["submissionId"])

	form, _ := repo.NewFormsRepository(app.db).FindById(formId)
	submissions, _ := repo.NewSubmissionsRepository(app.db).FindByFormId(form.Id)
	if !currentUser.CanDelete(form) {
		session.AddFlash("You are not authorized to access this resource.")
		session.Save(r, w)
		http.Redirect(w, r, fmt.Sprintf("/forms/%d", form.Id), 302)
		return
	}

	_, err := repo.NewSubmissionsRepository(app.db).Delete(submissionId)
	if err != nil {
		session.AddFlash("An error occured while deleting this submission")
		session.Save(r, w)
		app.Render(w, r, "forms/show", pongo2.Context{
			"form":        form,
			"submissions": submissions,
		})
		return
	}

	session.AddFlash("Successfully deleted submission!")
	session.Save(r, w)

	http.Redirect(w, r, fmt.Sprintf("/forms/%d", form.Id), 302)
}
