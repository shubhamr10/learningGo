package main

import "log"

func main() {
	type User struct {
		FirstName string
		LastName  string
		Email     string
		Age       int
	}
	var users []User
	users = append(users, User{"John", "Smith", "john@smith.com", 30})
	users = append(users, User{"John", "Cena", "john@cena.com", 30})
	users = append(users, User{"Mike", "Tyson", "mike@tyson.com", 30})
	users = append(users, User{"Ruchi", "Prakash", "ruchi@prakash.com", 30})

	for _, user := range users {
		log.Println("FirstName", user.FirstName, "LastName", user.LastName, "Email", user.Email, "Age", user.Age)
	}

	// string but it will show the byte
	//var myPoem = "Once upon a time, there was a guy named Shubham"
	//
	//for i, chars := range myPoem {
	//	log.Println("i", i, chars);
	//}

	// maps
	//animals := make(map[string]string)
	//animals["dog"] = "Tommy"
	//animals["horse"] = "Cecile"
	//
	//for animalType, name := range animals {
	//	log.Println("A ", animalType, " has a name ", name);
	//}

	// Slices
	//animals := []string{"dog", "horse", "cow", "cat", "tiger", "lion"}
	//
	//for i, animal := range animals {
	//	log.Println("i", i, "animal", animal)
	//}
	//
	//for _, animal := range animals {
	//	log.Println("animal", animal)
	//}
	//for i := 0; i <= 10; i++ {
	//	log.Println(i)
	//}
}
