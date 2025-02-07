package repository

import (
	"OzonTestTask/internal/models"
	"context"
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
	return inm
}

func (r *InMemoryRepository) GetPosts(ctx context.Context, limit, offset int) ([]*models.Post, error) {
	var posts []*models.Post

	counter := 0
	end := offset + limit
	r.posts.Range(func(key, value interface{}) bool {
		if counter >= offset {
			post := value.(*models.Post)
			posts = append(posts, post)
		}
		if counter >= end {
			return false
		}
		counter++
		return true
	})
	return posts, nil
}
