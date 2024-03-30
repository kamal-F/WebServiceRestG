package database

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func init() {
	var err error
	DB, err = sql.Open("mysql", "root:@tcp(localhost:3306)/rongsok")
	if err != nil {
		log.Fatal(err)
	}
}
