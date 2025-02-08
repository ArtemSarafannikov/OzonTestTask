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
}

func NewRootResolver(postService *service.PostService, userService *service.UserService, commentService *service.CommentService) Config {
	c := Config{
		Resolvers: &Resolver{
			PostService:    postService,
			UserService:    userService,
			CommentService: commentService,
		},
	}

	c.Directives.IsAuthenticated = func(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
		ctxUserID := ctx.Value(utils.UserIdCtxKey)
		if ctxUserID == nil {
			return nil, cstErrors.UnauthorizedError
		}
		return next(ctx)
	}

	c.Directives.Length = func(ctx context.Context, obj interface{}, next graphql.Resolver, max int) (interface{}, error) {
		return next(ctx)
	}

	return c
}
