package main

import (
	"log"
	"time"
)

var s = "seven"

type User struct {
	FirstName   string
	LastName    string
	PhoneNumber string
	Age         int
	BirthDate   time.Time
}

func main() {
	user := User{
		FirstName: "Shubham",
		LastName:  "Rauniyar",
	}
	log.Println(user.FirstName, user.LastName, user.BirthDate)
}

// small initial letter means Private
// and capital letter means having a public accesss
func whatever() {

}

//func main() {
//	//var s2 = "six"
//	//
//	//log.Println("s is ", s)
//	//log.Println("s2 is", s2)
//	//
//	//saySomething("xxxx")
//}

//func saySomething(s3 string) (string, string) {
//	log.Println("s from the saysomething is s", s)
//	return s3, "world"
//}
