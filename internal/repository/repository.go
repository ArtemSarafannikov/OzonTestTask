package repository

import (
	"context"
	"github.com/ArtemSarafannikov/OzonTestTask/internal/models"
)

type Repository interface {
	GetPosts(ctx context.Context, limit, offset int) ([]*models.Post, error)
	GetPostById(ctx context.Context, id string) (*models.Post, error)
	CreatePost(ctx context.Context, post *models.Post) (*models.Post, error)
	UpdatePost(ctx context.Context, post *models.Post)

	GetCommentsByPostId(ctx context.Context, postID string, limit, offset int) ([]*models.Comment, error)
	GetCommentById(ctx context.Context, id string) (*models.Comment, error)
	GetCommentsByCommentId(ctx context.Context, commentID string, limit, offset int) ([]*models.Comment, error)
	CreateComment(ctx context.Context, comment *models.Comment) (*models.Comment, error)

	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	GetUserByLogin(ctx context.Context, login string) (*models.User, error)
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	FixLastActivity(ctx context.Context, id string) error

	// For dataloader
	GetUsersByIDs(ctx context.Context, ids []string) ([]*models.User, error)
	GetPostsByIDs(ctx context.Context, ids []string) ([]*models.Post, error)
	GetCommentsByIDs(ctx context.Context, ids []string) ([]*models.Comment, error)
	GetCommentsByPostIDs(ctx context.Context, ids []string) ([]*models.Comment, error)
	GetCommentsByAuthorIDs(ctx context.Context, ids []string) ([]*models.Comment, error)
	GetCommentsByParentIDs(ctx context.Context, ids []string) ([]*models.Comment, error)
}
