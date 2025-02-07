package graphql

import (
	"OzonTestTask/internal/service"
	"context"
	"fmt"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	PostService *service.PostService
}

/*
type ResolverRoot interface {
    Mutation() MutationResolver
    Query() QueryResolver
}
*/

func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Posts(ctx context.Context, authorID *string, limit, offset *int) ([]*Post, error) {
	posts, err := r.PostService.GetPosts(ctx, *limit, *offset)
	if err != nil {
		panic(err)
	}

	// FIXME: remake to models.Post
	retPosts := make([]*Post, len(posts))
	for i, post := range posts {
		edited := post.EditedAt.Format("2006-01-02 15:04:05")
		retPosts[i] = &Post{
			ID:            post.ID,
			Title:         post.Title,
			Content:       post.Content,
			AllowComments: post.AllowComments,
			Author:        &User{},
			EditedAt:      &edited,
			CreatedAt:     post.CreatedAt.Format("2006-01-02 15:04:05"),
			Comments:      make([]*Comment, 0),
		}
	}
	return retPosts, nil
}

func (r *queryResolver) Post(ctx context.Context, postID string) (*Post, error) {
	panic(fmt.Errorf("not implemented: Post - post"))
}

// Comments is the resolver for the comments field.
func (r *queryResolver) Comments(ctx context.Context, postID string, authorID *string, limit *int, offset *int) ([]*Comment, error) {
	panic(fmt.Errorf("not implemented: Comments - comments"))
}

// Comment is the resolver for the comment field.
func (r *queryResolver) Comment(ctx context.Context, commentID string) (*Comment, error) {
	panic(fmt.Errorf("not implemented: Comment - comment"))
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreatePost(ctx context.Context, title string, content string, authorID string, allowComments *bool) (*Post, error) {
	panic(fmt.Errorf("not implemented: CreatePost - createPost"))
}

// CreateComment is the resolver for the createComment field.
func (r *mutationResolver) CreateComment(ctx context.Context, text string, postID string, authorID string, parentCommentID *string) (*Comment, error) {
	panic(fmt.Errorf("not implemented: CreateComment - createComment"))
}

// EditPost is the resolver for the editPost field.
func (r *mutationResolver) EditPost(ctx context.Context, postID string, title *string, content *string, allowComments *bool) (*Post, error) {
	panic(fmt.Errorf("not implemented: EditPost - editPost"))
}
