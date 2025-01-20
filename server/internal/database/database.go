package database

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
)

func NewSQL() (*pgx.Conn, error) {
	var db *pgx.Conn
	var err error

	connString := "postgres://alexnh:superpass1029@localhost:3308/service"

	maxRetries := 10
	delay := 3 * time.Second

	for i := range maxRetries {
		db, err = pgx.Connect(context.Background(), connString)
		if err == nil {
			break
		}

		log.Printf("Database connection retry: %d of %d", i+1, maxRetries)
		time.Sleep(delay)
	}

	if err != nil {
		log.Fatalf("Unable to connect; additional info: %s", err)
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
