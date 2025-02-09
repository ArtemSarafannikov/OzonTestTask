package graphql

import (
	"context"
	"github.com/ArtemSarafannikov/OzonTestTask/internal/models"
)

func (r *Resolver) Subscription() SubscriptionResolver { return &subscriptionResolver{r} }

type subscriptionResolver struct{ *Resolver }

func (r *subscriptionResolver) NewCommentPost(ctx context.Context, postID string) (<-chan *models.Comment, error) {
	ch, unsubscribe := r.PubSub.Subscribe(postID)

	go func() {
		<-ctx.Done()
		unsubscribe()
	}()

	return ch, nil
}
