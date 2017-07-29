package main

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

func main() {
	app := NewApplication()

	r := mux.NewRouter()
	r.HandleFunc("/", app.IndexHandler).Methods("GET")
	r.HandleFunc("/login", app.SessionsNewHandler).Methods("GET")
	r.HandleFunc("/login", app.SessionsCreateHandler).Methods("POST")
	r.HandleFunc("/logout", app.SessionsDestroyHandler).Methods("GET")
	r.HandleFunc("/signup", app.UsersNewHandler).Methods("GET")
	r.HandleFunc("/signup", app.UsersCreateHandler).Methods("POST")
	r.HandleFunc("/f/{uuid}", app.SubmissionsCreateHandler).Methods("POST")

	r.HandleFunc("/forms/{formId:[0-9]+}/submissions/{submissionId:[0-9]+}", app.RequireAuthentication(app.SubmissionsDestroyHandler)).Methods("DELETE")
	r.HandleFunc("/forms/{id:[0-9]+}/edit", app.RequireAuthentication(app.FormsEditHandler)).Methods("GET")
	r.HandleFunc("/forms/{id:[0-9]+}", app.RequireAuthentication(app.FormsShowHandler)).Methods("GET")
	r.HandleFunc("/forms/{id:[0-9]+}", app.RequireAuthentication(app.FormsDestroyHandler)).Methods("DELETE")
	r.HandleFunc("/forms/{id:[0-9]+}", app.RequireAuthentication(app.FormsUpdateHandler)).Methods("POST")
	r.HandleFunc("/forms/new", app.RequireAuthentication(app.FormsNewHandler)).Methods("GET")
	r.HandleFunc("/forms", app.RequireAuthentication(app.FormsIndexHandler)).Methods("GET")
	r.HandleFunc("/forms", app.RequireAuthentication(app.FormsCreateHandler)).Methods("POST")

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "3000"
	}

	log.Println("Listening on " + port)
	log.Fatal(http.ListenAndServe(":"+port,
		handlers.HTTPMethodOverrideHandler(handlers.LoggingHandler(os.Stdout, r))))
}

type AuthenticatedHandlerFunc func(http.ResponseWriter, *http.Request, *User)

func (app *Application) RequireAuthentication(next AuthenticatedHandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user *User
		session, err := app.GetSession(r)
		user, err = NewUsersRepository(app.db).FindById(session.Values["userId"].(int))
		if user == nil || err != nil {
			session.AddFlash("You must be logged in!")
			session.Save(r, w)
			http.Redirect(w, r, "/login", 307)
			return
		}

		next(w, r, user)
	})
}
