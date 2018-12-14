package store

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func init() {
	d, err := sql.Open("sqlite3", "uber_moto.db")
	if err != nil {
		log.Fatal("Couldn't establish database connection:", err)
	}
	db = d
}
