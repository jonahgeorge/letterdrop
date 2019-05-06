package main

import (
	"database/sql"
	"net/http"
	"os"
	"log"

	"github.com/flosch/pongo2"
	"github.com/gorilla/sessions"
	"github.com/haisum/recaptcha"
	repo "github.com/jonahgeorge/letterdrop/repositories"
	_ "github.com/lib/pq"
	"github.com/sendgrid/sendgrid-go"
)

type Application struct {
	db              *sql.DB
	sessions        *sessions.CookieStore
	emailClient     *sendgrid.Client
	recaptchaClient recaptcha.R
	hostName        string
}

func NewApplication() *Application {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	sessions := sessions.NewCookieStore([]byte(os.Getenv("SECRET_TOKEN")))
	emailClient := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	recaptchaClient := recaptcha.R{Secret: os.Getenv("RECAPTCHA_SECRET_TOKEN")}
	hostName := os.Getenv("HOST")

	return &Application{
		db:              db,
		sessions:        sessions,
		emailClient:     emailClient,
		recaptchaClient: recaptchaClient,
		hostName:        hostName,
	}
}

func (app *Application) Render(w http.ResponseWriter, r *http.Request, name string, data pongo2.Context) error {
	t, _ := pongo2.FromFile("templates/" + name + ".html")

	session, _ := app.GetSession(r)

	if session.Values["userId"] != nil {
		user, _ := repo.NewUsersRepository(app.db).FindById(session.Values["userId"].(int))
		data["currentUser"] = user
	}

	data["flashes"] = session.Flashes()
	data["recaptcha_site_key"] = os.Getenv("RECAPTCHA_SITE_KEY")
	session.Save(r, w)

	return t.ExecuteWriter(data, w)
}

func (app *Application) GetSession(r *http.Request) (*sessions.Session, error) {
	return app.sessions.Get(r, "letterdrop")
}
