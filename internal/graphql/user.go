package graphql

import (
	"context"
	"fmt"
	"github.com/ArtemSarafannikov/OzonTestTask/internal/models"
)

func (r *Resolver) User() UserResolver { return &userResolver{r} }

type userResolver struct{ *Resolver }

func (r *userResolver) LastActivity(ctx context.Context, obj *models.User) (string, error) {
	panic(fmt.Errorf("not implemented: LastActivity - lastActivity"))
}

// CreatedAt is the resolver for the createdAt field.
func (r *userResolver) CreatedAt(ctx context.Context, obj *models.User) (string, error) {
	panic(fmt.Errorf("not implemented: CreatedAt - createdAt"))
}
