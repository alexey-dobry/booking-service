package models

import (
	"time"

	_ "github.com/alexey-dobry/booking-service/server/internal/validator"
)

type User struct {
	Id        int       `json:"id" validate:"required,min=1,max=6"`
	Username  string    `json:"username" validate:"required,ascii,min=10,max64"`
	Password  string    `json:"password" validate:"required,ascii,min=10,max64"`
	CreatedAt time.Time `json:"created_at" validate:"required"`
	UpdatedAt time.Time `json:"updated_at" validate:"required"`
}
