package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

func (app *Application) SubmissionsCreateHandler(w http.ResponseWriter, r *http.Request) {
	uuid := mux.Vars(r)["uuid"]
	r.ParseForm()

	form, err := NewFormsRepository(app.db).FindByUuid(uuid)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json, err := json.Marshal(r.Form)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	_, err := NewSubmissionsRepository(app.db).Create(form.id, string(json))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	http.Redirect(w, r, "/", 302)
}
