package graphql

import (
	"context"
	"fmt"
	"github.com/ArtemSarafannikov/OzonTestTask/internal/models"
)

func (r *Resolver) Comment() CommentResolver { return &commentResolver{r} }

type commentResolver struct{ *Resolver }

func (r *commentResolver) Post(ctx context.Context, obj *models.Comment) (*models.Post, error) {
	panic(fmt.Errorf("not implemented: Post - post"))
}

// ParentComment is the resolver for the parentComment field.
func (r *commentResolver) ParentComment(ctx context.Context, obj *models.Comment) (*models.Comment, error) {
	panic(fmt.Errorf("not implemented: ParentComment - parentComment"))
}

// Author is the resolver for the author field.
func (r *commentResolver) Author(ctx context.Context, obj *models.Comment) (*models.User, error) {
	panic(fmt.Errorf("not implemented: Author - author"))
}

// CreatedAt is the resolver for the createdAt field.
func (r *commentResolver) CreatedAt(ctx context.Context, obj *models.Comment) (string, error) {
	panic(fmt.Errorf("not implemented: CreatedAt - createdAt"))
}

// Replies is the resolver for the replies field.
func (r *commentResolver) Replies(ctx context.Context, obj *models.Comment, limit *int, offset *int) ([]*models.Comment, error) {
	panic(fmt.Errorf("not implemented: Replies - replies"))
}
