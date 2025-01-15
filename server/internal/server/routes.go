package server

func (s *Server) initRoutes() {
	if s.router == nil {
		s.logger.Error("routes init error: router isn't initialized")
	}

	s.router.HandleFunc("/api/user", s.handleAddUser()).Methods("POST")
	s.router.HandleFunc("/api/user/{id}", s.handleDeleteUser()).Methods("DELETE")
	s.router.HandleFunc("/api/booking", s.handleAddBooking()).Methods("POST")
	s.router.HandleFunc("/api/booking/{id}", s.handleDeleteBooking()).Methods("DELTE")
	s.router.HandleFunc("/api/bookings", s.handleGetBookings()).Methods("GET")

	s.logger.Error("Server routes was initialized")
}
