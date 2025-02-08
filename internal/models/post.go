package models

import (
	"time"
)

type Post struct {
	ID            string     `json:"id"`
	Title         string     `json:"title"`
	Content       string     `json:"content"`
	AllowComments bool       `json:"allow_comments"`
	AuthorID      string     `json:"-"`
	EditedAt      *time.Time `json:"edited_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
}
