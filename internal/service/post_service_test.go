package service

import (
	"context"
	"errors"
	cstErrors "github.com/ArtemSarafannikov/OzonTestTask/internal/errors"
	"github.com/ArtemSarafannikov/OzonTestTask/internal/models"
	"github.com/ArtemSarafannikov/OzonTestTask/internal/testutils"
	"github.com/ArtemSarafannikov/OzonTestTask/internal/utils"
	"github.com/stretchr/testify/require"
	"log/slog"
	"os"
	"testing"
	"time"
)

func TestPostService_GetPosts_NoAuthor(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	fakeRepo := &testutils.FakeRepository{
		GetPostsFn: func(ctx context.Context, limit, offset int) ([]*models.Post, error) {
			return []*models.Post{
				{ID: "p1", Title: "Post 1", Content: "Content 1", AllowComments: true, AuthorID: "a1", CreatedAt: now},
				{ID: "p2", Title: "Post 2", Content: "Content 2", AllowComments: true, AuthorID: "a2", CreatedAt: now},
			}, nil
		},
	}
	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)
	svc := NewPostService(log, fakeRepo)
	posts, err := svc.GetPosts(ctx, nil, 10, 0)
	require.NoError(t, err)
	require.Len(t, posts, 2)
	require.Equal(t, "p1", posts[0].ID)
}

func TestPostService_GetPosts_WithAuthor(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	fakeRepo := &testutils.FakeRepository{
		GetPostsByAuthorIDFn: func(ctx context.Context, authorID string, limit, offset int) ([]*models.Post, error) {
			if authorID == "a1" {
				return []*models.Post{
					{ID: "p1", Title: "Post 1", Content: "Content 1", AllowComments: true, AuthorID: "a1", CreatedAt: now},
				}, nil
			}
			return nil, nil
		},
	}
	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)
	svc := NewPostService(log, fakeRepo)
	author := "a1"
	posts, err := svc.GetPosts(ctx, &author, 10, 0)
	require.NoError(t, err)
	require.Len(t, posts, 1)
	require.Equal(t, "p1", posts[0].ID)
}

func TestPostService_GetPostByID(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	fakeRepo := &testutils.FakeRepository{
		GetPostByIDFn: func(ctx context.Context, id string) (*models.Post, error) {
			if id == "p1" {
				return &models.Post{
					ID: "p1", Title: "Post 1", Content: "Content 1", AllowComments: true, AuthorID: "a1", CreatedAt: now,
				}, nil
			}
			return nil, errors.New("post not found")
		},
	}
	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)
	svc := NewPostService(log, fakeRepo)
	post, err := svc.GetPostByID(ctx, "p1")
	require.NoError(t, err)
	require.NotNil(t, post)
	require.Equal(t, "p1", post.ID)
}

func TestPostService_CreatePost_Success(t *testing.T) {
	ctx := context.WithValue(context.Background(), utils.UserIdCtxKey, "test_author")
	now := time.Now()
	fakeRepo := &testutils.FakeRepository{
		CreatePostFn: func(ctx context.Context, post *models.Post) (*models.Post, error) {
			post.ID = "new_post"
			if post.CreatedAt.IsZero() {
				post.CreatedAt = now
			}
			return post, nil
		},
	}
	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)
	svc := NewPostService(log, fakeRepo)
	post, err := svc.CreatePost(ctx, "New Title", "New Content", true)
	require.NoError(t, err)
	require.NotNil(t, post)
	require.Equal(t, "New Title", post.Title)
	require.Equal(t, "New Content", post.Content)
	require.True(t, post.AllowComments)
	require.Equal(t, "test_author", post.AuthorID)
	require.Equal(t, "new_post", post.ID)
}

func TestPostService_EditPost_Success(t *testing.T) {
	ctx := context.WithValue(context.Background(), utils.UserIdCtxKey, "test_author")
	now := time.Now()
	initialPost := &models.Post{
		ID: "p1", Title: "Old Title", Content: "Old Content", AllowComments: true,
		AuthorID: "test_author", CreatedAt: now,
	}
	fakeRepo := &testutils.FakeRepository{
		GetPostByIDFn: func(ctx context.Context, id string) (*models.Post, error) {
			if id == "p1" {
				return initialPost, nil
			}
			return nil, errors.New("post not found")
		},
	}
	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)
	svc := NewPostService(log, fakeRepo)
	newTitle := "New Title"
	newContent := "New Content"
	newAllow := false
	editedPost, err := svc.EditPost(ctx, "p1", &newTitle, &newContent, &newAllow)
	require.NoError(t, err)
	require.NotNil(t, editedPost)
	require.Equal(t, newTitle, editedPost.Title)
	require.Equal(t, newContent, editedPost.Content)
	require.Equal(t, newAllow, editedPost.AllowComments)
}

func TestPostService_EditPost_GetPostError(t *testing.T) {
	ctx := context.WithValue(context.Background(), utils.UserIdCtxKey, "test_author")
	fakeRepo := &testutils.FakeRepository{
		GetPostByIDFn: func(ctx context.Context, id string) (*models.Post, error) {
			return nil, errors.New("post not found")
		},
	}
	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)
	svc := NewPostService(log, fakeRepo)
	newTitle := "New Title"
	editedPost, err := svc.EditPost(ctx, "p1", &newTitle, nil, nil)
	require.Error(t, err)
	require.Nil(t, editedPost)
	require.EqualError(t, err, "post not found")
}

func TestPostService_EditPost_PermissionDenied(t *testing.T) {
	ctx := context.WithValue(context.Background(), utils.UserIdCtxKey, "test_author")
	now := time.Now()
	initialPost := &models.Post{
		ID: "p1", Title: "Old Title", Content: "Old Content", AllowComments: true,
		AuthorID: "other_author", CreatedAt: now,
	}
	fakeRepo := &testutils.FakeRepository{
		GetPostByIDFn: func(ctx context.Context, id string) (*models.Post, error) {
			return initialPost, nil
		},
	}
	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)
	svc := NewPostService(log, fakeRepo)
	newTitle := "New Title"
	editedPost, err := svc.EditPost(ctx, "p1", &newTitle, nil, nil)
	require.Error(t, err)
	require.Nil(t, editedPost)
	require.Equal(t, cstErrors.PermissionDeniedError, err)
}
