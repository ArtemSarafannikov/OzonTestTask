package service

import (
	"context"
	cstErrors "github.com/ArtemSarafannikov/OzonTestTask/internal/errors"
	"github.com/ArtemSarafannikov/OzonTestTask/internal/models"
	"github.com/ArtemSarafannikov/OzonTestTask/internal/repository"
	"github.com/ArtemSarafannikov/OzonTestTask/internal/utils"
)

type PostService struct {
	repo repository.Repository
}

func NewPostService(repo repository.Repository) *PostService {
	return &PostService{repo: repo}
}

func (s *PostService) GetPosts(ctx context.Context, authorID *string, limit, offset int) ([]*models.Post, error) {
	if authorID == nil {
		return s.repo.GetPosts(ctx, limit, offset)
	}
	return s.repo.GetPostsByAuthorID(ctx, *authorID, limit, offset)
}

func (s *PostService) GetPostByID(ctx context.Context, id string) (*models.Post, error) {
	return s.repo.GetPostByID(ctx, id)
}

func (s *PostService) CreatePost(ctx context.Context, title, content string, allowComment bool) (*models.Post, error) {
	post := &models.Post{
		Title:         title,
		Content:       content,
		AllowComments: allowComment,
		AuthorID:      utils.UserIDFromContext(ctx),
	}
	return s.repo.CreatePost(ctx, post)
}

func (s *PostService) EditPost(ctx context.Context, postID string, title, content *string, allowComment *bool) (*models.Post, error) {
	post, err := s.repo.GetPostByID(ctx, postID)
	if err != nil {
		return nil, err
	}
	userID := utils.UserIDFromContext(ctx)
	if post.AuthorID != userID {
		return nil, cstErrors.PermissionDeniedError
	}

	if title != nil {
		post.Title = *title
	}
	if content != nil {
		post.Content = *content
	}
	if allowComment != nil {
		post.AllowComments = *allowComment
	}
	go s.repo.UpdatePost(ctx, post)
	return post, nil
}
