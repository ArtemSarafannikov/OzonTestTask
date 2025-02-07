package models

import "time"

type User struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	LastActivity time.Time `json:"last_activity"`
	CreatedAt    time.Time `json:"created_at"`
}
