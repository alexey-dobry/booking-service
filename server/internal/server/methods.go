package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"reflect"

	"github.com/alexey-dobry/booking-service/server/internal/models"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

// handleAddUser
//
// @Summary Add new user to database
// @Description Creates functon which adds new user data to database
// @Accept json
//
// @Param Username formData string true "6 <= length <= 20"
// @Param password formData string true "length = 14"
// @Param CreatedAt formData string true "format = YYYY-MM-DD HH:MM:SS"
// @Param UpdatedAt formData string true "format = YYYY-MM-DD HH:MM:SS"
//
// @Success 200 {object} integer "ok"
// @Failure 400 {object} integer "Wrong ID"
// @Failure 500 {object} integer "Error scanning data from db response"
// @Router /user [post]
func (s *Server) handleAddUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var newUser models.User

		if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
			http.Error(w, fmt.Sprintf("Failed to decode json: %s", err), http.StatusBadRequest)
			s.logger.Error("Failed to decode json")
			return
		}

		password, _ := bcrypt.GenerateFromPassword([]byte(newUser.Password), 14)

		query := "INSERT INTO users (username,password,created_at,updated_at) VALUES ($1,$2,$3,$4)"

		time := time.Now()
		newUser.CreatedAt = time
		newUser.UpdatedAt = time

		_, err := s.database.Exec(context.Background(), query, newUser.Username, password, newUser.CreatedAt, newUser.UpdatedAt)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to add data to database; additional info: %s", err), http.StatusInternalServerError)
			s.logger.Error(fmt.Sprintf("Failed to add data to database; additional info: %s", err))
			return
		}

		w.WriteHeader(http.StatusOK)
		s.logger.Debug("Successefully added user data to database")
	}
}

// handleGetUser
//
// @Summary Get user data
// @Description Creates function which retrieves data of user specified by id from database
// @Produces json
//
// @Param id path int true "User ID "
//
// @Success 200 {object} models.User
// @Failure 400 {object} integer "Wrong ID"
// @Failure 500 {object} integer "Error scanning data from db response"
// @Router /user/{id} [get]
func (s *Server) handleGetUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		id, _ := strconv.Atoi(mux.Vars(r)["id"])

		var User models.User

		query := "SELECT id, username, password, created_at, updated_at FROM users WHERE id=$1;"

		data := s.database.QueryRow(context.Background(), query, id)

		err := data.Scan(&User.Id, &User.Username, &User.Password, &User.CreatedAt, &User.UpdatedAt)

		if err == pgx.ErrNoRows {
			http.Error(w, fmt.Sprintf("No entry with id {%d} was found in database", id), http.StatusBadRequest)
			s.logger.Error(fmt.Sprintf("No entry with id {%d} was found in database", id))
			return
		} else if err != nil {
			http.Error(w, fmt.Sprintf("Internal error; more info: %s", err), http.StatusInternalServerError)
			s.logger.Error(fmt.Sprintf("Internal error; more info: %s", err))
			return
		}

		json.NewEncoder(w).Encode(User)
		s.logger.Debug("Successfully retrieved user data")
	}
}

// handleUpdateUser
//
// @Summary Update user data
// @Description Creates function which updates data of user specified by id in database
// @Accept json
//
// @Param id path int true "User ID"
//
// @Success 200 {object} integer "ok"
// @Failure 400 {object} integer "Wrong ID"
// @Failure 500 {object} integer "Error scanning data from db response"
// @Router /user/{id} [put]
func (s *Server) handleUpdateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		id, _ := strconv.Atoi(mux.Vars(r)["id"])

		var newUserData models.User

		if err := json.NewDecoder(r.Body).Decode(&newUserData); err != nil {
			http.Error(w, fmt.Sprintf("Failed to decode json; additional info: %s", err), http.StatusBadRequest)
			s.logger.Error(fmt.Sprintf("Failed to decode json; additional info: %s", err))
			return
		}

		password, _ := bcrypt.GenerateFromPassword([]byte(newUserData.Password), 14)

		time := time.Now()
		newUserData.UpdatedAt = time

		var builder strings.Builder

		if newUserData.Password != "" {
			builder.WriteString("password='")
			builder.WriteString(string(password))
			builder.WriteString("',")
		}
		if newUserData.Username != "" {
			builder.WriteString("username=")
			builder.WriteString(newUserData.Username)
			builder.WriteString(",")
		}

		query := fmt.Sprintf("UPDATE users SET %s updated_at=$1 WHERE id=$2", builder.String())

		_, err := s.database.Exec(context.Background(), query, newUserData.UpdatedAt, id)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to add data to database; additional info: %s; querystr: %s", err, query), http.StatusInternalServerError)
			s.logger.Error(fmt.Sprintf("Failed to add data to database; additional info: %s", err))
			return
		}

		w.WriteHeader(http.StatusOK)
		s.logger.Debug("Successefully updated user data in database")
	}
}

// handleDeleteUser
//
// @Summary Delete specified user data
// @Description Creates function which deletes data of user specified by id from database
//
// @Param id path int true "User ID"
//
// @Success 200 {object} integer "ok"
// @Failure 400 {object} integer "Wrong Id"
// @Router /user/{id} [delete]
func (s *Server) handleDeleteUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		id := mux.Vars(r)["id"]

		query := "DELETE FROM users WHERE id=$1"
		_, err := s.database.Exec(context.Background(), query, id)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to delete specified user; additional info: %s", err), http.StatusBadRequest)
			s.logger.Error(fmt.Sprintf("Failed to delete specified user; additional info: %s", err))
			return
		}

		w.WriteHeader(http.StatusOK)
		s.logger.Debug("Successefully deleted specified user data from database")
	}
}

// handleAddBooking
//
// @Summary Adds new booking entry
// @Description Creates function which adds new user data to database
// @Accept json
//
// @Param UserId formData int true "integer >= 1"
// @Param StartTime formData string true "format = YYYY-MM-DD HH:MM:SS"
// @Param EndTime formData string true "format = YYYY-MM-DD HH:MM:SS"
//
// @Success 200 {object} integer "ok"
// @Failure 400 {object} integer "Wrong ID"
// @Failure 500 {object} integer "Error scanning data from db response"
// @Router /booking [post]
func (s *Server) handleAddBooking() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var newBooking models.Booking

		if err := json.NewDecoder(r.Body).Decode(&newBooking); err != nil {
			http.Error(w, fmt.Sprintf("Failed to decode json; additional info: %s", err), http.StatusBadRequest)
			s.logger.Error(fmt.Sprintf("Failed to decode json; additional info: %s", err))
			return
		}

		query := "INSERT INTO bookings (user_id,start_time,end_time,text) VALUES ($1,$2,$3,$4)"

		_, err := s.database.Exec(context.Background(), query, newBooking.UserId, newBooking.StartTime, newBooking.EndTime, newBooking.Text)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to add data to database; additional info: %s", err), http.StatusInternalServerError)
			s.logger.Error(fmt.Sprintf("Failed to add data to database; additional info: %s", err))
			return
		}

		w.WriteHeader(http.StatusCreated)
		s.logger.Debug("Successefully added booking data to database")
	}
}

// handleGetBooking
//
// @Summary Get booking data
// @Description Creates function which retrieves data of booking specified by id from database
// @Produces json
//
// @Param id path int true "Booking ID"
//
// @Success 200 {object} models.Booking "ok"
// @Failure 500 {object} integer "Error scanning data from db response"
// @Router /booking/{id} [get]
func (s *Server) handleGetBooking() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		id, _ := strconv.Atoi(mux.Vars(r)["id"])

		var Booking models.Booking

		query := "SELECT * FROM bookings WHERE id=$1"

		data := s.database.QueryRow(context.Background(), query, id)

		err := data.Scan(&Booking.Id, &Booking.UserId, &Booking.StartTime, &Booking.EndTime, &Booking.Text)
		if err == pgx.ErrNoRows {
			http.Error(w, fmt.Sprintf("No entry with id {%d} was found in database", id), http.StatusBadRequest)
			s.logger.Error(fmt.Sprintf("No entry with id {%d} was found in database", id))
			return
		} else if err != nil {
			http.Error(w, fmt.Sprintf("Internal error; more info: %s", err), http.StatusInternalServerError)
			s.logger.Error(fmt.Sprintf("Internal error; more info: %s", err))
			return
		}

		json.NewEncoder(w).Encode(Booking)
		s.logger.Debug("Successfully retrieved booking data")
	}
}

// handleGetBookings
//
// @Summary Get booking data
// @Description Creates function which retrieves data of all bookings from database
// @Produces json
//
// @Success 200 {array} models.Booking "ok"
// @Failure 500 {object} integer "Error scanning data from db response"
// @Router /bookings [get]
func (s *Server) handleGetBookings() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var bookingList []models.Booking

		query := "SELECT * FROM bookings"
		data, err := s.database.Query(context.Background(), query)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to retrieve data from database; additional info: %s", err), http.StatusInternalServerError)
			s.logger.Error(fmt.Sprintf("Failed to retrieve data from database; additional info: %s", err))
			return
		}

		for data.Next() {
			var Booking models.Booking
			err = data.Scan(&Booking.Id, &Booking.UserId, &Booking.StartTime, &Booking.EndTime, &Booking.Text)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to write data into object; additional info: %s", err), http.StatusInternalServerError)
				s.logger.Error(fmt.Sprintf("Failed to write data into object; additional info: %s", err))
				return
			}
			bookingList = append(bookingList, Booking)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(bookingList)
		s.logger.Debug("Successfully retrieved bookings data")
	}
}

// handleUpdateBooking
//
// @Summary Updates booking data
// @Description Creates function which updates data of booking specified by id in database
// @Accept json
//
// @Param id path int true "Booking ID"
//
// @Success 200 {object} integer "ok"
// @Failure 400 {object} integer "Wrong Id"
// @Failure 500 {object} integer "Error scanning data from db response"
// @Router /booking/{id} [put]
func (s *Server) handleUpdateBooking() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		id, _ := strconv.Atoi(mux.Vars(r)["id"])

		var newBookingData models.Booking

		if err := json.NewDecoder(r.Body).Decode(&newBookingData); err != nil {
			http.Error(w, fmt.Sprintf("Failed to decode json; additional info: %s", err), http.StatusBadRequest)
			s.logger.Error(fmt.Sprintf("Failed to decode json; additional info: %s", err))
			return
		}

		v := reflect.ValueOf(newBookingData)
		t := reflect.TypeOf(newBookingData)

		var builder strings.Builder

		for i := 0; i < v.NumField(); i++ {
			if v.Field(i).Interface() != "" && v.Field(i).Interface() != 0 && fmt.Sprint(v.Field(i).Interface()) != "0001-01-01 00:00:00 +0000 UTC" {
				builder.WriteString(t.Field(i).Tag.Get("json"))
				builder.WriteString("='")
				builder.WriteString(fmt.Sprint(v.Field(i).Interface()))
				builder.WriteString("',")
			}
		}

		query := fmt.Sprintf("UPDATE bookings SET %s WHERE id=$1", builder.String()[:len(builder.String())-1])

		_, err := s.database.Exec(context.Background(), query, id)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to execute sql command; additional info:%s; querystr: %s", err, query), http.StatusInternalServerError)
			s.logger.Error(fmt.Sprintf("Failed to execute sql command; additional info:%s", err))
			return
		}

		w.WriteHeader(http.StatusOK)
		s.logger.Debug("Successefully updated booking data in database")
	}
}

// handleDeleteBooking
//
// @Summary Delete specified booking data
// @Description Creates function which deletes data of booking specified by id from database
//
// @Param id path int true "Booking ID"
//
// @Success 200 {object} integer "ok"
// @Failure 400 {object} integer "Wrong Id"
// @Router /booking/{id} [delete]
func (s *Server) handleDeleteBooking() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		id := mux.Vars(r)["id"]

		query := "DELETE FROM bookings WHERE id=$1"
		_, err := s.database.Exec(context.Background(), query, id)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to delete specified booking; additional info: %s", err), http.StatusBadRequest)
			s.logger.Error(fmt.Sprintf("Failed to delete specified booking; additional info: %s", err))
			return
		}

		w.WriteHeader(http.StatusOK)
		s.logger.Debug("Successefully deleted specified booking data from database")
	}
}
