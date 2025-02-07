package graphql

import (
	"context"
	"fmt"
	"github.com/ArtemSarafannikov/OzonTestTask/internal/models"
)

func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }

func (r *queryResolver) Posts(ctx context.Context, authorID *string, limit, offset *int) ([]*models.Post, error) {
	return r.PostService.GetPosts(ctx, *limit, *offset)
}

func (r *queryResolver) Post(ctx context.Context, postID string) (*models.Post, error) {
	panic(fmt.Errorf("not implemented: Post - post"))
}

// Comments is the resolver for the comments field.
func (r *queryResolver) Comments(ctx context.Context, postID string, authorID *string, limit *int, offset *int) ([]*models.Comment, error) {
	panic(fmt.Errorf("not implemented: Comments - comments"))
}

// Comment is the resolver for the comment field.
func (r *queryResolver) Comment(ctx context.Context, commentID string) (*models.Comment, error) {
	panic(fmt.Errorf("not implemented: Comment - comment"))
}
