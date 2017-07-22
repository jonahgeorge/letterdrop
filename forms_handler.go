package main

import (
  "net/http"
  "github.com/flosch/pongo2"
  "fmt"
)

func (app *Application) FormsIndexHandler(w http.ResponseWriter, r *http.Request) {
  app.Render(w, r, "forms/index", pongo2.Context{})
}

func (app *Application) FormsNewHandler(w http.ResponseWriter, r *http.Request) {
  app.Render(w, r, "forms/new", pongo2.Context{})
}

func (app *Application) FormsCreateHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "Forms#create")
}

func (app *Application) FormsShowHandler(w http.ResponseWriter, r *http.Request) {
  app.Render(w, r, "forms/show", pongo2.Context{})
}

