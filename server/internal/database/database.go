package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

var schema = `
		CREATE TABLE IF NOT EXISTS users (
		id INT NOT NULL PRIMARY KEY,
		username TEXT NOT NULL,
		password TEXT NOT NULL,
		created_at TIMESTAMP,
		updated_at TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS bookings (
		id INT NOT NULL PRIMARY KEY,
		user_id INT NOT NULL,
		start_time TIMESTAMP,
		end_time TIMESTAMP,

		CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users (id)
			ON DELETE CASCADE
			ON UPDATE CASCADE
	);
`

// Init function creates a connection to the database and returns it. The function
// uses a simple retry logic if the connection could not be established.
func Init() (*pgx.Conn, error) {
	var db *pgx.Conn
	var err error

	err = godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Failed to load environment variables; additional info: %s", err)
	}

	db_user := os.Getenv("POSTGRES_USER")
	db_password := os.Getenv("POSTGRES_PASSWORD")
	db_name := os.Getenv("POSTGRES_DB")
	db_host := os.Getenv("POSTGRES_HOST")
	db_port := os.Getenv("POSTGRES_PORT")

	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", db_user, db_password, db_host, db_port, db_name)

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

	_, err = db.Exec(context.Background(), schema)
	if err != nil {
		log.Fatalf("Unable to execute init command; additional info:%s", err)
	} else {
		log.Print("Successefully executed init command")
	}

	return db, nil
}
