package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func (app *Application) SubmissionsCreateHandler(w http.ResponseWriter, r *http.Request) {
	formId, _ := strconv.Atoi(mux.Vars(r)["id"])
	r.ParseForm()

	j, _ := json.Marshal(r.Form)
	_, err := NewSubmissionsRepository(app.db).Create(formId, string(j))

	if err != nil {
		http.Error(w, "Something bad happened.", 500)
	} else {
		http.Redirect(w, r, "/", 302)
	}
}
