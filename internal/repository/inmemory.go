package repository

import (
	"context"
	"github.com/ArtemSarafannikov/OzonTestTask/internal/errors"
	"github.com/ArtemSarafannikov/OzonTestTask/internal/models"
	"github.com/ArtemSarafannikov/OzonTestTask/internal/utils"
	"sync"
	"time"
)

type InMemoryRepository struct {
	users    sync.Map // map[string]*models.User
	posts    sync.Map // map[string]*models.Post
	comments sync.Map // map[string]*models.Comment
}

func NewInMemoryRepository() *InMemoryRepository {
	inm := &InMemoryRepository{}
	now := time.Now()
	inm.posts.Store("1", &models.Post{
		ID:            "1",
		Title:         "First Post",
		Content:       "This is the first post",
		AllowComments: true,
		AuthorID:      "1",
		EditedAt:      &now,
		CreatedAt:     time.Now(),
	})
	inm.posts.Store("2", &models.Post{
		ID:            "2",
		Title:         "Second Post",
		Content:       "Test test test test",
		AllowComments: false,
		AuthorID:      "1",
		EditedAt:      &now,
		CreatedAt:     time.Now(),
	})
	return inm
}

func (r *InMemoryRepository) GetPosts(ctx context.Context, limit, offset int) ([]*models.Post, error) {
	var posts []*models.Post

	counter := 0
	last := offset + limit - 1
	r.posts.Range(func(key, value interface{}) bool {
		if counter >= offset {
			post := value.(*models.Post)
			posts = append(posts, post)
		}
		if counter >= last {
			return false
		}
		counter++
		return true
	})
	return posts, nil
}

func (r *InMemoryRepository) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	if _, err := r.GetUserByLogin(ctx, user.Username); err == nil {
		return nil, cstErrors.UserAlreadyExistsError
	}

	newId, err := utils.GenerateNewId()
	if err != nil {
		return nil, err
	}

	now := time.Now()

	user.ID = newId
	user.CreatedAt = now
	user.LastActivity = now
	r.users.Store(newId, user)
	return user, nil
}

func (r *InMemoryRepository) GetUserByLogin(ctx context.Context, login string) (*models.User, error) {
	user := &models.User{}
	r.users.Range(func(key, value interface{}) bool {
		userValue := value.(*models.User)
		if userValue.Username == login {
			user = userValue
			return false
		}
		return true
	})
	if user.ID == "" {
		return nil, cstErrors.UserNotFoundError
	}
	return user, nil
}
