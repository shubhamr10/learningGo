package handlers

import (
	"basic-web-app/pkg/render"
	"log"
	"net/http"
)

// Home is the homepage handler
func Home(w http.ResponseWriter, r *http.Request) {
	log.Println("Calling Home handler")
	render.RenderTemplate(w, "home.page.tmpl")
}

// About is the about page handler
func About(w http.ResponseWriter, r *http.Request) {
	log.Println("Callign About handler")
	render.RenderTemplate(w, "about.page.tmpl")
}
