package models

import (
	"time"
)

type Booking struct {
	Id        int       `json:"id" validate:"required"`
	UserId    int       `json:"user_id" validate:"required"`
	StartTime time.Time `json:"start_time" validate:"required"`
	EndTime   time.Time `json:"end_time" validate:"required"`
}
