package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/alexey-dobry/booking-service/server/internal/models"
	"github.com/gorilla/mux"
)

// User handler functions
func (s *Server) handleAddUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content type", "application/json")

		var newUser models.User

		if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
			http.Error(w, "Failed to decode json", http.StatusBadRequest)
			s.logger.Error("Failed to decode json")
		}

		query := "INSERT INTO users (id,username,password,created_at,updated_at) VALUES (?,?,?,?,?)"

		_, err := s.database.Exec(context.Background(), query, newUser.Id, newUser.Username, newUser.Password, newUser.CreatedAt, newUser.UpdatedAt)
		if err != nil {
			http.Error(w, "Failed to add data to database", http.StatusInternalServerError)
			s.logger.Error("Failed to add data to database")
		}

		w.WriteHeader(http.StatusCreated)
		s.logger.Debug("Successefully added user data to database")
	}
}

func (s *Server) handleDeleteUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content type", "application/json")

		id := mux.Vars(r)["id"]

		query := fmt.Sprintf("DELETE FROM users WHERE id=%s", id)
		_, err := s.database.Exec(context.Background(), query)
		if err != nil {
			http.Error(w, "Failed to delete specified user", http.StatusBadRequest)
			s.logger.Error("Failed to delete specified user")
		}

		w.WriteHeader(http.StatusOK)
		s.logger.Debug("Successefully deleted specified user data from database")
	}
}

// Booking handler functions
func (s *Server) handleAddBooking() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content type", "application/json")

		var newBooking models.Booking

		if err := json.NewDecoder(r.Body).Decode(&newBooking); err != nil {
			http.Error(w, "Failed to decode json", http.StatusBadRequest)
			s.logger.Error("Failed to decode json")
		}

		query := "INSERT INTO bookings (id,user_id,start_time,end_time) VALUES (?,?,?,?)"

		_, err := s.database.Exec(context.Background(), query, newBooking.Id, newBooking.UserId, newBooking.StartTime, newBooking.EndTime)
		if err != nil {
			http.Error(w, "Failed to add data to database", http.StatusInternalServerError)
			s.logger.Error("Failed to add data to database")
		}

		w.WriteHeader(http.StatusCreated)
		s.logger.Debug("Successefully added booking data to database")
	}
}

func (s *Server) handleDeleteBooking() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content type", "application/json")

		id := mux.Vars(r)["id"]

		query := fmt.Sprintf("DELETE FROM bookings WHERE id=%s", id)
		_, err := s.database.Exec(context.Background(), query)
		if err != nil {
			http.Error(w, "Failed to delete specified booking", http.StatusBadRequest)
			s.logger.Error("Failed to delete specified booking")
		}

		w.WriteHeader(http.StatusOK)
		s.logger.Debug("Successefully deleted specified booking data from database")
	}
}

func (s *Server) handleGetBookings() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content type", "application/json")

		var bookingList []models.Booking

		query := "SELECT * FROM bookings"
		data, err := s.database.Query(context.Background(), query)
		if err != nil {
			http.Error(w, "Failed to retrieve data from database", http.StatusInternalServerError)
			s.logger.Error(fmt.Sprintf("Failed to retrieve data from database; additional info: %s", err))
		}

		for data.Next() {
			var booking models.Booking
			if err := data.Scan(&booking.Id, &booking.UserId, &booking.StartTime, &booking.EndTime); err != nil {
				http.Error(w, "Failed to parse data", http.StatusInternalServerError)
				s.logger.Error(fmt.Sprintf("Failed to parse data; additional info: %s", err))
			}
			bookingList = append(bookingList, booking)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(bookingList)
		s.logger.Debug("Successfully got and sent booking data")
	}
}
