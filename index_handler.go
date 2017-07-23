package main

import (
	"github.com/flosch/pongo2"
	"net/http"
)

func (app *Application) IndexHandler(w http.ResponseWriter, r *http.Request) {
	app.Render(w, r, "index", pongo2.Context{})
}
