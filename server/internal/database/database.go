package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

// Init function creates a connection to the database and returns it. The function
// uses a simple retry logic if the connection could not be established.
func Init() (*pgx.Conn, error) {
	var db *pgx.Conn
	var err error

	err = godotenv.Load("../.env")
	if err != nil {
		log.Printf("Failed to load environment variables; additional info: %s", err)
	}

	db_user := os.Getenv("POSTGRES_USER")
	db_password := os.Getenv("POSTGRES_PASSWORD")
	db_name := os.Getenv("POSTGRES_DB")
	db_host := os.Getenv("POSTGRES_HOST")
	db_port := os.Getenv("POSTGRES_PORT")

	if db_user == "" {
		db_user = "user"
	}
	if db_password == "" {
		db_password = "password"
	}
	if db_name == "" {
		db_name = "postgres"
	}
	if db_host == "" {
		db_host = "localhost"
	}
	if db_port == "" {
		db_port = "5432"
	}

	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", db_user, db_password, db_host, db_port, db_name)

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
		log.Fatalf("Unable to connect; additional info: %s; conn string:%s", err, connString)
	} else {
		log.Print("Successfully connected")
	}

	err = db.Ping(context.Background())
	if err != nil {
		log.Fatalf("Bad connection: %s", err)
	} else {
		log.Print("Connection is good")
	}

	db_goose, err := sql.Open("postgres", connString)
	if err != nil {
		log.Fatalf("Error converting *pgx.Conn to *sql.DB: %v", err)
	}

	err = goose.Up(db_goose, "../migrations")
	if err != nil {
		log.Fatalf("Migrations error: %s", err)
	}
	db_goose.Close()

	return db, nil
}
