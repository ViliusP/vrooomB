package util

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

//Connect connects to local mysql database
func Connect() *sql.DB {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/dev_vroom")
	DB = db
	if err != nil {
		log.Fatal("Could not connect to database")
	}

	return DB
}
