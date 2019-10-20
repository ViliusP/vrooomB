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

}
