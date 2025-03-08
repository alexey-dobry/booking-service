package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/alexey-dobry/booking-service/server/internal/models"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
)

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
		//сделать обработку возврата пустого результата
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

		var builder strings.Builder

		if newBookingData.EndTime.String() != "0001-01-01 00:00:00 +0000 UTC" {
			builder.WriteString("end_time='")
			builder.WriteString(newBookingData.EndTime.Format("2006-01-02 15:04:05.000"))
			builder.WriteString("',")
		}
		if newBookingData.StartTime.String() != "0001-01-01 00:00:00 +0000 UTC" {
			builder.WriteString("start_time='")
			builder.WriteString(newBookingData.StartTime.Format("2006-01-02 15:04:05.000"))
			builder.WriteString("',")
		}
		if newBookingData.Text != "" {
			builder.WriteString("text='")
			builder.WriteString(newBookingData.Text)
			builder.WriteString("',")
		}

		query := fmt.Sprintf("UPDATE bookings SET %s WHERE id=$1", builder.String()[:len(builder.String())-1])

		_, err := s.database.Exec(context.Background(), query, id)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to execute sql command; additional info:%s; querystr: %s", err, query), http.StatusInternalServerError)
			s.logger.Error(fmt.Sprintf("Failed to execute sql command; additional info:%s", err))
			return
		}
		s.logger.Debug(query)
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
