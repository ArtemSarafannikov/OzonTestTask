package graphql

import (
	"context"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	cstErrors "github.com/ArtemSarafannikov/OzonTestTask/internal/errors"
	"github.com/ArtemSarafannikov/OzonTestTask/internal/service"
)

var (
	UserIdCtxKey = "userID"
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
		ctxUserID := ctx.Value(UserIdCtxKey)
		if ctxUserID == nil {
			return nil, cstErrors.UnauthorizedError
		}
		return next(ctx)
	}

	c.Directives.Length = func(ctx context.Context, obj interface{}, next graphql.Resolver, max int) (interface{}, error) {
		fieldCtx := graphql.GetFieldContext(ctx)
		if fieldCtx == nil {
			// TODO: replace error
			return nil, cstErrors.InternalError
		}
		var argValue interface{}
		for _, arg := range fieldCtx.Field.Arguments {
			if arg.Name == "text" {
				argValue = arg.Value
				break
			}
		}

		if argValue == nil {
			// TODO: replace error
			return nil, cstErrors.InternalError
		}
		var str string
		switch v := argValue.(type) {
		case string:
			str = v
		case *string:
			if v != nil {
				str = *v
			} else {
				// TODO: replace error
				return nil, cstErrors.InternalError
			}
		default:
			// TODO: replace error
			return nil, cstErrors.InternalError
		}

		if len(str) > max {
			return nil, fmt.Errorf("%s: text must be less than 2000 symbols", cstErrors.TooLongContentError)
		}
		return next(ctx)
	}

	return c
}
