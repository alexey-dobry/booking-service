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

// New creates an App instance with the given database and logger.
// It returns a pointer to the new App.
func New(database *pgx.Conn, logger *logger.Logger) *App {
	a := App{
		server: *server.New(database, logger),
	}
	log.Print("App instance created")
	return &a
}

// Run starts the application by running the server instance.

func (a *App) Run() {
	log.Print("App is started")
	a.server.Run()
}
