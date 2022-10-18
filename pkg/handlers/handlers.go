package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/shubhamr10/learningGo/pkg/config"
	"github.com/shubhamr10/learningGo/pkg/models"
	"github.com/shubhamr10/learningGo/pkg/render"
	"log"
	"net/http"
)

type Repository struct {
	App *config.AppConfig
}

// Repo is the variable used by the handlers
var Repo *Repository

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHandlers set the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the homepage handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIp := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIp)

	render.RenderTemplate(w, r, "home.page.tmpl", &models.TemplateData{})
}

// About is the about page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	// perform some logic
	stringmap := make(map[string]string)
	stringmap["test"] = "Hello! Again"

	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")

	stringmap["remote_ip"] = remoteIP
	render.RenderTemplate(w, r, "about.page.tmpl", &models.TemplateData{
		StringMap: stringmap,
	})
}

// Generals is the homepage handler
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	remoteIp := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIp)

	render.RenderTemplate(w, r, "generals.page.tmpl", &models.TemplateData{})
}

// Majors is the homepage handler
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	remoteIp := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIp)

	render.RenderTemplate(w, r, "majors.page.tmpl", &models.TemplateData{})
}

// MakeReservations is the homepage handler
func (m *Repository) MakeReservations(w http.ResponseWriter, r *http.Request) {
	remoteIp := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIp)

	render.RenderTemplate(w, r, "make-reservation.page.tmpl", &models.TemplateData{})
}

// SearchAvailability is the homepage handler
func (m *Repository) SearchAvailability(w http.ResponseWriter, r *http.Request) {
	remoteIp := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIp)

	render.RenderTemplate(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

// PostSearchAvailability is the homepage handler
func (m *Repository) PostSearchAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")
	w.Write([]byte(fmt.Sprintf("posted to search from %s till %s", start, end)))
}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

// PostSearchAvailabilityJSON handles request for availability and send JSON response
func (m *Repository) PostSearchAvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	resp := jsonResponse{
		OK:      true,
		Message: "available!",
	}
	out, err := json.MarshalIndent(resp, "", "     ")
	if err != nil {
		log.Println(err)
	}
	log.Println(string(out))
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// Contact is the homepage handler
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	remoteIp := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIp)

	render.RenderTemplate(w, r, "contact.page.tmpl", &models.TemplateData{})
}
