package service

import (
	"context"
	"github.com/ArtemSarafannikov/OzonTestTask/internal/models"
	"github.com/ArtemSarafannikov/OzonTestTask/internal/repository"
)

type PostService struct {
	repo repository.Repository
}

func NewPostService(repo repository.Repository) *PostService {
	return &PostService{repo: repo}
}

func (s *PostService) GetPosts(ctx context.Context, limit, offset int) ([]*models.Post, error) {
	return s.repo.GetPosts(ctx, limit, offset)
}
