package graphql

import (
	"context"
	cstErrors "github.com/ArtemSarafannikov/OzonTestTask/internal/errors"
	"github.com/ArtemSarafannikov/OzonTestTask/internal/models"
)

func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }

func (r *queryResolver) Posts(ctx context.Context, authorID *string, limit, offset *int) ([]*models.Post, error) {
	posts, err := r.PostService.GetPosts(ctx, authorID, *limit, *offset)
	return posts, cstErrors.GetCustomError(err)
}

func (r *queryResolver) Post(ctx context.Context, postID string) (*models.Post, error) {
	post, err := r.PostService.GetPostByID(ctx, postID)
	return post, cstErrors.GetCustomError(err)
}

// Comments is the resolver for the comments field.
func (r *queryResolver) Comments(ctx context.Context, postID string, authorID *string, limit *int, offset *int) ([]*models.Comment, error) {
	comment, err := r.CommentService.GetComments(ctx, postID, authorID, *limit, *offset)
	return comment, cstErrors.GetCustomError(err)
}

// Comment is the resolver for the comment field.
func (r *queryResolver) Comment(ctx context.Context, commentID string) (*models.Comment, error) {
	comment, err := r.CommentService.GetCommentByID(ctx, commentID)
	return comment, cstErrors.GetCustomError(err)
}
