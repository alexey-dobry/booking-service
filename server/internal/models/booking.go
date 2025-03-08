package models

import (
	"time"

	_ "github.com/alexey-dobry/booking-service/server/internal/validator"
)

// @Description Booking is a struct which contains Id, UserId, StartTime and EndTime
// needs rework: text field
type Booking struct {
	Id        int       `json:"id"`
	UserId    int       `json:"user_id"`
	StartTime time.Time `json:"start_time" validate:"required"`
	EndTime   time.Time `json:"end_time" validate:"required"`
	Text      string    `json:"text" validate:"required,max=100,excludesall=/\\#@$"`
}
