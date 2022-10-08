package main

import "log"

func main() {

	myVar := "cat"

	switch myVar {
	case "cat":
		log.Println("cat is set to cat")

	case "dog":
		log.Println("is it a dog")

	default:
		log.Println("sd")
	}

	//
	//cat := "cat"
	//if cat == "cat" {
	//	log.Println("Cat is cat!")
	//} else {
	//	log.Println("Cat is not cat!")
	//}
	//// if else
	//var isTrue bool
	//isTrue = false
	//
	//if isTrue {
	//	log.Println("isTrue is", isTrue)
	//} else {
	//	log.Println("isTrue is ", isTrue)
	//}
	//
	//// 3
	//myNum := 100
	//
	//if myNum > 99 && !isTrue {
	//	log.Println("my nums is true to true")
	//} else if myNum < 100 && isTrue {
	//	log.Println("Go there 1")
	//} else if myNum == 101 || isTrue {
	//	log.Println("Go there 2")
	//}
}
