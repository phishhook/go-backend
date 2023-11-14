package models

import "time"

type User struct {
	ID          int       `json:"id,omitempty"`
	Username    string    `json:"username"`
	PhoneNumber string    `json:"phone_number"`
	CreatedAt   time.Time `json:"created_at"`
}
