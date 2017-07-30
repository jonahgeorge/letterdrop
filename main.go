package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
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
	r.HandleFunc("/email_confirmation/new", app.EmailConfirmationsNewHandler).Methods("GET")
	r.HandleFunc("/email_confirmation", app.EmailConfirmationsCreateHandler).Methods("POST")
	r.HandleFunc("/email_confirmation", app.EmailConfirmationsShowHandler).Methods("GET")

	r.HandleFunc("/forms/{formId:[0-9]+}/submissions/{submissionId:[0-9]+}", app.RequireAuthentication(app.RequireEmailConfirmation(app.SubmissionsDestroyHandler))).Methods("DELETE")
	r.HandleFunc("/forms/{id:[0-9]+}/edit", app.RequireAuthentication(app.RequireEmailConfirmation(app.FormsEditHandler))).Methods("GET")
	r.HandleFunc("/forms/{id:[0-9]+}", app.RequireAuthentication(app.RequireEmailConfirmation(app.FormsShowHandler))).Methods("GET")
	r.HandleFunc("/forms/{id:[0-9]+}", app.RequireAuthentication(app.RequireEmailConfirmation(app.FormsDestroyHandler))).Methods("DELETE")
	r.HandleFunc("/forms/{id:[0-9]+}", app.RequireAuthentication(app.RequireEmailConfirmation(app.FormsUpdateHandler))).Methods("POST")
	r.HandleFunc("/forms/new", app.RequireAuthentication(app.RequireEmailConfirmation(app.FormsNewHandler))).Methods("GET")
	r.HandleFunc("/forms", app.RequireAuthentication(app.RequireEmailConfirmation(app.FormsIndexHandler))).Methods("GET")
	r.HandleFunc("/forms", app.RequireAuthentication(app.RequireEmailConfirmation(app.FormsCreateHandler))).Methods("POST")

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "3000"
	}

	log.Println("Listening on " + port)
	log.Fatal(http.ListenAndServe(":"+port,
		handlers.HTTPMethodOverrideHandler(
			handlers.LoggingHandler(os.Stdout, r))))
}

type AuthenticatedHandlerFunc func(http.ResponseWriter, *http.Request, *User)

func (app *Application) RequireAuthentication(next AuthenticatedHandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := app.GetSession(r)
		user, err := NewUsersRepository(app.db).FindById(session.Values["userId"].(int))
		if user == nil || err != nil {
			session.AddFlash("You must be logged in!")
			session.Save(r, w)
			http.Redirect(w, r, "/login", 307)
			return
		}

		next(w, r, user)
	})
}

func (app *Application) RequireEmailConfirmation(next AuthenticatedHandlerFunc) AuthenticatedHandlerFunc {
	return AuthenticatedHandlerFunc(func(w http.ResponseWriter, r *http.Request, currentUser *User) {
		session, _ := app.GetSession(r)

		if !currentUser.isEmailConfirmed {
			session.AddFlash("You must confirm your email address before continuing")
			session.Save(r, w)
			http.Redirect(w, r, "/email_confirmation/new", 302)
			return
		}

		next(w, r, currentUser)
	})
}
