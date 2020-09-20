package Model

import (
	"database/sql"
	"log"
)

var db *sql.DB

func Init() {
	var err error
	db, err = sql.Open("mysql", "root:123456@tcp(127.0.0.1:3361)/")
	if err != nil {
		log.Fatal("Database Connected Failed!")
	}
}

func Close() {
	db.Close()
}