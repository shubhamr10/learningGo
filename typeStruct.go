package main

import "log"

var s = "seven"

func main() {
	var s2 = "six"

	log.Println("s is ", s)
	log.Println("s2 is", s2)

	saySomething("xxxx")
}

func saySomething(s3 string) (string, string) {
	log.Println("s from the saysomething is s", s)
	return s3, "world"
}
