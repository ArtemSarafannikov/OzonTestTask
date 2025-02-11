package service

import (
	"context"
	cstErrors "github.com/ArtemSarafannikov/OzonTestTask/internal/errors"
	"github.com/ArtemSarafannikov/OzonTestTask/internal/models"
	"github.com/ArtemSarafannikov/OzonTestTask/internal/repository"
	"github.com/ArtemSarafannikov/OzonTestTask/internal/utils"
	"log/slog"
)

type PostService struct {
	repo repository.Repository
	log  *slog.Logger
}

func NewPostService(log *slog.Logger, repo repository.Repository) *PostService {
	return &PostService{log: log, repo: repo}
}

func (s *PostService) GetPosts(ctx context.Context, authorID *string, limit, offset int) ([]*models.Post, error) {
	const op = "PostService.GetPosts"
	log := s.log.With(
		slog.String("op", op),
	)

	var posts []*models.Post
	var err error

	if authorID == nil {
		posts, err = s.repo.GetPosts(ctx, limit, offset)
	} else {
		posts, err = s.repo.GetPostsByAuthorID(ctx, *authorID, limit, offset)
	}
	if err != nil && !cstErrors.IsCustomError(err) {
		log.Error(err.Error())
	}
	return posts, err
}

func (s *PostService) GetPostByID(ctx context.Context, id string) (*models.Post, error) {
	const op = "PostService.GetPostByID"
	log := s.log.With(
		slog.String("op", op),
	)

	post, err := s.repo.GetPostByID(ctx, id)
	if err != nil && !cstErrors.IsCustomError(err) {
		log.Error(err.Error())
	}
	return post, err
}

func (s *PostService) CreatePost(ctx context.Context, title, content string, allowComment bool) (*models.Post, error) {
	const op = "PostService.CreatePost"
	log := s.log.With(
		slog.String("op", op),
	)

	userID := utils.UserIDFromContext(ctx)
	post := &models.Post{
		Title:         title,
		Content:       content,
		AllowComments: allowComment,
		AuthorID:      userID,
	}
	post, err := s.repo.CreatePost(ctx, post)
	if err != nil {
		if !cstErrors.IsCustomError(err) {
			log.Error(err.Error())
		}
		return nil, err
	}
	go s.repo.FixLastActivity(ctx, userID)
	return post, err
}

func (s *PostService) EditPost(ctx context.Context, postID string, title, content *string, allowComment *bool) (*models.Post, error) {
	const op = "PostService.EditPost"
	log := s.log.With(
		slog.String("op", op),
	)

	post, err := s.repo.GetPostByID(ctx, postID)
	if err != nil {
		if !cstErrors.IsCustomError(err) {
			log.Error(err.Error())
		}
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
	go s.repo.FixLastActivity(ctx, userID)
	return post, nil
}
