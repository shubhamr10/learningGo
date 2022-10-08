package main

import (
	"log"
	"sort"
)

//type Users struct {
//	FirstName string
//	LastName  string
//}

func main() {
	// Array ==> slice in golang
	var myString string
	myString = "Fish"
	log.Println("mySeeing", myString)

	var mySliceString []string
	mySliceString = append(mySliceString, "Shubham")
	mySliceString = append(mySliceString, "Akshay")
	mySliceString = append(mySliceString, "Gaurav")

	sort.Strings(mySliceString)

	log.Println(mySliceString)

	myIntSlice := []int{1, 2, 3, 4, 5, 6, 7}
	log.Println(myIntSlice[0:4])
	//// Creating a object(map) in golang
	//myMap := make(map[string]Users)
	//
	//me := Users{
	//	FirstName: "Shubham",
	//	LastName:  "Rauniyar",
	//}
	//
	//myMap["me"] = me
	//
	//log.Println("me", myMap["me"].FirstName)

	//myMap["dog"] = "Samson"
	//myMap["otherDog"] = "Cassey"
	//log.Println(myMap["dog"])
}

//func main() {
//	var myString string
//	var myInt int
//
//	myString = "Hi"
//	myInt = 11
//
//	mySecString := "another string"
//
//	log.Println("myString", myString, myInt, mySecString)
//}
