package graphql

import (
	"context"
	cstErrors "github.com/ArtemSarafannikov/OzonTestTask/internal/errors"
	"github.com/ArtemSarafannikov/OzonTestTask/internal/models"
)

func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreatePost(ctx context.Context, post CreatePostInput) (*models.Post, error) {
	return r.PostService.CreatePost(ctx, post.Title, post.Content, *post.AllowComments)
}

// CreateComment is the resolver for the createComment field.
func (r *mutationResolver) CreateComment(ctx context.Context, comment CreateCommentInput) (*models.Comment, error) {
	// TODO: relocate 2000 to const
	if len(comment.Text) > 2000 {
		return nil, cstErrors.TooLongContentError
	}
	return r.CommentService.CreateComment(ctx, comment.Text, comment.PostID, comment.ParentCommentID)
}

// EditPost is the resolver for the editPost field.
func (r *mutationResolver) EditPost(ctx context.Context, newPost EditPostInput) (*models.Post, error) {
	return r.PostService.EditPost(ctx, newPost.PostID, newPost.Title, newPost.Content, newPost.AllowComments)
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
