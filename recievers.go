package main

import "log"

type myStruct struct {
	FirstName string
}

//func printFirstName() string{
//
//}

func (m *myStruct) printFirstName() string {
	return m.FirstName
}

func main() {
	var myVar myStruct
	myVar.FirstName = "John"

	myVar2 := myStruct{
		FirstName: "Miro",
	}
	log.Println("myvar is set to ", myVar.printFirstName())
	log.Println("myvar2 is set to ", myVar2.printFirstName())

}
