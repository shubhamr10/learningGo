package main

import "fmt"

func main() {
	fmt.Println("Hello world.")
	var whatToSay string
	var i int

	whatToSay = "Goodbye, cruel world!"
	fmt.Println(whatToSay)

	i = 10
	fmt.Println("i is set to ", i)

	//whatWasSaid, theOtherthingthatWasSaid := saySomething()
	//fmt.Println("This function returned", whatWasSaid, theOtherthingthatWasSaid)
}

//
//func saySomething() (string, string) {
//	return "something", "else"
//}
