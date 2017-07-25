package main

import (
	"github.com/flosch/pongo2"
	"net/http"
)

func (app *Application) SessionsNewHandler(w http.ResponseWriter, r *http.Request) {
	app.Render(w, r, "sessions/new", pongo2.Context{})
}

func (app *Application) SessionsCreateHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := app.GetSession(r)

	user, _ := NewUsersRepository(app.db).FindByEmailAndPassword(r.PostFormValue("email"), r.PostFormValue("password"))
	if user != nil {
		session.AddFlash("Either your email or password was invalid.")
		session.Save(r, w)
		app.Render(w, r, "sessions/new", pongo2.Context{
			"email": r.PostFormValue("email"),
		})
	}

	session.Values["userId"] = user.id
	session.AddFlash("Successfully logged in!")
	session.Save(r, w)
	http.Redirect(w, r, "/", 302)
}

func (app *Application) SessionsDestroyHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := app.GetSession(r)
	session.Values["userId"] = nil
	session.AddFlash("Successfully logged out!")
	session.Save(r, w)

	http.Redirect(w, r, "/", 302)
}
