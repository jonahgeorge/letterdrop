package main

import (
	"fmt"
	"github.com/flosch/pongo2"
	"net/http"
)

func (app *Application) UsersNewHandler(w http.ResponseWriter, r *http.Request) {
	app.Render(w, r, "users/new", pongo2.Context{})
}

func (app *Application) UsersCreateHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := app.GetSession(r)

	user, err := NewUsersRepository(app.db).Create(
		r.PostFormValue("name"),
		r.PostFormValue("email"),
		r.PostFormValue("password"))

	if err != nil {
		session.AddFlash("Woah, something bad happened.")
		session.Save(r, w)
		app.Render(w, r, "users/new", pongo2.Context{
			"name":  r.PostFormValue("name"),
			"email": r.PostFormValue("email"),
		})
	}

	session.Values["userId"] = user.id
	session.AddFlash(fmt.Sprintf("Welcome, %s!", user.name))
	session.Save(r, w)
	http.Redirect(w, r, "/", 302)
}
