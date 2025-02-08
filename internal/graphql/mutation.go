package graphql

import (
	"context"
	"fmt"
	"github.com/ArtemSarafannikov/OzonTestTask/internal/models"
)

func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreatePost(ctx context.Context, post CreatePostInput) (*models.Post, error) {
	panic(fmt.Errorf("not implemented: CreatePost - createPost"))
}

// CreateComment is the resolver for the createComment field.
func (r *mutationResolver) CreateComment(ctx context.Context, comment CreateCommentInput) (*models.Comment, error) {
	panic(fmt.Errorf("not implemented: CreateComment - createComment"))
}

// EditPost is the resolver for the editPost field.
func (r *mutationResolver) EditPost(ctx context.Context, newPost EditPostInput) (*models.Post, error) {
	panic(fmt.Errorf("not implemented: EditPost - editPost"))
}

func (r *mutationResolver) Register(ctx context.Context, login string, password string) (*AuthPayload, error) {
	token, user, err := r.UserService.Register(ctx, login, password)
	payload := &AuthPayload{
		Token: token,
		User:  user,
	}
	return payload, err
}

func (r *mutationResolver) Login(ctx context.Context, login string, password string) (*AuthPayload, error) {
	token, user, err := r.UserService.Login(ctx, login, password)
	payload := &AuthPayload{
		Token: token,
		User:  user,
	}
	return payload, err
}
