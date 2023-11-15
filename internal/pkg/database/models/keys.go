package models

import "time"

type Key struct {
	ID         int       `json:"id,omitempty"`
	Key        string    `json:"key"`
	CreatedAt  time.Time `json:"created_at"`
	LastUsedAt time.Time `json:"last_used_at"`
	Active     bool      `json:"active"`
}
