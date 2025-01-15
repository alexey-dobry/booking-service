package server

import (
	"log"
	"net/http"

	"github.com/alexey-dobry/booking-service/server/internal/logger"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx"
)

type Server struct {
	router   *mux.Router
	database *pgx.Conn
	logger   *logger.Logger
}

func New(database *pgx.Conn, logger *logger.Logger) *Server {
	s := Server{
		router:   mux.NewRouter(),
		database: database,
		logger:   logger,
	}

	s.initRoutes()

	s.logger.Debug("Server instanse created")
	return &s
}

func (s *Server) Run() {
	log.Fatal(http.ListenAndServe(":8000", s.router))
}
