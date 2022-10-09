package main

import (
	"fmt"
	"net/http"
)

const portNumber = ":8080"

func main() {

	http.HandleFunc("/", Home)
	http.HandleFunc("/about", About)

	fmt.Println(fmt.Sprintf("Starting application at port number: %s", portNumber))
	err := http.ListenAndServe(portNumber, nil)
	if err != nil {
		return
	}

}
