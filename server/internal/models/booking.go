package models

import "time"

type Booking struct {
	Id        int       `json:"id"`
	UserId    int       `json:"user_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}
