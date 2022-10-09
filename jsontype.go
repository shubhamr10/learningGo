package main

import (
	"encoding/json"
	"log"
)

type Person struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	HairColor string `json:"hair_color"`
	HasDog    bool   `json:"has_dog"`
}

func main() {
	myJSON := `[{
	"first_name":"Clark",
	"last_name":"Kent",
	"hair_color":"black",
	"has_dog":true
},{
	"first_name":"Bruce",
	"last_name":"Wayne",
	"hair_color":"black",
	"has_dog":false
}]`

	// Read JSON and convert to struct

	var unMarshaled []Person

	// string to bytes
	err := json.Unmarshal([]byte(myJSON), &unMarshaled)
	if err != nil {
		log.Println("Error unmarshalling")
	}

	log.Printf("unmarshalled %v", unMarshaled)

	var mySlice []Person
	mySlice = append(mySlice, Person{
		FirstName: "Wally",
		LastName:  "West",
		HairColor: "red",
		HasDog:    false,
	})
	mySlice = append(mySlice, Person{
		FirstName: "Diana",
		LastName:  "Prince",
		HairColor: "Black",
		HasDog:    true,
	})

	// in production use json.Marshal
	newJson, err := json.MarshalIndent(mySlice, "", "  ")
	if err != nil {
		log.Println("Error while marshalling", err)
	}

	log.Println(string(newJson))
}
