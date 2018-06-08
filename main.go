package main

import (
	"log"
)

type User struct {
	userID int    `json:"userID"`
	email  string `json:"email"`
}

func main() {
	db := connect()

	defer db.Close()

	// Execute the query
	results, err := db.Query("SELECT userID, Email FROM users_table")
	if err != nil {
		log.Fatal(err.Error()) // proper error handling instead of panic in your app
	}

	for results.Next() {
		var user User
		// for each row, scan the result into our tag composite object
		err = results.Scan(&user.userID, &user.email)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		// and then print out the tag's Name attribute
		log.Printf(user.email)
	}

}
