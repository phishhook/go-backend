package models

import "time"

type Link struct {
	ID         int       `json:"link_id,omitempty"`
	UserId     int       `json:"user_id"`
	Url        string    `json:"url"`
	ClickedAt  time.Time `json:"clicked_at"`
	IsPhishing string    `json:"is_phishing"`
}
