package models

import (
	"time"

	_ "github.com/alexey-dobry/booking-service/server/internal/validator"
)

// @Description User is a struct which contains Id, Username, Password, CreatedAt and UpdatedAt
type User struct {
	Id        int       `json:"id" validate:"required"`
	Username  string    `json:"username" validate:"required,min=6,max=20"`
	Password  string    `json:"password" validate:"required,len=14"`
	CreatedAt time.Time `json:"created_at" validate:"required"`
	UpdatedAt time.Time `json:"updated_at" validate:"required"`
}
