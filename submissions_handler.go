package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/flosch/pongo2"
	"github.com/gorilla/mux"
	"github.com/haisum/recaptcha"
)

func (app *Application) SubmissionsCreateHandler(w http.ResponseWriter, r *http.Request) {
	uuid := mux.Vars(r)["uuid"]
	r.ParseForm()

	form, err := NewFormsRepository(app.db).FindByUuid(uuid)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Process Recaptcha if enabled
	if form.recaptchaSecretKey != nil {
		recaptchaClient := recaptcha.R{
			Secret: *form.recaptchaSecretKey,
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

	_, err = NewSubmissionsRepository(app.db).Create(form.id, string(json))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	app.SendSubmissionNotification(form, json)
	http.Redirect(w, r, r.Referer(), 302)
}

func (app *Application) SubmissionsDestroyHandler(w http.ResponseWriter, r *http.Request, currentUser *User) {
	session, _ := app.GetSession(r)

	formId, _ := strconv.Atoi(mux.Vars(r)["formId"])
	submissionId, _ := strconv.Atoi(mux.Vars(r)["submissionId"])

	form, _ := NewFormsRepository(app.db).FindById(formId)
	submissions, _ := NewSubmissionsRepository(app.db).FindByFormId(form.id)

	// TODO Check authorization on form

	_, err := NewSubmissionsRepository(app.db).Delete(submissionId)
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
	http.Redirect(w, r, fmt.Sprintf("/forms/%d", form.id), 302)
}
