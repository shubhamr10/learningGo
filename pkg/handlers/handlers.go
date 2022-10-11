package handlers

import (
	"basicwebapp/pkg/config"
	"basicwebapp/pkg/models"
	"basicwebapp/pkg/render"
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
	render.RenderTemplate(w, "home.page.tmpl", &models.TemplateData{})
}

// About is the about page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	// perform some logic
	stringmap := make(map[string]string)
	stringmap["test"] = "Hello! Again"
	render.RenderTemplate(w, "about.page.tmpl", &models.TemplateData{
		StringMap: stringmap,
	})
}
