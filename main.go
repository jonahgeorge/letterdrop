package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	app := NewApplication()

	router := mux.NewRouter()
	router.HandleFunc("/", app.IndexHandler).Methods("GET")

	router.HandleFunc("/login", app.SessionsNewHandler).Methods("GET")
	router.HandleFunc("/login", app.SessionsCreateHandler).Methods("POST")
	router.HandleFunc("/logout", app.SessionsDestroyHandler).Methods("GET")

	router.HandleFunc("/signup", app.UsersNewHandler).Methods("GET")
	router.HandleFunc("/signup", app.UsersCreateHandler).Methods("POST")

	router.HandleFunc("/forms", app.FormsIndexHandler).Methods("GET")
	router.HandleFunc("/forms/new", app.FormsNewHandler).Methods("GET")
	router.HandleFunc("/forms", app.FormsCreateHandler).Methods("POST")
	router.HandleFunc("/forms/{id:[0-9]+}", app.FormsShowHandler).Methods("GET")
	router.HandleFunc("/forms/{id:[0-9]+}/edit", app.FormsEditHandler).Methods("GET")
	router.HandleFunc("/forms/{id:[0-9]+}", app.FormsUpdateHandler).Methods("POST")

	router.HandleFunc("/f/{uuid}", app.SubmissionsCreateHandler).Methods("POST")

	log.Fatal(http.ListenAndServe(":5000", router))
}
