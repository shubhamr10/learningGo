package main

import "log"

func main() {
	var myString string
	myString = "Green"

	log.Println("my string is set to ", myString)
	changeUsingPointer(&myString)
	log.Println("after calling the function my string is set to ", myString)
}

func changeUsingPointer(s *string) {
	log.Println("s is set to", s)
	newValue := "Red"
	*s = newValue
}
