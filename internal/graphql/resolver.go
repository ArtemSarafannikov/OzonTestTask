package graphql

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
	cstErrors "github.com/ArtemSarafannikov/OzonTestTask/internal/errors"
	"github.com/ArtemSarafannikov/OzonTestTask/internal/service"
	"github.com/ArtemSarafannikov/OzonTestTask/internal/utils"
)

type Resolver struct {
	PostService    *service.PostService
	UserService    *service.UserService
	CommentService *service.CommentService
	PubSub         *service.PubSub
}

func NewRootResolver(postService *service.PostService,
	userService *service.UserService,
	commentService *service.CommentService,
	pubSub *service.PubSub) Config {
	c := Config{
		Resolvers: &Resolver{
			PostService:    postService,
			UserService:    userService,
			CommentService: commentService,
			PubSub:         pubSub,
		},
	}

	countComplexity := func(childComplexity int, limit *int, offset *int) int {
		return *limit * childComplexity
	}
	c.Complexity.Comment.Replies = countComplexity
	c.Complexity.Post.Comments = countComplexity

	countComplexityComments := func(childComplexity int, postId string, authorID *string, limit *int, offset *int) int {
		return *limit * childComplexity
	}
	c.Complexity.Query.Comments = countComplexityComments

	c.Directives.IsAuthenticated = func(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
		ctxUserID := ctx.Value(utils.UserIdCtxKey)
		if ctxUserID == nil {
			return nil, cstErrors.UnauthorizedError
		}
		return next(ctx)
	}

	return c
}
