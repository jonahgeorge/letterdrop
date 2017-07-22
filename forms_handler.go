package main

import (
  "net/http"
  "github.com/flosch/pongo2"
)

func (app *Application) FormsIndexHandler(w http.ResponseWriter, r *http.Request) {
  session, _ := app.sessions.Get(r, "letterbox")
  id := session.Values["userId"]

  forms := NewFormsRepository(app.db).FindByUserId(id.(int))

  app.Render(w, r, "forms/index", pongo2.Context{
    "forms": forms,
  })
}

func (app *Application) FormsNewHandler(w http.ResponseWriter, r *http.Request) {
  app.Render(w, r, "forms/new", pongo2.Context{})
}

func (app *Application) FormsCreateHandler(w http.ResponseWriter, r *http.Request) {
  session, _ := app.sessions.Get(r, "letterbox")
  id := session.Values["userId"]

  _, err := NewFormsRepository(app.db).Create(id.(int), r.PostFormValue("name"))
  if err != nil {
    // Add flash
    app.Render(w, r, "forms/new", pongo2.Context{
      "name": r.PostFormValue("name"),
    })
  } else {
    // Add flash
    http.Redirect(w, r, "/forms", 302)
  }
}

func (app *Application) FormsShowHandler(w http.ResponseWriter, r *http.Request) {
  app.Render(w, r, "forms/show", pongo2.Context{})
}

