package server

func (s *Server) initRoutes() {
	if s.router == nil {
		s.logger.Error("routes init error: router isn't initialized")
	}

	s.router.HandleFunc("/user", s.handleAddUser()).Methods("POST")
	s.router.HandleFunc("/user/{id}", s.handleGetUser()).Methods("GET")
	s.router.HandleFunc("/user/{id}", s.handleUpdateUser()).Methods("PUT")
	s.router.HandleFunc("/user/{id}", s.handleDeleteUser()).Methods("DELETE")

	s.router.HandleFunc("/booking", s.handleAddBooking()).Methods("POST")
	s.router.HandleFunc("/booking/{id}", s.handleGetBooking()).Methods("GET")
	s.router.HandleFunc("/bookings", s.handleGetBookings()).Methods("GET")
	s.router.HandleFunc("/booking/{id}", s.handleUpdateBooking()).Methods("PUT")
	s.router.HandleFunc("/booking/{id}", s.handleDeleteBooking()).Methods("DELTE")

	s.logger.Debug("Server routes was initialized")
}
