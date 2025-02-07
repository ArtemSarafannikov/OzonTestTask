package graphql

import (
	"context"
	"fmt"
	"github.com/ArtemSarafannikov/OzonTestTask/internal/models"
)

func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreatePost(ctx context.Context, title string, content string, authorID string, allowComments *bool) (*models.Post, error) {
	panic(fmt.Errorf("not implemented: CreatePost - createPost"))
}

// CreateComment is the resolver for the createComment field.
func (r *mutationResolver) CreateComment(ctx context.Context, text string, postID string, authorID string, parentCommentID *string) (*models.Comment, error) {
	panic(fmt.Errorf("not implemented: CreateComment - createComment"))
}

// EditPost is the resolver for the editPost field.
func (r *mutationResolver) EditPost(ctx context.Context, postID string, title *string, content *string, allowComments *bool) (*models.Post, error) {
	panic(fmt.Errorf("not implemented: EditPost - editPost"))
}
