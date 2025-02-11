package service

import (
	"context"
	"errors"
	cstErrors "github.com/ArtemSarafannikov/OzonTestTask/internal/errors"
	"github.com/ArtemSarafannikov/OzonTestTask/internal/models"
	"github.com/ArtemSarafannikov/OzonTestTask/internal/testutils"
	"github.com/ArtemSarafannikov/OzonTestTask/internal/utils"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// Тест получения комментариев по postID без фильтра по authorID.
func TestCommentService_GetComments_NoAuthor(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	fakeRepo := &testutils.FakeRepository{
		GetCommentsByPostIDFn: func(ctx context.Context, postID string, limit, offset int) ([]*models.Comment, error) {
			return []*models.Comment{
				{
					ID:        "c1",
					PostID:    postID,
					AuthorID:  "author1",
					Text:      "Комментарий 1",
					CreatedAt: now,
				},
				{
					ID:        "c2",
					PostID:    postID,
					AuthorID:  "author2",
					Text:      "Комментарий 2",
					CreatedAt: now,
				},
			}, nil
		},
	}
	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)
	svc := NewCommentService(log, fakeRepo)
	comments, err := svc.GetComments(ctx, "post1", nil, 10, 0)
	require.NoError(t, err)
	require.Len(t, comments, 2)
	require.Equal(t, "c1", comments[0].ID)
}

// Тест получения комментариев по postID с фильтром по authorID.
func TestCommentService_GetComments_WithAuthor(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	fakeRepo := &testutils.FakeRepository{
		GetCommentsByPostAuthorIDFn: func(ctx context.Context, postID string, authorID string, limit, offset int) ([]*models.Comment, error) {
			if authorID == "a1" {
				return []*models.Comment{
					{
						ID:        "c3",
						PostID:    postID,
						AuthorID:  authorID,
						Text:      "Комментарий автора",
						CreatedAt: now,
					},
				}, nil
			}
			return nil, nil
		},
	}
	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)
	svc := NewCommentService(log, fakeRepo)
	authorID := "a1"
	comments, err := svc.GetComments(ctx, "post1", &authorID, 10, 0)
	require.NoError(t, err)
	require.Len(t, comments, 1)
	require.Equal(t, "c3", comments[0].ID)
}

// Тест получения комментария по ID (успешный случай).
func TestCommentService_GetCommentByID_Success(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	expected := &models.Comment{
		ID:        "c1",
		PostID:    "p1",
		AuthorID:  "author1",
		Text:      "Текст комментария",
		CreatedAt: now,
	}
	fakeRepo := &testutils.FakeRepository{
		GetCommentByIDFn: func(ctx context.Context, id string) (*models.Comment, error) {
			if id == "c1" {
				return expected, nil
			}
			return nil, errors.New("не найден")
		},
	}
	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)
	svc := NewCommentService(log, fakeRepo)
	comment, err := svc.GetCommentByID(ctx, "c1")
	require.NoError(t, err)
	require.Equal(t, expected, comment)
}

// Тест получения ответов (replies) для комментария.
func TestCommentService_GetReplies_Success(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	fakeRepo := &testutils.FakeRepository{
		GetCommentsByCommentIDFn: func(ctx context.Context, commentID string, limit, offset int) ([]*models.Comment, error) {
			return []*models.Comment{
				{
					ID:        "r1",
					PostID:    "p1",
					ParentID:  &commentID,
					AuthorID:  "author1",
					Text:      "Ответ 1",
					CreatedAt: now,
				},
				{
					ID:        "r2",
					PostID:    "p1",
					ParentID:  &commentID,
					AuthorID:  "author2",
					Text:      "Ответ 2",
					CreatedAt: now,
				},
			}, nil
		},
	}
	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)
	svc := NewCommentService(log, fakeRepo)
	replies, err := svc.GetReplies(ctx, "c1", 10, 0)
	require.NoError(t, err)
	require.Len(t, replies, 2)
	require.Equal(t, "r1", replies[0].ID)
}

// Тест успешного создания комментария.
func TestCommentService_CreateComment_Success(t *testing.T) {
	// Создаём контекст с userID
	ctx := context.WithValue(context.Background(), utils.UserIdCtxKey, "test_author")
	now := time.Now()

	fakeRepo := &testutils.FakeRepository{
		GetPostByIDFn: func(ctx context.Context, id string) (*models.Post, error) {
			// Имитируем, что пост найден и разрешает комментарии
			return &models.Post{
				ID:            id,
				Title:         "Test Post",
				Content:       "Test Content",
				AllowComments: true,
				CreatedAt:     now,
			}, nil
		},
		// Если передан parentID, возвращаем корректный родительский комментарий
		GetCommentByIDFn: func(ctx context.Context, id string) (*models.Comment, error) {
			return &models.Comment{
				ID:        id,
				PostID:    "post1",
				AuthorID:  "some_author",
				Text:      "Parent Comment",
				CreatedAt: now,
			}, nil
		},
		CreateCommentFn: func(ctx context.Context, comment *models.Comment) (*models.Comment, error) {
			// Присваиваем новый ID и, если время создания не задано, устанавливаем его
			comment.ID = "new_comment"
			if comment.CreatedAt.IsZero() {
				comment.CreatedAt = now
			}
			return comment, nil
		},
	}

	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)
	svc := NewCommentService(log, fakeRepo)
	parentID := "parent1"
	comment, err := svc.CreateComment(ctx, "New Comment", "post1", &parentID)
	require.NoError(t, err)
	require.NotNil(t, comment)
	require.Equal(t, "new_comment", comment.ID)
	require.Equal(t, "New Comment", comment.Text)
	require.Equal(t, "post1", comment.PostID)
	require.NotNil(t, comment.ParentID)
	require.Equal(t, parentID, *comment.ParentID)
	// Проверяем, что author_id установлен из контекста
	require.Equal(t, "test_author", comment.AuthorID)
}

// Тест создания комментария с возникновением ошибки.
func TestCommentService_CreateComment_Error(t *testing.T) {
	now := time.Now()
	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	t.Run("GetPostByID error", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), utils.UserIdCtxKey, "test_author")
		fakeRepo := &testutils.FakeRepository{
			GetPostByIDFn: func(ctx context.Context, id string) (*models.Post, error) {
				return nil, errors.New("post not found")
			},
		}
		svc := NewCommentService(log, fakeRepo)
		comment, err := svc.CreateComment(ctx, "Comment text", "post1", nil)
		require.Error(t, err)
		require.Nil(t, comment)
		require.EqualError(t, err, "post not found")
	})

	t.Run("Post does not allow comments", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), utils.UserIdCtxKey, "test_author")
		fakeRepo := &testutils.FakeRepository{
			GetPostByIDFn: func(ctx context.Context, id string) (*models.Post, error) {
				return &models.Post{
					ID:            id,
					Title:         "Test Post",
					Content:       "Test Content",
					AllowComments: false, // комментарии запрещены
					CreatedAt:     now,
				}, nil
			},
		}
		svc := NewCommentService(log, fakeRepo)
		comment, err := svc.CreateComment(ctx, "Comment text", "post1", nil)
		require.Error(t, err)
		require.Nil(t, comment)
		// Ожидаем PermissionDeniedError
		require.Equal(t, cstErrors.PermissionDeniedError, err)
	})

	t.Run("Parent comment GetCommentByID error", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), utils.UserIdCtxKey, "test_author")
		fakeRepo := &testutils.FakeRepository{
			GetPostByIDFn: func(ctx context.Context, id string) (*models.Post, error) {
				return &models.Post{
					ID:            id,
					Title:         "Test Post",
					Content:       "Test Content",
					AllowComments: true,
					CreatedAt:     now,
				}, nil
			},
			GetCommentByIDFn: func(ctx context.Context, id string) (*models.Comment, error) {
				return nil, errors.New("parent comment not found")
			},
		}
		svc := NewCommentService(log, fakeRepo)
		parentID := "invalid_parent"
		comment, err := svc.CreateComment(ctx, "Comment text", "post1", &parentID)
		require.Error(t, err)
		require.Nil(t, comment)
		require.EqualError(t, err, "parent comment not found")
	})

	t.Run("Parent comment post mismatch", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), utils.UserIdCtxKey, "test_author")
		fakeRepo := &testutils.FakeRepository{
			GetPostByIDFn: func(ctx context.Context, id string) (*models.Post, error) {
				return &models.Post{
					ID:            id,
					Title:         "Test Post",
					Content:       "Test Content",
					AllowComments: true,
					CreatedAt:     now,
				}, nil
			},
			GetCommentByIDFn: func(ctx context.Context, id string) (*models.Comment, error) {
				// Возвращаем комментарий, привязанный к другому посту
				return &models.Comment{
					ID:        id,
					PostID:    "different_post",
					AuthorID:  "some_author",
					Text:      "Parent Comment",
					CreatedAt: now,
				}, nil
			},
		}
		svc := NewCommentService(log, fakeRepo)
		parentID := "parent1"
		comment, err := svc.CreateComment(ctx, "Comment text", "post1", &parentID)
		require.Error(t, err)
		require.Nil(t, comment)
		// Ожидаем CommentNotInPostError
		require.Equal(t, cstErrors.CommentNotInPostError, err)
	})

	t.Run("CreateComment error", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), utils.UserIdCtxKey, "test_author")
		fakeRepo := &testutils.FakeRepository{
			GetPostByIDFn: func(ctx context.Context, id string) (*models.Post, error) {
				return &models.Post{
					ID:            id,
					Title:         "Test Post",
					Content:       "Test Content",
					AllowComments: true,
					CreatedAt:     now,
				}, nil
			},
			CreateCommentFn: func(ctx context.Context, comment *models.Comment) (*models.Comment, error) {
				return nil, errors.New("create comment failed")
			},
		}
		svc := NewCommentService(log, fakeRepo)
		comment, err := svc.CreateComment(ctx, "Comment text", "post1", nil)
		require.Error(t, err)
		require.Nil(t, comment)
		require.EqualError(t, err, "create comment failed")
	})
}
