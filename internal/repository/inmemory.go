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
		AllowComments: false,
		AuthorID:      "0194e611-bfed-7a3c-9586-6c77012fbf7a",
		EditedAt:      &now,
		CreatedAt:     time.Now(),
	})
	inm.posts.Store("2", &models.Post{
		ID:            "2",
		Title:         "Second Post",
		Content:       "Test test test test",
		AllowComments: true,
		AuthorID:      "0194e611-bfed-7a3c-9586-6c77012fbf7a",
		EditedAt:      &now,
		CreatedAt:     time.Now(),
	})

	inm.users.Store("0194e611-bfed-7a3c-9586-6c77012fbf7a", &models.User{
		ID:           "0194e611-bfed-7a3c-9586-6c77012fbf7a",
		Username:     "admin",
		Password:     "$2a$14$f2PBTnj0UOY.P/gp4z/Qr.ESQe35DBW0vnMrObRo/5kw4zpr95S/6", // admin
		LastActivity: time.Now(),
		CreatedAt:    time.Now(),
	})

	inm.comments.Store("0194e628-e4f2-75e9-bc32-0601f9a6d4bc", &models.Comment{
		ID:        "0194e628-e4f2-75e9-bc32-0601f9a6d4bc",
		PostID:    "2",
		ParentID:  nil,
		AuthorID:  "0194e611-bfed-7a3c-9586-6c77012fbf7a",
		Text:      "under second post hahaha",
		CreatedAt: time.Now(),
	})
	inm.comments.Store("0194e628-e4f2-75e9-bc33-0601f9a6d4bc", &models.Comment{
		ID:        "0194e628-e4f2-75e9-bc33-0601f9a6d4bc",
		PostID:    "2",
		ParentID:  nil,
		AuthorID:  "0194e611-bfed-7a3c-9587-6c77012fbf7a",
		Text:      "from other user",
		CreatedAt: time.Now(),
	})
	parentID := "0194e628-e4f2-75e9-bc32-0601f9a6d4bc"
	inm.comments.Store("0194efaf-1171-7c9a-8ea1-f76a2dfd67b4", &models.Comment{
		ID:        "0194efaf-1171-7c9a-8ea1-f76a2dfd67b4",
		PostID:    "2",
		ParentID:  &parentID,
		AuthorID:  "0194e611-bfed-7a3c-9586-6c77012fbf7a",
		Text:      "Comment under comment",
		CreatedAt: time.Now(),
	})
	inm.comments.Store("0194efaf-1171-7b9a-8ea1-f76a2dfd67b4", &models.Comment{
		ID:        "0194efaf-1171-7b9a-8ea1-f76a2dfd67b4",
		PostID:    "2",
		ParentID:  &parentID,
		AuthorID:  "0194e611-bfed-7a3c-9586-6c77012fbf7a",
		Text:      "comment under comment 2",
		CreatedAt: time.Now(),
	})
	parentID2 := "0194efaf-1171-7b9a-8ea1-f76a2dfd67b4"
	inm.comments.Store("0004efaf-1171-7b9a-8ea1-f76a2dfd67b4", &models.Comment{
		ID:        "0004efaf-1171-7b9a-8ea1-f76a2dfd67b4",
		PostID:    "2",
		ParentID:  &parentID2,
		AuthorID:  "0194e611-bfed-7a3c-9586-6c77012fbf7a",
		Text:      "Comment under reply comment",
		CreatedAt: time.Now(),
	})

	//ticker := time.NewTicker(2 * time.Second)
	//go func() {
	//	for {
	//		<-ticker.C
	//		inm.users.Range(func(key, value interface{}) bool {
	//			fmt.Printf("Key: %v, Value: %+v\n", key, value.(*models.User))
	//			return true
	//		})
	//		fmt.Println()
	//	}
	//}()
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
		counter++
		return counter < last
	})
	return posts, nil
}

func (r *InMemoryRepository) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	if _, err := r.GetUserByLogin(ctx, user.Username); err == nil {
		return nil, cstErrors.UserAlreadyExistsError
	}

	newId, err := utils.GenerateNewID()
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
		return nil, cstErrors.NotFoundError
	}
	return user, nil
}

func (r *InMemoryRepository) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	user, exists := r.users.Load(id)
	if !exists {
		return nil, cstErrors.NotFoundError
	}
	return user.(*models.User), nil
}

func (r *InMemoryRepository) FixLastActivity(ctx context.Context, id string) error {
	user, err := r.GetUserByID(ctx, id)
	if err != nil {
		return err
	}

	user.LastActivity = time.Now()
	return nil
}

func (r *InMemoryRepository) GetPostByID(ctx context.Context, id string) (*models.Post, error) {
	post, exists := r.posts.Load(id)
	if !exists {
		return nil, cstErrors.NotFoundError
	}
	return post.(*models.Post), nil
}

func (r *InMemoryRepository) GetCommentsByPostID(ctx context.Context, postID string, limit, offset int) ([]*models.Comment, error) {
	var comments []*models.Comment

	counter := 0
	last := offset + limit
	r.comments.Range(func(key, value interface{}) bool {
		if counter >= offset {
			comment := value.(*models.Comment)
			if comment.PostID == postID && comment.ParentID == nil {
				counter++
				comments = append(comments, comment)
			}
		}
		return counter < last
	})
	return comments, nil
}

func (r *InMemoryRepository) GetCommentByID(ctx context.Context, id string) (*models.Comment, error) {
	comment, exists := r.comments.Load(id)
	if !exists {
		return nil, cstErrors.NotFoundError
	}
	return comment.(*models.Comment), nil
}

func (r *InMemoryRepository) GetCommentsByCommentID(ctx context.Context, commentID string, limit, offset int) ([]*models.Comment, error) {
	var comments []*models.Comment

	counter := 0
	last := offset + limit
	r.comments.Range(func(key, value interface{}) bool {
		if counter >= offset {
			comment := value.(*models.Comment)
			if comment.ParentID == nil {
				return true
			}
			if *comment.ParentID == commentID {
				counter++
				comments = append(comments, comment)
			}
		}
		return counter < last
	})
	return comments, nil
}

func (r *InMemoryRepository) CreatePost(ctx context.Context, post *models.Post) (*models.Post, error) {
	newId, err := utils.GenerateNewID()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	post.CreatedAt = now
	post.ID = newId
	r.posts.Store(newId, post)
	return post, nil
}

func (r *InMemoryRepository) UpdatePost(ctx context.Context, post *models.Post) {
	now := time.Now()
	post.EditedAt = &now
	r.posts.Store(post.ID, post)
}

func (r *InMemoryRepository) CreateComment(ctx context.Context, comment *models.Comment) (*models.Comment, error) {
	newId, err := utils.GenerateNewID()
	if err != nil {
		return nil, err
	}
	comment.ID = newId
	comment.CreatedAt = time.Now()
	r.comments.Store(newId, comment)
	return comment, nil
}

func (r *InMemoryRepository) GetUsersByIDs(ctx context.Context, ids []string) ([]*models.User, error) {
	users := make([]*models.User, 0, len(ids))

	for _, id := range ids {
		if user, err := r.GetUserByID(ctx, id); err == nil {
			users = append(users, user)
		}
	}
	return users, nil
}

func (r *InMemoryRepository) GetCommentsByIDs(ctx context.Context, ids []string) ([]*models.Comment, error) {
	comments := make([]*models.Comment, 0, len(ids))

	for _, id := range ids {
		if comment, err := r.GetCommentByID(ctx, id); err == nil {
			comments = append(comments, comment)
		}
	}
	return comments, nil
}

func (r *InMemoryRepository) GetCommentsByPostIDs(ctx context.Context, ids []string) ([]*models.Comment, error) {
	var comments []*models.Comment

	idsSet := map[string]struct{}{}
	for _, id := range ids {
		if _, err := r.GetPostByID(ctx, id); err == nil {
			idsSet[id] = struct{}{}
		}
	}
	r.comments.Range(func(key, value interface{}) bool {
		comment := value.(*models.Comment)
		if _, ok := idsSet[comment.PostID]; ok && comment.ParentID == nil {
			comments = append(comments, comment)
		}
		return true
	})
	return comments, nil
}

func (r *InMemoryRepository) GetCommentsByAuthorIDs(ctx context.Context, ids []string) ([]*models.Comment, error) {
	var comments []*models.Comment

	idsSet := map[string]struct{}{}
	for _, id := range ids {
		if _, err := r.GetUserByID(ctx, id); err == nil {
			idsSet[id] = struct{}{}
		}
	}
	r.comments.Range(func(key, value interface{}) bool {
		comment := value.(*models.Comment)
		if _, ok := idsSet[comment.AuthorID]; ok {
			comments = append(comments, comment)
		}
		return true
	})
	return comments, nil
}

func (r *InMemoryRepository) GetPostsByIDs(ctx context.Context, ids []string) ([]*models.Post, error) {
	posts := make([]*models.Post, 0, len(ids))

	for _, id := range ids {
		var post *models.Post
		post, err := r.GetPostByID(ctx, id)
		if err == nil {
			posts = append(posts, post)
		}
	}
	return posts, nil
}

func (r *InMemoryRepository) GetCommentsByParentIDs(ctx context.Context, ids []string) ([]*models.Comment, error) {
	var comments []*models.Comment

	idsSet := map[string]struct{}{}
	for _, id := range ids {
		if _, err := r.GetCommentByID(ctx, id); err == nil {
			idsSet[id] = struct{}{}
		}
	}
	r.comments.Range(func(key, value interface{}) bool {
		comment := value.(*models.Comment)
		if comment.ParentID != nil {
			if _, ok := idsSet[*comment.ParentID]; ok {
				comments = append(comments, comment)
			}
		}
		return true
	})
	return comments, nil
}

func (r *InMemoryRepository) GetPostsByAuthorID(ctx context.Context, authorID string, limit, offset int) ([]*models.Post, error) {
	var posts []*models.Post

	if _, err := r.GetUserByID(ctx, authorID); err != nil {
		return nil, cstErrors.NotFoundError
	}

	counter := 0
	last := offset + limit - 1
	r.posts.Range(func(key, value interface{}) bool {
		if counter >= offset {
			post := value.(*models.Post)
			if post.AuthorID == authorID {
				posts = append(posts, post)
				counter++
			}
		}
		return counter < last
	})
	return posts, nil
}

func (r *InMemoryRepository) GetCommentsByPostAuthorID(ctx context.Context, postID string, authorID string, limit, offset int) ([]*models.Comment, error) {
	var comments []*models.Comment

	if _, err := r.GetUserByID(ctx, authorID); err != nil {
		return nil, cstErrors.NotFoundError
	}

	counter := 0
	last := offset + limit
	r.comments.Range(func(key, value interface{}) bool {
		if counter >= offset {
			comment := value.(*models.Comment)
			if comment.PostID == postID && comment.AuthorID == authorID && comment.ParentID == nil {
				counter++
				comments = append(comments, comment)
			}
		}
		return counter < last
	})
	return comments, nil
}
