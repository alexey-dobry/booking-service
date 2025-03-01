package main

import (
	"context"
	"log"

	"github.com/alexey-dobry/booking-service/server/internal/app"
	"github.com/alexey-dobry/booking-service/server/internal/database"
	"github.com/alexey-dobry/booking-service/server/internal/logger"
)

// @BasePath /api

// @title RESTful API test project for MireaCyberZone
// @description This project works with PostgresSQL. It has functionality to create users and bookings. One user can have multiple bookings.

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
