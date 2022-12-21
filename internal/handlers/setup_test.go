package handlers

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/justinas/nosurf"
	"github.com/shubhamr10/learningGo/internal/config"
	"github.com/shubhamr10/learningGo/internal/driver"
	"github.com/shubhamr10/learningGo/internal/models"
	"github.com/shubhamr10/learningGo/internal/render"
)

var app config.AppConfig
var session *scs.SessionManager
var pathToTemplate = "./../../templates"
var functions = template.FuncMap{
	"humanDate":  render.HumanDate,
	"formatDate": render.FormatDate,
	"iterate":    render.Iterate,
	"add":        render.Add,
}
var infoLog *log.Logger
var errorLog *log.Logger
var db *driver.DB

func TestMain(m *testing.M) {

	// What am I going to put in the sessions
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.RoomRestriction{})
	gob.Register(map[string]int{})

	// change this to true when in production
	app.InProduction = false
	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog
	// Session
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	mailChan := make(chan models.MailData)
	app.MailChan = mailChan

	defer close(mailChan)

	listenforMail()

	tc, err := CreateTestTemplateCache()
	if err != nil {
		log.Fatal(err)
	}
	app.TemplateCache = tc
	app.UseCache = true
	repo := NewTestRepo(&app)
	NewHandlers(repo)

	render.NewRenderer(&app)

	os.Exit(m.Run())
}

func getRoutes() http.Handler {

	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	//mux.Use(WriteToConsole)
	//mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", Repo.Home)
	mux.Get("/about", Repo.About)
	mux.Get("/generals-quarters", Repo.Generals)
	mux.Get("/majors-suite", Repo.Majors)

	mux.Get("/make-reservations", Repo.MakeReservations)
	mux.Post("/make-reservations", Repo.PostMakeReservations)
	mux.Get("/reservation-summary", Repo.ReservationSummary)

	mux.Get("/search-availability", Repo.SearchAvailability)
	mux.Post("/search-availability", Repo.PostSearchAvailability)
	mux.Post("/search-availability-json", Repo.PostSearchAvailabilityJSON)

	mux.Get("/contact", Repo.Contact)

	mux.Get("/choose-room/{id}", Repo.ChooseRoom)
	mux.Get("/book-room", Repo.BookRoom)

	mux.Get("/contact", Repo.Contact)
	mux.Get("/users/login", Repo.ShowLogin)
	mux.Post("/users/login", Repo.PostShowLogin)
	mux.Get("/users/logout", Repo.Logout)

	mux.Get("/admin/dashboard", Repo.AdminDashboard)
	mux.Get("/admin/reservations-new", Repo.AdminNewReservations)
	mux.Get("/admin/reservations-all", Repo.AdminAllReservations)
	mux.Get("/admin/reservation-calendars", Repo.AdminReservationsCalender)
	mux.Post("/admin/reservation-calendars", Repo.AdminPostReservationsCalender)
	mux.Get("/admin/process-reservation/{src}/{id}/do", Repo.AdminProcessReservation)
	mux.Get("/admin/delete-reservation/{src}/{id}/do", Repo.AdminDeleteReservation)

	mux.Get("/admin/reservations/{src}/{id}/show", Repo.AdminShowReservation)
	mux.Post("/admin/reservations/{src}/{id}", Repo.AdminPostShowReservation)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))
	return mux
}

// NoSurf adds CSRF protection
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})

	return csrfHandler
}

// SessionLoad loads and saves the session on every request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

func CreateTestTemplateCache() (map[string]*template.Template, error) {
	//myCache := make(map[string]*template.Template)
	log.Println("creating template cache¯")
	myCache := map[string]*template.Template{}

	// get all the files names *page.tmpl from ./templates
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplate))
	if err != nil {
		return myCache, err
	}
	// range through all files ending with *.page.tmpl
	for _, page := range pages {
		// return last element of path, we need to get the file name
		name := filepath.Base(page)
		// ts means template set
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}
		//
		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplate))
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			//
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplate))
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}

	return myCache, nil
}

func listenforMail() {
	go func() {
		for {

			_ = <-app.MailChan
		}
	}()
}
