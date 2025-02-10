package graphql

import (
	"context"
	"github.com/ArtemSarafannikov/OzonTestTask/internal/dataloaders"
	"github.com/ArtemSarafannikov/OzonTestTask/internal/models"
	"github.com/ArtemSarafannikov/OzonTestTask/internal/utils"
)

func (r *Resolver) Comment() CommentResolver { return &commentResolver{r} }

type commentResolver struct{ *Resolver }

func (r *commentResolver) Post(ctx context.Context, obj *models.Comment) (*models.Post, error) {
	loader := ctx.Value(utils.DataLoadersCtxKey).(*dataloaders.DataLoaders).PostLoader
	return loader.Load(obj.PostID)
}

// ParentComment is the resolver for the parentComment field.
func (r *commentResolver) ParentComment(ctx context.Context, obj *models.Comment) (*models.Comment, error) {
	if obj.ParentID == nil {
		return nil, nil
	}
	return r.CommentService.GetCommentByID(ctx, *obj.ParentID)
}

// Author is the resolver for the author field.
func (r *commentResolver) Author(ctx context.Context, obj *models.Comment) (*models.User, error) {
	loader := ctx.Value(utils.DataLoadersCtxKey).(*dataloaders.DataLoaders).UserLoader
	return loader.Load(obj.AuthorID)
}

// CreatedAt is the resolver for the createdAt field.
func (r *commentResolver) CreatedAt(ctx context.Context, obj *models.Comment) (string, error) {
	timeStr := utils.ConvertTimeToString(obj.CreatedAt)
	return timeStr, nil
}

// Replies is the resolver for the replies field.
func (r *commentResolver) Replies(ctx context.Context, obj *models.Comment, limit *int, offset *int) ([]*models.Comment, error) {
	//return r.CommentService.GetReplies(ctx, obj.ID, *limit, *offset)
	loader := ctx.Value(utils.DataLoadersCtxKey).(*dataloaders.DataLoaders).CommentByParentIDLoader
	return loader.Load(obj.ID)
}
