package graphql

import (
	"context"
	"github.com/ArtemSarafannikov/OzonTestTask/internal/models"
)

func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }

func (r *queryResolver) Posts(ctx context.Context, authorID *string, limit, offset *int) ([]*models.Post, error) {
	return r.PostService.GetPosts(ctx, authorID, *limit, *offset)
}

func (r *queryResolver) Post(ctx context.Context, postID string) (*models.Post, error) {
	return r.PostService.GetPostByID(ctx, postID)
}

// Comments is the resolver for the comments field.
func (r *queryResolver) Comments(ctx context.Context, postID string, authorID *string, limit *int, offset *int) ([]*models.Comment, error) {
	return r.CommentService.GetComments(ctx, postID, authorID, *limit, *offset)
}

// Comment is the resolver for the comment field.
func (r *queryResolver) Comment(ctx context.Context, commentID string) (*models.Comment, error) {
	return r.CommentService.GetCommentByID(ctx, commentID)
}
