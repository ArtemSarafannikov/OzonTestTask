package service

import "github.com/ArtemSarafannikov/OzonTestTask/internal/repository"

type CommentService struct {
	repo repository.Repository
}

func NewCommentService(repo repository.Repository) *CommentService {
	return &CommentService{repo: repo}
}
