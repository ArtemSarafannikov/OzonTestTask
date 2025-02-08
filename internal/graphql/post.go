package graphql

import (
	"context"
	"github.com/ArtemSarafannikov/OzonTestTask/internal/models"
	"github.com/ArtemSarafannikov/OzonTestTask/internal/utils"
)

func (r *Resolver) Post() PostResolver { return &postResolver{r} }

type postResolver struct{ *Resolver }

func (r *postResolver) Author(ctx context.Context, obj *models.Post) (*models.User, error) {
	return r.UserService.GetUserByID(ctx, obj.AuthorID)
}

// EditedAt is the resolver for the editedAt field.
func (r *postResolver) EditedAt(ctx context.Context, obj *models.Post) (*string, error) {
	if obj.EditedAt == nil {
		return nil, nil
	}
	timeStr := utils.ConvertTimeToString(*obj.EditedAt)
	return &timeStr, nil
}

// CreatedAt is the resolver for the createdAt field.
func (r *postResolver) CreatedAt(ctx context.Context, obj *models.Post) (string, error) {
	timeStr := utils.ConvertTimeToString(obj.CreatedAt)
	return timeStr, nil
}

// Comments is the resolver for the comments field.
func (r *postResolver) Comments(ctx context.Context, obj *models.Post, limit *int, offset *int) ([]*models.Comment, error) {
	return r.CommentService.GetComments(ctx, obj.ID, *limit, *offset)
}
