package repository

import (
	"context"
	"github.com/ArtemSarafannikov/OzonTestTask/internal/models"
	"sync"
	"time"
)

type InMemoryRepository struct {
	users    sync.Map // map[string]*models.Post
	posts    sync.Map // map[string]*models.Comment
	comments sync.Map
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
