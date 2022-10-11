package main

import (
	"basicwebapp/pkg/config"
	"basicwebapp/pkg/handlers"
	"basicwebapp/pkg/render"
	"fmt"
	"log"
	"net/http"
)

const portNumber = ":8080"

func main() {
	var app config.AppConfig

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal(err)
	}
	app.TemplateCache = tc
	app.UseCache = false
	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

	http.HandleFunc("/", handlers.Repo.Home)
	http.HandleFunc("/about", handlers.Repo.About)

	fmt.Println(fmt.Sprintf("Starting application at port number: %s", portNumber))
	http.ListenAndServe(portNumber, nil)

}
