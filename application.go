package main

import (
	"database/sql"
	"github.com/flosch/pongo2"
	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
	"github.com/sendgrid/sendgrid-go"
	"net/http"
	"os"
)

type Application struct {
	db          *sql.DB
	sessions    *sessions.CookieStore
	emailClient *sendgrid.Client
}

func NewApplication() *Application {
	db, _ := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	sessions := sessions.NewCookieStore([]byte(os.Getenv("SECRET_TOKEN")))
	emailClient := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))

	return &Application{
		db:          db,
		sessions:    sessions,
		emailClient: emailClient,
	}
}

func (app *Application) Render(w http.ResponseWriter, r *http.Request, name string, data pongo2.Context) error {
	t, _ := pongo2.FromFile("templates/" + name + ".html")

	session, _ := app.GetSession(r)

	if session.Values["userId"] != nil {
		user, _ := NewUsersRepository(app.db).FindById(session.Values["userId"].(int))
		data["currentUser"] = user
	}

	data["flashes"] = session.Flashes()
	session.Save(r, w)

	return t.ExecuteWriter(data, w)
}

func (app *Application) GetSession(r *http.Request) (*sessions.Session, error) {
	return app.sessions.Get(r, "letterdrop")
}
