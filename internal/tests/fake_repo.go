package tests

import (
	"context"
	"github.com/ArtemSarafannikov/OzonTestTask/internal/models"
)

// FakeRepository реализует интерфейс Repository для тестирования.
// Для каждого метода определён function-поле (с окончанием Fn), которое позволяет
// задавать необходимое поведение в тестах.
type FakeRepository struct {
	// Методы для работы с постами:
	GetPostsFn           func(ctx context.Context, limit, offset int) ([]*models.Post, error)
	GetPostByIDFn        func(ctx context.Context, id string) (*models.Post, error)
	GetPostsByAuthorIDFn func(ctx context.Context, authorID string, limit, offset int) ([]*models.Post, error)
	CreatePostFn         func(ctx context.Context, post *models.Post) (*models.Post, error)
	UpdatePostFn         func(ctx context.Context, post *models.Post)

	// Методы для работы с комментариями:
	GetCommentsByPostIDFn       func(ctx context.Context, postID string, limit, offset int) ([]*models.Comment, error)
	GetCommentsByPostAuthorIDFn func(ctx context.Context, postID string, authorID string, limit, offset int) ([]*models.Comment, error)
	GetCommentByIDFn            func(ctx context.Context, id string) (*models.Comment, error)
	GetCommentsByCommentIDFn    func(ctx context.Context, commentID string, limit, offset int) ([]*models.Comment, error)
	CreateCommentFn             func(ctx context.Context, comment *models.Comment) (*models.Comment, error)

	// Методы для работы с пользователями:
	CreateUserFn      func(ctx context.Context, user *models.User) (*models.User, error)
	GetUserByLoginFn  func(ctx context.Context, login string) (*models.User, error)
	GetUserByIDFn     func(ctx context.Context, id string) (*models.User, error)
	FixLastActivityFn func(ctx context.Context, id string) error

	// Методы для dataloader:
	GetUsersByIDsFn          func(ctx context.Context, ids []string) ([]*models.User, error)
	GetPostsByIDsFn          func(ctx context.Context, ids []string) ([]*models.Post, error)
	GetCommentsByIDsFn       func(ctx context.Context, ids []string) ([]*models.Comment, error)
	GetCommentsByPostIDsFn   func(ctx context.Context, ids []string) ([]*models.Comment, error)
	GetCommentsByAuthorIDsFn func(ctx context.Context, ids []string) ([]*models.Comment, error)
	GetCommentsByParentIDsFn func(ctx context.Context, ids []string) ([]*models.Comment, error)
}

// Реализация методов для работы с постами:
func (f *FakeRepository) GetPosts(ctx context.Context, limit, offset int) ([]*models.Post, error) {
	if f.GetPostsFn != nil {
		return f.GetPostsFn(ctx, limit, offset)
	}
	return nil, nil
}

func (f *FakeRepository) GetPostByID(ctx context.Context, id string) (*models.Post, error) {
	if f.GetPostByIDFn != nil {
		return f.GetPostByIDFn(ctx, id)
	}
	return nil, nil
}

func (f *FakeRepository) GetPostsByAuthorID(ctx context.Context, authorID string, limit, offset int) ([]*models.Post, error) {
	if f.GetPostsByAuthorIDFn != nil {
		return f.GetPostsByAuthorIDFn(ctx, authorID, limit, offset)
	}
	return nil, nil
}

func (f *FakeRepository) CreatePost(ctx context.Context, post *models.Post) (*models.Post, error) {
	if f.CreatePostFn != nil {
		return f.CreatePostFn(ctx, post)
	}
	return nil, nil
}

func (f *FakeRepository) UpdatePost(ctx context.Context, post *models.Post) {
	if f.UpdatePostFn != nil {
		f.UpdatePostFn(ctx, post)
	}
}

// Реализация методов для работы с комментариями:
func (f *FakeRepository) GetCommentsByPostID(ctx context.Context, postID string, limit, offset int) ([]*models.Comment, error) {
	if f.GetCommentsByPostIDFn != nil {
		return f.GetCommentsByPostIDFn(ctx, postID, limit, offset)
	}
	return nil, nil
}

func (f *FakeRepository) GetCommentsByPostAuthorID(ctx context.Context, postID string, authorID string, limit, offset int) ([]*models.Comment, error) {
	if f.GetCommentsByPostAuthorIDFn != nil {
		return f.GetCommentsByPostAuthorIDFn(ctx, postID, authorID, limit, offset)
	}
	return nil, nil
}

func (f *FakeRepository) GetCommentByID(ctx context.Context, id string) (*models.Comment, error) {
	if f.GetCommentByIDFn != nil {
		return f.GetCommentByIDFn(ctx, id)
	}
	return nil, nil
}

func (f *FakeRepository) GetCommentsByCommentID(ctx context.Context, commentID string, limit, offset int) ([]*models.Comment, error) {
	if f.GetCommentsByCommentIDFn != nil {
		return f.GetCommentsByCommentIDFn(ctx, commentID, limit, offset)
	}
	return nil, nil
}

func (f *FakeRepository) CreateComment(ctx context.Context, comment *models.Comment) (*models.Comment, error) {
	if f.CreateCommentFn != nil {
		return f.CreateCommentFn(ctx, comment)
	}
	return nil, nil
}

// Реализация методов для работы с пользователями:
func (f *FakeRepository) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	if f.CreateUserFn != nil {
		return f.CreateUserFn(ctx, user)
	}
	return nil, nil
}

func (f *FakeRepository) GetUserByLogin(ctx context.Context, login string) (*models.User, error) {
	if f.GetUserByLoginFn != nil {
		return f.GetUserByLoginFn(ctx, login)
	}
	return nil, nil
}

func (f *FakeRepository) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	if f.GetUserByIDFn != nil {
		return f.GetUserByIDFn(ctx, id)
	}
	return nil, nil
}

func (f *FakeRepository) FixLastActivity(ctx context.Context, id string) error {
	if f.FixLastActivityFn != nil {
		return f.FixLastActivityFn(ctx, id)
	}
	return nil
}

// Реализация методов для dataloader:
func (f *FakeRepository) GetUsersByIDs(ctx context.Context, ids []string) ([]*models.User, error) {
	if f.GetUsersByIDsFn != nil {
		return f.GetUsersByIDsFn(ctx, ids)
	}
	return nil, nil
}

func (f *FakeRepository) GetPostsByIDs(ctx context.Context, ids []string) ([]*models.Post, error) {
	if f.GetPostsByIDsFn != nil {
		return f.GetPostsByIDsFn(ctx, ids)
	}
	return nil, nil
}

func (f *FakeRepository) GetCommentsByIDs(ctx context.Context, ids []string) ([]*models.Comment, error) {
	if f.GetCommentsByIDsFn != nil {
		return f.GetCommentsByIDsFn(ctx, ids)
	}
	return nil, nil
}

func (f *FakeRepository) GetCommentsByPostIDs(ctx context.Context, ids []string) ([]*models.Comment, error) {
	if f.GetCommentsByPostIDsFn != nil {
		return f.GetCommentsByPostIDsFn(ctx, ids)
	}
	return nil, nil
}

func (f *FakeRepository) GetCommentsByAuthorIDs(ctx context.Context, ids []string) ([]*models.Comment, error) {
	if f.GetCommentsByAuthorIDsFn != nil {
		return f.GetCommentsByAuthorIDsFn(ctx, ids)
	}
	return nil, nil
}

func (f *FakeRepository) GetCommentsByParentIDs(ctx context.Context, ids []string) ([]*models.Comment, error) {
	if f.GetCommentsByParentIDsFn != nil {
		return f.GetCommentsByParentIDsFn(ctx, ids)
	}
	return nil, nil
}
