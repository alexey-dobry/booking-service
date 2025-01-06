package database

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx"
)

func NewSQL() (*pgx.Conn, error) {
	var db *pgx.Conn
	var err error

	dsn, err := pgx.ParseConnectionString("")
	if err != nil {
		log.Fatal("Error parsing dsn")
	}

	maxRetries := 10
	delay := 3 * time.Second

	for i := range maxRetries {
		db, err = pgx.Connect(dsn)
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

	err = db.Ping(context.Background())
	if err != nil {
		log.Fatalf("Bad connection: %s", err)
	} else {
		log.Print("Connection is good")
	}

	return db, nil
}
