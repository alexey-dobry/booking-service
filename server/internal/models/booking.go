package models

import (
	_ "github.com/alexey-dobry/booking-service/server/internal/validator"
)

type Booking struct {
	Id        int    `json:"id" validate:"required,min=1,max=6"`
	UserId    int    `json:"user_id" validate:"required,min=1,max=6"`
	StartTime string `json:"start_time" validate:"required"`
	EndTime   string `json:"end_time" validate:"required"`
}
