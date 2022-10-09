package main

import (
	"basicwebapp/pkg/handlers"
	"fmt"
	"net/http"
)

const portNumber = ":8080"

func main() {

	http.HandleFunc("/", handlers.Home)
	http.HandleFunc("/about", handlers.About)

	fmt.Println(fmt.Sprintf("Starting application at port number: %s", portNumber))
	err := http.ListenAndServe(portNumber, nil)
	if err != nil {
		return
	}

}
