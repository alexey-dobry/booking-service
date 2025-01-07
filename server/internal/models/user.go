package models

import (
	"time"

	_ "github.com/alexey-dobry/booking-service/server/internal/validator"
)

type User struct {
	Id        int       `json:"id" validate:"required"`
	Username  string    `json:"username" validate:"required"`
	Password  string    `json:"password" validate:"required"`
	CreatedAt time.Time `json:"created_at" validate:"required"`
	UpdatedAt time.Time `json:"updated_at" validate:"required"`
}
