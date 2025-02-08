package graphql

import (
	"context"
	"github.com/ArtemSarafannikov/OzonTestTask/internal/models"
	"github.com/ArtemSarafannikov/OzonTestTask/internal/utils"
)

func (r *Resolver) User() UserResolver { return &userResolver{r} }

type userResolver struct{ *Resolver }

func (r *userResolver) LastActivity(ctx context.Context, obj *models.User) (string, error) {
	timeStr := utils.ConvertTimeToString(obj.LastActivity)
	return timeStr, nil
}

// CreatedAt is the resolver for the createdAt field.
func (r *userResolver) CreatedAt(ctx context.Context, obj *models.User) (string, error) {
	timeStr := utils.ConvertTimeToString(obj.CreatedAt)
	return timeStr, nil
}
