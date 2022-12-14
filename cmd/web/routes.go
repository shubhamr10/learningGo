package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/shubhamr10/learningGo/internal/config"
	"github.com/shubhamr10/learningGo/internal/handlers"
	"net/http"
)

// muxroutes create routing using "github.com/bmizerany/pat"
//func muxroutes(app *config.AppConfig) http.Handler {
//	mux := pat.New()
//
//	mux.Get("/", http.HandlerFunc(handlers.Repo.Home))
//	mux.Get("/about", http.HandlerFunc(handlers.Repo.About))
//
//	return mux
//}

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	//mux.Use(WriteToConsole)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)
	mux.Get("/generals-quarters", handlers.Repo.Generals)
	mux.Get("/majors-suite", handlers.Repo.Majors)

	mux.Get("/make-reservations", handlers.Repo.MakeReservations)
	mux.Post("/make-reservations", handlers.Repo.PostMakeReservations)
	mux.Get("/reservation-summary", handlers.Repo.ReservationSummary)

	mux.Get("/search-availability", handlers.Repo.SearchAvailability)
	mux.Post("/search-availability", handlers.Repo.PostSearchAvailability)
	mux.Post("/search-availability-json", handlers.Repo.PostSearchAvailabilityJSON)
	mux.Get("/choose-room/{id}", handlers.Repo.ChooseRoom)
	mux.Get("/book-room", handlers.Repo.BookRoom)

	mux.Get("/contact", handlers.Repo.Contact)
	mux.Get("/users/login", handlers.Repo.ShowLogin)
	mux.Post("/users/login", handlers.Repo.PostShowLogin)

	mux.Get("/users/logout", handlers.Repo.Logout)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	mux.Route("/admin", func(mux chi.Router) {
		//mux.Use(Auth)

		mux.Get("/dashboard", handlers.Repo.AdminDashboard)
		mux.Get("/reservations-new", handlers.Repo.AdminNewReservations)
		mux.Get("/reservations-all", handlers.Repo.AdminAllReservations)
		mux.Get("/reservation-calendars", handlers.Repo.AdminReservationsCalender)
		mux.Post("/reservation-calendars", handlers.Repo.AdminPostReservationsCalender)
		mux.Get("/process-reservation/{src}/{id}/do", handlers.Repo.AdminProcessReservation)
		mux.Get("/delete-reservation/{src}/{id}/do", handlers.Repo.AdminDeleteReservation)

		mux.Get("/reservations/{src}/{id}/show", handlers.Repo.AdminShowReservation)
		mux.Post("/reservations/{src}/{id}", handlers.Repo.AdminPostShowReservation)

	})
	return mux
}
