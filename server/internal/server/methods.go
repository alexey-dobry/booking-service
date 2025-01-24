package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/alexey-dobry/booking-service/server/internal/models"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

// User handler functions
func (s *Server) handleAddUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var newUser models.User

		if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
			http.Error(w, "Failed to decode json", http.StatusBadRequest)
			s.logger.Error("Failed to decode json")
		}

		password, _ := bcrypt.GenerateFromPassword([]byte(newUser.Password), 14)

		query := "INSERT INTO users (id,username,password,created_at,updated_at) VALUES (?,?,?,?,?)"

		time := time.Now().Format("2006-01-02 15:04:05")
		newUser.CreatedAt = time
		newUser.UpdatedAt = time

		_, err := s.database.Exec(context.Background(), query, newUser.Id, newUser.Username, password, newUser.CreatedAt, newUser.UpdatedAt)
		if err != nil {
			http.Error(w, "Failed to add data to database", http.StatusInternalServerError)
			s.logger.Error("Failed to add data to database")
		}

		w.WriteHeader(http.StatusOK)
		s.logger.Debug("Successefully added user data to database")
	}
}

func (s *Server) handleGetUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		id := mux.Vars(r)["id"]

		var User models.User

		query := fmt.Sprintf("SELECT * FROM users WHERE id=%s", id)

		data, err := s.database.Query(context.Background(), query)
		if err != nil {
			http.Error(w, "Failed to retrieve data from database", http.StatusInternalServerError)
			s.logger.Error(fmt.Sprintf("Failed to retrieve data from database; additional info: %s", err))
		}

		err = data.Scan(&User.Id, &User.Username, &User.Password, &User.CreatedAt, &User.UpdatedAt)
		if err != nil {
			http.Error(w, "Internal error", http.StatusInternalServerError)
			s.logger.Error(fmt.Sprintf("Failed to write data into object; additional info: %s", err))
		}

		json.NewEncoder(w).Encode(User)
		s.logger.Debug("Successfully got user data")
	}
}

func (s *Server) handleUpdateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var newUserData models.User

		if err := json.NewDecoder(r.Body).Decode(&newUserData); err != nil {
			http.Error(w, "Failed to decode json", http.StatusBadRequest)
			s.logger.Error("Failed to decode json")
		}

		password, _ := bcrypt.GenerateFromPassword([]byte(newUserData.Password), 14)

		time := time.Now().Format("2006-01-02 15:04:05")
		newUserData.UpdatedAt = time

		query := fmt.Sprintf("UPDATE cars SET username=%s, password=%s,updated_at='%s' WHERE id=%d", newUserData.Username, password, newUserData.UpdatedAt, newUserData.Id)

		_, err := s.database.Exec(context.Background(), query)
		if err != nil {
			http.Error(w, "Failed to add data to database", http.StatusInternalServerError)
			s.logger.Error("Failed to add data to database")
		}

		w.WriteHeader(http.StatusOK)
		s.logger.Debug("Successefully updated user data in database")
	}
}

func (s *Server) handleDeleteUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

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
		w.Header().Set("Content-Type", "application/json")

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

func (s *Server) handleGetBooking() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		id := mux.Vars(r)["id"]

		var Booking models.Booking

		query := fmt.Sprintf("SELECT * FROM bookings WHERE id=%s", id)

		data, err := s.database.Query(context.Background(), query)
		if err != nil {
			http.Error(w, "Failed to retrieve data from database", http.StatusInternalServerError)
			s.logger.Error(fmt.Sprintf("Failed to retrieve data from database; additional info: %s", err))
		}

		err = data.Scan(&Booking.Id, &Booking.UserId, &Booking.StartTime, &Booking.EndTime)
		if err != nil {
			http.Error(w, "Internal error", http.StatusInternalServerError)
			s.logger.Error(fmt.Sprintf("Failed to write data into object; additional info: %s", err))
		}

		json.NewEncoder(w).Encode(Booking)
		s.logger.Debug("Successfully got booking data")
	}
}

func (s *Server) handleGetBookings() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

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
		s.logger.Debug("Successfully got bookings data")
	}
}

func (s *Server) handleUpdateBooking() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var newBookingData models.Booking

		if err := json.NewDecoder(r.Body).Decode(&newBookingData); err != nil {
			http.Error(w, "Failed to decode json", http.StatusBadRequest)
			s.logger.Error(fmt.Sprintf("Failed to decode json; additional info: %s", err))
		}

		query := fmt.Sprintf("UPDATE bookings SET start_time='%s' end_time='%s' WHERE id=%d", newBookingData.StartTime, newBookingData.EndTime, newBookingData.Id)

		_, err := s.database.Exec(context.Background(), query)
		if err != nil {
			http.Error(w, "Failed to execute sql command", http.StatusInternalServerError)
			s.logger.Error(fmt.Sprintf("Failed to execute sql command; additional info:%s", err))
		}

		w.WriteHeader(http.StatusOK)
		s.logger.Debug("Successefully updated booking data in database")
	}
}

func (s *Server) handleDeleteBooking() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

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
