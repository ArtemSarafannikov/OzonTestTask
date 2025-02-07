package service

import (
	"OzonTestTask/internal/models"
	"OzonTestTask/internal/repository"
	"context"
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
