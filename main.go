package main

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
)

func main() {
	// connect to a database
	conn, err := sql.Open("pgx", "host=localhost port=55000 password=postgrespw dbname=test_connect user=postgres")
	if err != nil {
		log.Fatal(fmt.Sprintf("Unable to connect: %v\n", err))
	}
	defer conn.Close()

	log.Println("Connected to database!")
	// test my connection
	err = conn.Ping()
	if err != nil {
		log.Fatal("cannot ping database!")
	}
	log.Println("pinged database successfully")
	// get rows from table
	err = getAllRows(conn)
	if err != nil {
		log.Fatal(err)
	}
	// insert a row
	query := "INSERT INTO users (first_name, last_name) values  " +
		"($1, $2)"
	_, err = conn.Exec(query, "Jack", "Brown")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Inserted a row")
	// get rows from table
	err = getAllRows(conn)
	if err != nil {
		log.Fatal(err)
	}
	// update a row
	stmt := "UPDATE users SET " +
		"first_name = $1  " +
		"where id = $2"
	_, err = conn.Exec(stmt, "Jackie", 4)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("updated a row")
	// get the row from table
	err = getAllRows(conn)
	if err != nil {
		log.Fatal(err)
	}
	// get one row by id
	query = "SELECT id, first_name, last_name from users " +
		"where id = $1"
	var firstName, lastName string
	var id int
	row := conn.QueryRow(query, 1)
	err = row.Scan(&id, &firstName, &lastName)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Query row returns", id, firstName, lastName)

	// delete a row
	query = "DELETE FROM users where id = $1"
	_, err = conn.Exec(query, 6)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("deleted a row")
	// get rows
	err = getAllRows(conn)
	if err != nil {
		log.Fatal(err)
	}
}

func getAllRows(conn *sql.DB) error {
	rows, err := conn.Query("select id, first_name, last_name from users")
	if err != nil {
		log.Println(err)
		return err
	}
	defer rows.Close()
	var firstName, lastName string
	var id int

	for rows.Next() {
		err := rows.Scan(&id, &firstName, &lastName)
		if err != nil {
			log.Println(err)
			return err
		}
		fmt.Println("Record is", id, firstName, lastName)
	}
	if err = rows.Err(); err != nil {
		log.Fatal("error scanning rows", err)
	}

	fmt.Println("-------------------------------------")
	return nil
}
