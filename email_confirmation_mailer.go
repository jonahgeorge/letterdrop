package main

import (
	"github.com/flosch/pongo2"
	"github.com/jonahgeorge/letterdrop/models"
	repo "github.com/jonahgeorge/letterdrop/repositories"
	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

const (
	EMAIL_CONFIRMATION_TEXT_TEMPLATE = `TODO Plaintext Emails`

	EMAIL_CONFIRMATION_HTML_TEMPLATE = `
<p>Click and confirm that you want to create an account on LetterDrop. This link can only be used once.</p>
<p><a href="{{ host }}/email_confirmation?token={{ user.EmailConfirmationToken }}">Verify Email</a></p>
<p>Or</p>
<p>{{ host }}/email_confirmation?token={{ user.EmailConfirmationToken }}</p>`

	SUBMISSION_NOTIFICATION_TEXT_TEMPLATE = `TODO Plaintext Emails`

	SUBMISSION_NOTIFICATION_HTML_TEMPLATE = `
<h1>New Form Submission</h1> 
<br>
<pre>{{ json }}</pre>`
)

func (app *Application) SendEmailConfirmation(user *models.User) (*rest.Response, error) {
	c := pongo2.Context{
		"user": user,
		"host": app.hostName,
	}

	htmlContent, _ := pongo2.FromString(EMAIL_CONFIRMATION_HTML_TEMPLATE)
	html, _ := htmlContent.Execute(c)

	plainTextContent, _ := pongo2.FromString(EMAIL_CONFIRMATION_TEXT_TEMPLATE)
	plainText, _ := plainTextContent.Execute(c)

	message := mail.NewSingleEmail(
		mail.NewEmail("LetterDrop", "team@letterdrop.herokuapp.com"),
		"Verify your email",
		mail.NewEmail(user.Name, user.Email),
		plainText,
		html,
	)

	return app.emailClient.Send(message)
}

func (app *Application) SendSubmissionNotification(form *models.Form, json []byte) (*rest.Response, error) {
	user, _ := repo.NewUsersRepository(app.db).FindById(form.UserId)

	c := pongo2.Context{"json": string(json)}

	htmlContent, _ := pongo2.FromString(SUBMISSION_NOTIFICATION_HTML_TEMPLATE)
	html, _ := htmlContent.Execute(c)

	plainTextContent, _ := pongo2.FromString(SUBMISSION_NOTIFICATION_TEXT_TEMPLATE)
	plainText, _ := plainTextContent.Execute(c)

	message := mail.NewSingleEmail(
		mail.NewEmail("LetterDrop", "team@letterdrop.herokuapp.com"),
		"New Form Submission",
		mail.NewEmail(user.Name, user.Email),
		plainText,
		html,
	)

	return app.emailClient.Send(message)
}
