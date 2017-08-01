package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/flosch/pongo2"
	"github.com/jonahgeorge/letterdrop/mailers"
	"github.com/jonahgeorge/letterdrop/models"
	repo "github.com/jonahgeorge/letterdrop/repositories"
	"github.com/tuvistavie/securerandom"
)

func (app *Application) UsersNewHandler(w http.ResponseWriter, r *http.Request) {
	app.Render(w, r, "users/new", pongo2.Context{})
}

func (app *Application) UsersCreateHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := app.GetSession(r)

	token, _ := securerandom.UrlSafeBase64(10, true)
	newUser := &models.User{
		Name:                   r.PostFormValue("name"),
		Email:                  r.PostFormValue("email"),
		PasswordDigest:         r.PostFormValue("password"),
		EmailConfirmationToken: &token,
	}

	if !app.recaptchaClient.Verify(*r) {
		session.AddFlash("Invalid ReCaptcha")
		session.Save(r, w)
		app.Render(w, r, "users/new", pongo2.Context{"user": newUser})
		return
	}

	user, err := repo.NewUsersRepository(app.db).Create(newUser)
	if err != nil {
		log.Println(err)
		session.AddFlash("Woah, something bad happened.")
		session.Save(r, w)
		app.Render(w, r, "users/new", pongo2.Context{"user": newUser})
		return
	}

	mailers.SendEmailConfirmation(app.emailClient, app.hostName, user)

	session.Values["userId"] = user.Id
	session.AddFlash(fmt.Sprintf("Welcome, %s!", user.Name))
	session.Save(r, w)

	http.Redirect(w, r, "/", 302)
}
