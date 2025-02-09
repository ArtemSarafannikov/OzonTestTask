package dataloaders

import (
	"context"
	"github.com/ArtemSarafannikov/OzonTestTask/internal/models"
	"github.com/ArtemSarafannikov/OzonTestTask/internal/repository"
	"time"
)

type DataLoaders struct {
	UserLoader              *UserLoader
	PostLoader              *PostLoader
	CommentLoader           *CommentLoader
	CommentByPostIDLoader   *CommentByPostIDLoader
	CommentByAuthorIDLoader *CommentByAuthorIDLoader
	CommentByParentIDLoader *CommentByParentIDLoader
}

func NewDataLoaders(repo repository.Repository) *DataLoaders {
	userConfig := UserLoaderConfig{
		Wait:     1 * time.Millisecond,
		MaxBatch: 100,
		Fetch: func(keys []string) ([]*models.User, []error) {
			users, err := repo.GetUsersByIDs(context.Background(), keys)
			if err != nil {
				errs := make([]error, len(keys))
				for i := range errs {
					errs[i] = err
				}
				return nil, errs
			}
			return users, nil
		},
	}
	userLoader := NewUserLoader(userConfig)

	postConfig := PostLoaderConfig{
		Wait:     1 * time.Millisecond,
		MaxBatch: 100,
		Fetch: func(keys []string) ([]*models.Post, []error) {
			posts, err := repo.GetPostsByIDs(context.Background(), keys)
			if err != nil {
				errs := make([]error, len(keys))
				for i := range errs {
					errs[i] = err
				}
				return nil, errs
			}
			return posts, nil
		},
	}
	postLoader := NewPostLoader(postConfig)

	commentConfig := CommentLoaderConfig{
		Wait:     1 * time.Millisecond,
		MaxBatch: 100,
		Fetch: func(keys []string) ([]*models.Comment, []error) {
			comments, err := repo.GetCommentsByIDs(context.Background(), keys)
			if err != nil {
				errs := make([]error, len(keys))
				for i := range errs {
					errs[i] = err
				}
				return nil, errs
			}
			return comments, nil
		},
	}
	commentLoader := NewCommentLoader(commentConfig)

	commentByPostConfig := CommentByPostIDLoaderConfig{
		Wait:     1 * time.Millisecond,
		MaxBatch: 100,
		Fetch: func(keys []string) ([][]*models.Comment, []error) {
			comments, err := repo.GetCommentsByPostIDs(context.Background(), keys)
			if err != nil {
				errs := make([]error, len(keys))
				for i := range errs {
					errs[i] = err
				}
				return nil, errs
			}
			commentsByPost := make(map[string][]*models.Comment)
			for _, comment := range comments {
				commentsByPost[comment.PostID] = append(commentsByPost[comment.PostID], comment)
			}

			result := make([][]*models.Comment, len(keys))
			for i, key := range keys {
				result[i] = commentsByPost[key]
			}
			return result, nil
		},
	}
	commentByPostLoader := NewCommentByPostIDLoader(commentByPostConfig)

	commentByAuthorConfig := CommentByAuthorIDLoaderConfig{
		Wait:     1 * time.Millisecond,
		MaxBatch: 100,
		Fetch: func(keys []string) ([][]*models.Comment, []error) {
			comments, err := repo.GetCommentsByAuthorIDs(context.Background(), keys)
			if err != nil {
				errs := make([]error, len(keys))
				for i := range errs {
					errs[i] = err
				}
				return nil, errs
			}
			commentsByAuthor := make(map[string][]*models.Comment)
			for _, comment := range comments {
				commentsByAuthor[comment.AuthorID] = append(commentsByAuthor[comment.AuthorID], comment)
			}

			result := make([][]*models.Comment, len(keys))
			for i, key := range keys {
				result[i] = commentsByAuthor[key]
			}
			return result, nil
		},
	}
	commentByAuthorLoader := NewCommentByAuthorIDLoader(commentByAuthorConfig)

	commentByParentConfig := CommentByParentIDLoaderConfig{
		Wait:     1 * time.Millisecond,
		MaxBatch: 100,
		Fetch: func(keys []string) ([][]*models.Comment, []error) {
			comments, err := repo.GetCommentsByParentIDs(context.Background(), keys)
			if err != nil {
				errs := make([]error, len(keys))
				for i := range errs {
					errs[i] = err
				}
				return nil, errs
			}
			commentsByParent := make(map[string][]*models.Comment)
			for _, comment := range comments {
				commentsByParent[comment.AuthorID] = append(commentsByParent[*comment.ParentID], comment)
			}

			result := make([][]*models.Comment, len(keys))
			for i, key := range keys {
				result[i] = commentsByParent[key]
			}
			return result, nil
		},
	}
	commentByParentLoader := NewCommentByParentIDLoader(commentByParentConfig)

	return &DataLoaders{
		UserLoader:              userLoader,
		PostLoader:              postLoader,
		CommentLoader:           commentLoader,
		CommentByPostIDLoader:   commentByPostLoader,
		CommentByAuthorIDLoader: commentByAuthorLoader,
		CommentByParentIDLoader: commentByParentLoader,
	}
}
