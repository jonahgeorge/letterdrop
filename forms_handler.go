package main

import (
	"github.com/flosch/pongo2"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
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

	_, err := NewFormsRepository(app.db).Create(id.(int), r.PostFormValue("name"), r.PostFormValue("description"))
	if err != nil {
		session.AddFlash("An error occured while creating your form")
		session.Save(r, w)
		app.Render(w, r, "forms/new", pongo2.Context{
			"name": r.PostFormValue("name"),
			// TODO Inject form instead
		})
	} else {
		session.AddFlash("Successfully created form!")
		session.Save(r, w)
		http.Redirect(w, r, "/forms", 302)
	}
}

func (app *Application) FormsShowHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	form, _ := NewFormsRepository(app.db).FindById(id)
	submissions, _ := NewSubmissionsRepository(app.db).FindByFormId(form.id)

	app.Render(w, r, "forms/show", pongo2.Context{
		"form":        form,
		"submissions": submissions,
	})
}

func (app *Application) FormsEditHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	form, _ := NewFormsRepository(app.db).FindById(id)
	app.Render(w, r, "forms/edit", pongo2.Context{
		"form": form,
	})
}

func (app *Application) FormsUpdateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	form, err := NewFormsRepository(app.db).FindById(id)

	// TODO Fix this

	if err != nil {
		app.Render(w, r, "forms/edit", pongo2.Context{
			"form": form,
		})
	} else {
		http.Redirect(w, r, "/forms/"+string(form.id), 302)
	}
}
