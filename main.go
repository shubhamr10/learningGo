package main

import (
	"errors"
	"fmt"
	"net/http"
)

const portNumber = ":8080"

func main() {

	http.HandleFunc("/", Home)
	http.HandleFunc("/about", About)
	http.HandleFunc("/divide", Divide)

	fmt.Println(fmt.Sprintf("Starting application at port number: %s", portNumber))
	err := http.ListenAndServe(portNumber, nil)
	if err != nil {
		return
	}

}

// Home is the homepage handler
func Home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is a home page!")
}

// About is the about page handler
func About(w http.ResponseWriter, r *http.Request) {
	sum := addValues(2, 2)
	fmt.Fprintf(w, "This is a about page!: %d", sum)
	//fmt.Fprintf(w, fmt.Sprintf("This is about page and 2 +2 is :%d", sum))
}

// Divide is the divide page which divides two numbers
func Divide(w http.ResponseWriter, r *http.Request) {
	values, err := divideValues(2, 0)
	if err != nil {
		fmt.Fprintf(w, "cannot divide by zero")
		return
	}
	fmt.Fprintf(w, "Dividing %f from %f will be %f", 2.0, 0.0, values)
}

// divideValues values x and y
func divideValues(x, y float32) (float32, error) {
	if y <= 0 {
		err := errors.New("cannot divide by zero")
		return 0, err
	}
	result := x / y
	return result, nil
}

// addValues adds two integer and return the sum
func addValues(x, y int) int {
	return x + y
}
