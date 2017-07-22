package main

import (
  "net/http"
  "github.com/flosch/pongo2"
  "fmt"
)

func (app *Application) UsersNewHandler(w http.ResponseWriter, r *http.Request) {
  app.Render(w, r, "users/new", pongo2.Context{})
}

func (app *Application) UsersCreateHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "Users#create")
}

