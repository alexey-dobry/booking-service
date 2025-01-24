package main

import (
	"context"
	"log"

	"github.com/alexey-dobry/booking-service/server/internal/app"
	"github.com/alexey-dobry/booking-service/server/internal/database"
	"github.com/alexey-dobry/booking-service/server/internal/logger"
)

// main starts the application
//
// It creates a database connection, creates a logger instance, creates an App instance
// and runs it.
func main() {
	db, err := database.Init()
	if err != nil {
		log.Fatal("Failed to create database connection")
	}
	defer db.Close(context.Background())

	logger := logger.NewLogger()

	a := app.New(db, logger)

	a.Run()
}
