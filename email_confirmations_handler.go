package main

import (
	"log"
	"net/http"

	"github.com/flosch/pongo2"
	"github.com/jonahgeorge/letterdrop/mailers"
	repo "github.com/jonahgeorge/letterdrop/repositories"
	"github.com/tuvistavie/securerandom"
)

func (app *Application) EmailConfirmationsNewHandler(w http.ResponseWriter, r *http.Request) {
	app.Render(w, r, "email_confirmations/new", pongo2.Context{})
}

func (app *Application) EmailConfirmationsCreateHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := app.GetSession(r)

	usersRepo := repo.NewUsersRepository(app.db)

	ok := app.recaptchaClient.Verify(*r)
	if !ok {
		session.AddFlash("Invalid ReCaptcha")
		session.Save(r, w)
		app.Render(w, r, "email_confirmations/new", pongo2.Context{"email": r.PostFormValue("email")})
		return
	}

	user, _ := usersRepo.FindByEmail(r.PostFormValue("email"))

	token, _ := securerandom.UrlSafeBase64(10, true)
	user.EmailConfirmationToken = &token
	_, err := usersRepo.Update(user)
	if err != nil {
		log.Println(err)
		session.AddFlash("Something bad happened :(")
		session.Save(r, w)
		app.Render(w, r, "email_confirmations/new", pongo2.Context{"email": r.PostFormValue("email")})
		return
	}

	mailers.SendEmailConfirmation(app.emailClient, app.hostName, user)
	session.AddFlash("Please check your email for confirmation instructions.")
	session.Save(r, w)
	http.Redirect(w, r, "/email_confirmation/new", 302)

}

func (app *Application) EmailConfirmationsShowHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := app.GetSession(r)

	token := r.URL.Query().Get("token")

	usersRepo := repo.NewUsersRepository(app.db)

	user, _ := usersRepo.FindByEmailConfirmationToken(token)
	if user == nil {
		session.AddFlash("This token is not associated with a user.")
		session.Save(r, w)
		http.Redirect(w, r, "/email_confirmation/new", 302)
		return
	}

	user.EmailConfirmationToken = nil
	user.IsEmailConfirmed = true
	usersRepo.Update(user) // TODO Error handling

	session.Values["userId"] = user.Id
	session.AddFlash("Successfully confirmed your email address.")
	session.Save(r, w)
	http.Redirect(w, r, "/", 302)
}
