package app

import (
	"log"

	"github.com/alexey-dobry/booking-service/server/internal/logger"
	"github.com/alexey-dobry/booking-service/server/internal/server"
	"github.com/jackc/pgx/v5"
)

type App struct {
	server server.Server
}

func New(database *pgx.Conn, logger *logger.Logger) *App {
	a := App{
		server: *server.New(database, logger),
	}
	log.Print("App instance created")
	return &a
}

func (a *App) Run() {
	log.Print("App is started")
	a.server.Run()
}
