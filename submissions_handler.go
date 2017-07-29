package main

import (
	"encoding/json"
	"fmt"
	"github.com/flosch/pongo2"
	"github.com/gorilla/mux"
	"github.com/haisum/recaptcha"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"net/http"
	"strconv"
)

const (
	plainTextTemplate = `
Plaintext New Form Submission

{{ json }}`

	htmlTemplate = `
<h1>New Form Submission</h1> 
<br>
<pre>{{ json }}</pre>`
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

	err = app.SendNotification(form, json)
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

func (app *Application) SendNotification(form *Form, json []byte) error {
	user, _ := NewUsersRepository(app.db).FindById(form.userId)

	c := pongo2.Context{
		"json": string(json),
	}

	htmlContent, _ := pongo2.FromString(htmlTemplate)
	html, _ := htmlContent.Execute(c)

	plainTextContent, _ := pongo2.FromString(plainTextTemplate)
	plainText, _ := plainTextContent.Execute(c)

	message := mail.NewSingleEmail(
		mail.NewEmail("Letterdrop Team", "team@letterdrop.herokuapp.com"),
		"New Form Submission",
		mail.NewEmail(user.name, user.email),
		plainText,
		html,
	)

	_, err := app.emailClient.Send(message)
	return err
}
