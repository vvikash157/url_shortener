package config

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

func InitDB() *sql.DB {

	var err error
	db, err = sql.Open("postgres", os.Getenv("ConnectionString"))
	if err != nil {
		log.Fatalf("error opening connection string %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("error connecting Database %v", err)
	}

	log.Println("database connected successfully")
	return db
}
