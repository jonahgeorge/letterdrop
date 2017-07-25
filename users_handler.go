package main

import (
	"github.com/flosch/pongo2"
	"net/http"
)

func (app *Application) UsersNewHandler(w http.ResponseWriter, r *http.Request) {
	app.Render(w, r, "users/new", pongo2.Context{})
}

func (app *Application) UsersCreateHandler(w http.ResponseWriter, r *http.Request) {
	_, err := NewUsersRepository(app.db).Create(
		r.PostFormValue("email"), r.PostFormValue("password"))
	if err != nil {
		session.AddFlash("Woah, something bad happened.")
		app.Render(w, r, "users/new", pongo2.Context{
			"email": r.PostFormValue("email"),
		})
	}

	session.AddFlash("Successfully signed up!")
	http.Redirect(w, r, "/", 302)
}
