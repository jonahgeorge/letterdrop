package main

import (
	"database/sql"
	"github.com/flosch/pongo2"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

type Application struct {
	db       *sql.DB
	sessions *sessions.CookieStore
}

func NewApplication() *Application {
	db, _ := sql.Open("postgres", "postgres://jonahgeorge@localhost/letterbox_development?sslmode=disable")
	sessions := sessions.NewCookieStore([]byte("something-very-secret"))

	return &Application{
		db:       db,
		sessions: sessions,
	}
}

func (app *Application) Render(w http.ResponseWriter, r *http.Request, name string, data pongo2.Context) error {
	t, _ := pongo2.FromFile("templates/" + name + ".html")

	session, _ := app.sessions.Get(r, "letterbox")

	if session.Values["userId"] != nil {
		user := NewUsersRepository(app.db).FindById(session.Values["userId"].(int))
		data["currentUser"] = user
	}

	data["flashes"] = session.Flashes()
	session.Save(r, w)

	return t.ExecuteWriter(data, w)
}

func (app *Application) GetSession(r *http.Request) (*sessions.Session, error) {
	return app.sessions.Get(r, "letterbox")
}

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
	router.HandleFunc("/forms/{uuid}", app.FormsShowHandler).Methods("GET")

	log.Fatal(http.ListenAndServe(":5000", router))
}
