package main

import (
	"fmt"
	"net/http"

	"github.com/flosch/pongo2"
)

func (app *Application) UsersNewHandler(w http.ResponseWriter, r *http.Request) {
	app.Render(w, r, "users/new", pongo2.Context{})
}

func (app *Application) UsersCreateHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := app.GetSession(r)

	newUser := &User{
		name:           r.PostFormValue("name"),
		email:          r.PostFormValue("email"),
		passwordDigest: r.PostFormValue("password"),
	}

	user, err := NewUsersRepository(app.db).Create(newUser.name, newUser.email, newUser.passwordDigest)
	if err != nil {
		session.AddFlash("Woah, something bad happened.")
		session.Save(r, w)
		app.Render(w, r, "users/new", pongo2.Context{"user": newUser})
		return
	}

	session.Values["userId"] = user.id
	session.AddFlash(fmt.Sprintf("Welcome, %s!", user.name))
	session.Save(r, w)
	http.Redirect(w, r, "/", 302)
}
