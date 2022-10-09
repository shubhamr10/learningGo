package main

import (
	"log"
	"myackage/helpers"
)

func main() {
	intChan := make(chan int)
	defer close(intChan)

	go calculateValue(intChan)
	num := <-intChan
	log.Println("num", num)
}

const numPool = 10

func calculateValue(intChan chan int) {
	randomNymber := helpers.RandomNumbers(numPool)
	intChan <- randomNymber
}

//func main() {
//	log.Println("Hello");
//	var myVar helpers.SomeType
//	myVar.TypeName = "Shubham"
//
//	log.Println("myVar", myVar);
//}
