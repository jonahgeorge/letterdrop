package main

import (
  "net/http"
  "github.com/flosch/pongo2"
  "fmt"
)

func (app *Application) SessionsNewHandler(w http.ResponseWriter, r *http.Request) {
  app.Render(w, r, "sessions/new", pongo2.Context{})
}

func (app *Application) SessionsCreateHandler(w http.ResponseWriter, r *http.Request) {
  session, _ := app.sessions.Get(r, "letterbox")
  // defer session.Save(r, w)

  userRepo := NewUsersRepository(app.db)
  user := userRepo.FindByEmailAndPassword(r.PostFormValue("email"), r.PostFormValue("password"))

  if user != nil {
    session.Values["userId"] = user.id
    session.AddFlash("Successfully logged in!")
    session.Save(r, w)
    http.Redirect(w, r, "/", 302)
  } else {
    session.AddFlash("Either your email or password was invalid.")
    app.Render(w, r, "sessions/new", pongo2.Context{
      "email": r.PostFormValue("email"),
    })
  }
}

func (app *Application) SessionsDestroyHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "Sessions#destroy")
}

