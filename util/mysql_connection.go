package util

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

//Connect connects to local mysql database
func Connect() *sql.DB {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/dev_vrooom")

	if err != nil {
		log.Fatal("Could not connect to database")
	}

	return db
}
