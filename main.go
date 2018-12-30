package main

import (
	"log"
	"net/http"

	"./routes"
	"./util"
)

func main() {
	db := util.Connect()

	defer db.Close()
	router := routes.NewRouter()
	log.Fatal(http.ListenAndServe(":8080", router))

	/*	// Execute the query
		results, err := db.Query("SELECT user_ID, email FROM users_table")
		if err != nil {
			log.Fatal(err.Error()) // proper error handling instead of panic in your app
		}

		for results.Next() {
			var user users.User
			// for each row, scan the result into our tag composite object
			err = results.Scan(&user.UserID, &user.Email)
			if err != nil {
				panic(err.Error()) // proper error handling instead of panic in your app
			}
			// and then print out the tag's Name attribute
			log.Printf(user.Email)
		}*/

}
