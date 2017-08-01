package mailers

import (
	"github.com/flosch/pongo2"
	"github.com/jonahgeorge/letterdrop/models"
	"github.com/sendgrid/rest"
	sendgrid "github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

const (
	EMAIL_CONFIRMATION_TEXT_TEMPLATE = `TODO Plaintext Emails`

	EMAIL_CONFIRMATION_HTML_TEMPLATE = `
<p>Click and confirm that you want to create an account on LetterDrop. This link can only be used once.</p>
<p><a href="{{ host }}/email_confirmation?token={{ user.EmailConfirmationToken }}">Verify Email</a></p>
<p>Or</p>
<p>{{ host }}/email_confirmation?token={{ user.EmailConfirmationToken }}</p>`
)

func SendEmailConfirmation(emailClient *sendgrid.Client, hostName string, user *models.User) (*rest.Response, error) {
	c := pongo2.Context{
		"user": user,
		"host": hostName,
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

	return emailClient.Send(message)
}
