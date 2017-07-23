package main

import (
	"github.com/flosch/pongo2"
	"net/http"
	// "log"
)

func (app *Application) UsersNewHandler(w http.ResponseWriter, r *http.Request) {
	app.Render(w, r, "users/new", pongo2.Context{})
}

func (app *Application) UsersCreateHandler(w http.ResponseWriter, r *http.Request) {
	usersRepo := NewUsersRepository(app.db)
	_, err := usersRepo.Create(r.PostFormValue("email"), r.PostFormValue("password"))
	if err == nil {
		// session.AddFlash("")
		http.Redirect(w, r, "/", 302)
	} else {
		// session.AddFlash("")
		app.Render(w, r, "users/new", pongo2.Context{
			"email": r.PostFormValue("email"),
		})
	}
}
