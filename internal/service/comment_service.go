package service

import (
	"context"
	cstErrors "github.com/ArtemSarafannikov/OzonTestTask/internal/errors"
	"github.com/ArtemSarafannikov/OzonTestTask/internal/models"
	"github.com/ArtemSarafannikov/OzonTestTask/internal/repository"
	"github.com/ArtemSarafannikov/OzonTestTask/internal/utils"
)

type CommentService struct {
	repo repository.Repository
}

func NewCommentService(repo repository.Repository) *CommentService {
	return &CommentService{repo: repo}
}

func (s *CommentService) GetComments(ctx context.Context, postID string, authorID *string, limit, offset int) ([]*models.Comment, error) {
	if authorID == nil {
		return s.repo.GetCommentsByPostID(ctx, postID, limit, offset)
	}
	return s.repo.GetCommentsByPostAuthorID(ctx, postID, *authorID, limit, offset)
}

func (s *CommentService) GetCommentByID(ctx context.Context, commentID string) (*models.Comment, error) {
	return s.repo.GetCommentByID(ctx, commentID)
}

func (s *CommentService) GetReplies(ctx context.Context, commentID string, limit, offset int) ([]*models.Comment, error) {
	return s.repo.GetCommentsByCommentID(ctx, commentID, limit, offset)
}

func (s *CommentService) CreateComment(ctx context.Context, text, postID string, parentID *string) (*models.Comment, error) {
	post, err := s.repo.GetPostByID(ctx, postID)
	if err != nil {
		return nil, err
	}
	if !post.AllowComments {
		return nil, cstErrors.PermissionDeniedError
	}

	// Needed????
	if parentID != nil {
		parent, err := s.repo.GetCommentByID(ctx, *parentID)
		if err != nil {
			return nil, err
		}
		if parent.PostID != post.ID {
			return nil, cstErrors.CommentNotInPostError
		}
	}
	comment := &models.Comment{
		PostID:   postID,
		ParentID: parentID,
		AuthorID: utils.UserIDFromContext(ctx),
		Text:     text,
	}
	return s.repo.CreateComment(ctx, comment)
}
