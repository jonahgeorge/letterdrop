package main

import (
	"net/http"

	"github.com/flosch/pongo2"
)

func (app *Application) IndexHandler(w http.ResponseWriter, r *http.Request) {
	app.Render(w, r, "index", pongo2.Context{})
}
