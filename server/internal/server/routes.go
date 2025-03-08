package server

import (
	"net/http"

	_ "github.com/alexey-dobry/booking-service/server/docs"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func (s *Server) initRoutes() {
	if s.router == nil {
		s.logger.Error("routes init error: router isn't initialized")
	}

	s.router.HandleFunc("/user", s.handleAddUser()).Methods("POST")
	s.router.HandleFunc("/user/{id}", s.handleGetUser()).Methods("GET")
	s.router.HandleFunc("/users", s.handleGetUsers()).Methods("GET")
	s.router.HandleFunc("/user/{id}", s.handleUpdateUser()).Methods("PUT")
	s.router.HandleFunc("/user/{id}", s.handleDeleteUser()).Methods("DELETE")

	s.router.HandleFunc("/booking", s.handleAddBooking()).Methods("POST")
	s.router.HandleFunc("/booking/{id}", s.handleGetBooking()).Methods("GET")
	s.router.HandleFunc("/bookings", s.handleGetBookings()).Methods("GET")
	s.router.HandleFunc("/booking/{id}", s.handleUpdateBooking()).Methods("PUT")
	s.router.HandleFunc("/booking/{id}", s.handleDeleteBooking()).Methods("DELETE")

	s.router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8000/swagger/doc.json"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)).Methods(http.MethodGet)

	s.logger.Debug("Server routes was initialized")
}
