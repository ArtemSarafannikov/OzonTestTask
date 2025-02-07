package repository

import (
	"context"
	"github.com/ArtemSarafannikov/OzonTestTask/internal/models"
)

type Repository interface {
	GetPosts(ctx context.Context, limit, offset int) ([]*models.Post, error)
	//GetPostById(ctx context.Context, id string) (*models.Post, error)
	//CreatePost(ctx context.Context, post *models.Post) (*models.Post, error)
	//GetComments(ctx context.Context, postId string, limit, offset int) ([]*models.Comment, error)
	//CreateComment(ctx context.Context, comment *models.Comment) (*models.Comment, error)
}
