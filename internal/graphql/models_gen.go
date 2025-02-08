// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package graphql

import (
	"github.com/ArtemSarafannikov/OzonTestTask/internal/models"
)

type AuthPayload struct {
	Token string       `json:"token"`
	User  *models.User `json:"user"`
}

type CreateCommentInput struct {
	Text            string  `json:"text"`
	PostID          string  `json:"postID"`
	ParentCommentID *string `json:"parentCommentID,omitempty"`
}

type CreatePostInput struct {
	Title         string `json:"title"`
	Content       string `json:"content"`
	AllowComments *bool  `json:"allowComments,omitempty"`
}

type EditPostInput struct {
	PostID        string  `json:"postID"`
	Title         *string `json:"title,omitempty"`
	Content       *string `json:"content,omitempty"`
	AllowComments *bool   `json:"allowComments,omitempty"`
}

type Mutation struct {
}

type Query struct {
}
