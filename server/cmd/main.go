package main

import (
	"context"
	"log"

	"github.com/alexey-dobry/booking-service/server/internal/app"
	"github.com/alexey-dobry/booking-service/server/internal/database"
	"github.com/alexey-dobry/booking-service/server/internal/logger"
)

func main() {
	db, err := database.NewSQL()
	if err != nil {
		log.Fatal("Failed to create database connection")
	}
	defer db.Close(context.Background())

	logger := logger.NewLogger()

	a := app.New(db, logger)

	a.Run()
}
