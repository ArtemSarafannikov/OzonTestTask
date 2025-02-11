package service

import (
	"context"
	cstErrors "github.com/ArtemSarafannikov/OzonTestTask/internal/errors"
	"github.com/ArtemSarafannikov/OzonTestTask/internal/models"
	"github.com/ArtemSarafannikov/OzonTestTask/internal/repository"
	"github.com/ArtemSarafannikov/OzonTestTask/internal/utils"
	"log/slog"
)

type CommentService struct {
	repo repository.Repository
	log  *slog.Logger
}

func NewCommentService(log *slog.Logger, repo repository.Repository) *CommentService {
	return &CommentService{log: log, repo: repo}
}

func (s *CommentService) GetComments(ctx context.Context, postID string, authorID *string, limit, offset int) ([]*models.Comment, error) {
	const op = "CommentService.GetComments"
	log := s.log.With(
		slog.String("op", op),
	)

	var comments []*models.Comment
	var err error

	if authorID == nil {
		comments, err = s.repo.GetCommentsByPostID(ctx, postID, limit, offset)
	} else {
		comments, err = s.repo.GetCommentsByPostAuthorID(ctx, postID, *authorID, limit, offset)
	}
	if err != nil && !cstErrors.IsCustomError(err) {
		log.Error(err.Error())
	}
	return comments, err
}

func (s *CommentService) GetCommentByID(ctx context.Context, commentID string) (*models.Comment, error) {
	const op = "CommentService.GetCommentByID"
	log := s.log.With(
		slog.String("op", op),
	)

	comment, err := s.repo.GetCommentByID(ctx, commentID)
	if err != nil && !cstErrors.IsCustomError(err) {
		log.Error(err.Error())
	}
	return comment, err
}

func (s *CommentService) GetReplies(ctx context.Context, commentID string, limit, offset int) ([]*models.Comment, error) {
	const op = "CommentService.GetReplies"
	log := s.log.With(
		slog.String("op", op),
	)

	comments, err := s.repo.GetCommentsByCommentID(ctx, commentID, limit, offset)
	if err != nil && !cstErrors.IsCustomError(err) {
		log.Error(err.Error())
	}
	return comments, err
}

func (s *CommentService) CreateComment(ctx context.Context, text, postID string, parentID *string) (*models.Comment, error) {
	const op = "CommentService.CreateComment"
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
	if !post.AllowComments {
		return nil, cstErrors.PermissionDeniedError
	}

	// Needed????
	if parentID != nil {
		parent, err := s.repo.GetCommentByID(ctx, *parentID)
		if err != nil {
			if !cstErrors.IsCustomError(err) {
				log.Error(err.Error())
			}
			return nil, err
		}
		if parent.PostID != post.ID {
			return nil, cstErrors.CommentNotInPostError
		}
	}
	userID := utils.UserIDFromContext(ctx)
	comment := &models.Comment{
		PostID:   postID,
		ParentID: parentID,
		AuthorID: userID,
		Text:     text,
	}

	comment, err = s.repo.CreateComment(ctx, comment)
	if err != nil {
		if !cstErrors.IsCustomError(err) {
			log.Error(err.Error())
		}
		return nil, err
	}
	go s.repo.FixLastActivity(ctx, userID)
	return comment, err
}
