package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/flosch/pongo2"
	"github.com/tuvistavie/securerandom"
)

func (app *Application) UsersNewHandler(w http.ResponseWriter, r *http.Request) {
	app.Render(w, r, "users/new", pongo2.Context{})
}

func (app *Application) UsersCreateHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := app.GetSession(r)

	token, _ := securerandom.UrlSafeBase64(10, true)
	newUser := &User{
		name:                   r.PostFormValue("name"),
		email:                  r.PostFormValue("email"),
		passwordDigest:         r.PostFormValue("password"),
		emailConfirmationToken: &token,
	}

	if !app.recaptchaClient.Verify(*r) {
		session.AddFlash("Invalid ReCaptcha")
		session.Save(r, w)
		app.Render(w, r, "users/new", pongo2.Context{"user": newUser})
		return
	}

	user, err := NewUsersRepository(app.db).Create(newUser)
	if err != nil {
		log.Println(err)
		session.AddFlash("Woah, something bad happened.")
		session.Save(r, w)
		app.Render(w, r, "users/new", pongo2.Context{"user": newUser})
		return
	}

	app.SendEmailConfirmation(user)

	session.Values["userId"] = user.id
	session.AddFlash(fmt.Sprintf("Welcome, %s!", user.name))
	session.Save(r, w)
	http.Redirect(w, r, "/", 302)
}
