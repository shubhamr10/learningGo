package main

import (
	"basic-web-app/pkg/handlers"
	"fmt"
	"log"
	"net/http"
)

const portNumber = ":8080"

func main() {
	log.Println("calling main")
	http.HandleFunc("/", handlers.Home)
	http.HandleFunc("/about", handlers.About)
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {})

	fmt.Println(fmt.Sprintf("Starting application on port %s", portNumber))
	http.ListenAndServe(portNumber, nil)
}
