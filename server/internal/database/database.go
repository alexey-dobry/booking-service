package database

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
)

func NewSQL() (*sql.DB, error) {
	dsn := ""

	var db *sql.DB
	var err error

	maxRetries := 10
	delay := 3 * time.Second

	for i := range maxRetries {
		db, err = sql.Open("postgres", dsn)
		if err == nil {
			break
		}

		log.Printf("Database connection retry: %d of %d", i+1, maxRetries)

		time.Sleep(delay)
	}

	if err != nil {
		log.Fatal("Unable to connect")
	} else {
		log.Print("Successfully connected")
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Bad connection: %s", err)
	} else {
		log.Print("Connection is good")
	}

	return db, nil
}
