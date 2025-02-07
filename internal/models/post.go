package models

import "time"

type Post struct {
	ID            string     `json:"id"`
	Title         string     `json:"title"`
	Content       string     `json:"content"`
	AllowComments bool       `json:"allow_comments"`
	AuthorID      string     `json:"author_id"`
	EditedAt      *time.Time `json:"edited_at"`
	CreatedAt     time.Time  `json:"created_at"`
}
