package repository

import (
	"context"
	"database/sql"
	"fmt"
	cstErrors "github.com/ArtemSarafannikov/OzonTestTask/internal/errors"
	"github.com/ArtemSarafannikov/OzonTestTask/internal/models"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"strings"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository() (*PostgresRepository, error) {
	// TODO: add config params
	// TODO: add logger for all methods
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/post_db?sslmode=disable")
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresRepository{db: db}, nil
}

func (p *PostgresRepository) GetParamsString(params []string) string {
	sb := strings.Builder{}
	length := len(params)
	for i := 1; i < length; i++ {
		sb.WriteString(fmt.Sprintf("$%d, ", i))
	}
	if length > 0 {
		sb.WriteString(fmt.Sprintf("$%d", length))
	}
	return sb.String()
}

func (p *PostgresRepository) GetPosts(ctx context.Context, limit, offset int) ([]*models.Post, error) {
	const op = "postgres.GetPosts"
	const query = `SELECT id, author_id, title, content, allowed_comments, edited_at, created_at
       FROM posts
       ORDER BY created_at DESC
       LIMIT $1 OFFSET $2`

	var err error
	rows, err := p.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		var (
			post models.Post
		)
		if err = rows.Scan(&post.ID,
			&post.AuthorID,
			&post.Title,
			&post.Content,
			&post.AllowComments,
			&post.EditedAt,
			&post.CreatedAt); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		posts = append(posts, &post)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return posts, nil
}

func (p *PostgresRepository) GetPostByID(ctx context.Context, id string) (*models.Post, error) {
	const op = "postgres.GetPostByID"
	const query = `SELECT id, author_id, title, content, allowed_comments, edited_at, created_at
       FROM posts
       WHERE id=$1`

	var post models.Post

	row := p.db.QueryRowContext(ctx, query, id)
	if row.Err() != nil {
		return nil, fmt.Errorf("%s: %w", op, row.Err())
	}
	if err := row.Scan(&post.ID,
		&post.AuthorID,
		&post.Title,
		&post.Content,
		&post.AllowComments,
		&post.EditedAt,
		&post.CreatedAt); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &post, nil
}

func (p *PostgresRepository) GetPostsByAuthorID(ctx context.Context, authorID string, limit, offset int) ([]*models.Post, error) {
	const op = "postgres.GetPostsByAuthorID"
	const query = `SELECT id, author_id, title, content, allowed_comments, edited_at, created_at
       FROM posts
       WHERE author_id=$1
       ORDER BY created_at DESC
       LIMIT $2 OFFSET $3`

	var err error
	rows, err := p.db.QueryContext(ctx, query, authorID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		var (
			post models.Post
		)
		if err = rows.Scan(&post.ID,
			&post.AuthorID,
			&post.Title,
			&post.Content,
			&post.AllowComments,
			&post.EditedAt,
			&post.CreatedAt); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		posts = append(posts, &post)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return posts, nil
}

func (p *PostgresRepository) CreatePost(ctx context.Context, post *models.Post) (*models.Post, error) {
	const op = "postgres.CreatePost"
	const query = `INSERT INTO posts (author_id, title, content, allowed_comments)
					VALUES ($1, $2, $3, $4)
					RETURNING id, created_at`

	row := p.db.QueryRowContext(ctx, query, post.AuthorID, post.Title)
	if row.Err() != nil {
		return nil, fmt.Errorf("%s: %w", op, row.Err())
	}
	if err := row.Scan(&post.ID,
		&post.CreatedAt); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return post, nil
}

func (p *PostgresRepository) UpdatePost(ctx context.Context, post *models.Post) {
	const op = "postgres.UpdatePost"
	const query = `UPDATE posts
					SET title=$1, content=$2, allowed_comments=$3, edited_at=now()
					WHERE id=$4`

	stmt, err := p.db.PrepareContext(ctx, query)
	if err != nil {
		// TODO: add logger
		return
	}

	if _, err = stmt.ExecContext(ctx, post.Title, post.Content, post.AllowComments, post.ID); err != nil {
		// TODO: add logger
	}
}

func (p *PostgresRepository) GetCommentsByPostID(ctx context.Context, postID string, limit, offset int) ([]*models.Comment, error) {
	const op = "postgres.GetCommentsByPostID"
	const query = `SELECT id, post_id, parent_comment_id, author_id, text, created_at
       FROM comments
       WHERE post_id=$1 AND parent_comment_id IS NULL
       ORDER BY created_at DESC
       LIMIT $2 OFFSET $3`

	var err error
	rows, err := p.db.QueryContext(ctx, query, postID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var comments []*models.Comment
	for rows.Next() {
		var (
			comment models.Comment
		)
		if err = rows.Scan(&comment.ID,
			&comment.PostID,
			&comment.ParentID,
			&comment.AuthorID,
			&comment.Text,
			&comment.CreatedAt); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		comments = append(comments, &comment)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return comments, nil
}

func (p *PostgresRepository) GetCommentsByPostAuthorID(ctx context.Context, postID string, authorID string, limit, offset int) ([]*models.Comment, error) {
	const op = "postgres.GetCommentsByPostAuthorID"
	const query = `SELECT id, post_id, parent_comment_id, author_id, text, created_at
       FROM comments
       WHERE post_id=$1 AND author_id=$2 AND parent_comment_id IS NULL
       ORDER BY created_at DESC
       LIMIT $3 OFFSET $4`

	var err error
	rows, err := p.db.QueryContext(ctx, query, postID, authorID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var comments []*models.Comment
	for rows.Next() {
		var (
			comment models.Comment
		)
		if err = rows.Scan(&comment.ID,
			&comment.PostID,
			&comment.ParentID,
			&comment.AuthorID,
			&comment.Text,
			&comment.CreatedAt); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		comments = append(comments, &comment)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return comments, nil
}

func (p *PostgresRepository) GetCommentByID(ctx context.Context, id string) (*models.Comment, error) {
	const op = "postgres.GetCommentByID"
	const query = `SELECT id, post_id, parent_comment_id, author_id, text, created_at
       FROM comments
       WHERE id=$1`

	var comment models.Comment

	row := p.db.QueryRowContext(ctx, query, id)
	if row.Err() != nil {
		return nil, fmt.Errorf("%s: %w", op, row.Err())
	}
	if err := row.Scan(&comment.ID,
		&comment.PostID,
		&comment.ParentID,
		&comment.AuthorID,
		&comment.Text,
		&comment.CreatedAt); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &comment, nil
}

func (p *PostgresRepository) GetCommentsByCommentID(ctx context.Context, commentID string, limit, offset int) ([]*models.Comment, error) {
	const op = "postgres.GetCommentsByCommentID"
	const query = `SELECT id, post_id, parent_comment_id, author_id, text, created_at
       FROM comments
       WHERE parent_comment_id=$1
       ORDER BY created_at DESC
       LIMIT $2 OFFSET $3`

	var err error
	rows, err := p.db.QueryContext(ctx, query, commentID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var comments []*models.Comment
	for rows.Next() {
		var (
			comment models.Comment
		)
		if err = rows.Scan(&comment.ID,
			&comment.PostID,
			&comment.ParentID,
			&comment.AuthorID,
			&comment.Text,
			&comment.CreatedAt); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		comments = append(comments, &comment)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return comments, nil
}

func (p *PostgresRepository) CreateComment(ctx context.Context, comment *models.Comment) (*models.Comment, error) {
	const op = "postgres.CreateComment"
	const query = `INSERT INTO comments (post_id, parent_comment_id, author_id, text)
					VALUES ($1, $2, $3, $4)
					RETURNING id, created_at`

	row := p.db.QueryRowContext(ctx, query, comment.PostID, comment.ParentID, comment.AuthorID, comment.Text)
	if row.Err() != nil {
		return nil, fmt.Errorf("%s: %w", op, row.Err())
	}
	if err := row.Scan(&comment.ID,
		&comment.CreatedAt); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return comment, nil
}

func (p *PostgresRepository) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	const op = "postgres.CreateUser"
	const query = `INSERT INTO users (login, password, last_activity)
					VALUES ($1, $2, $3)
					RETURNING id, created_at`

	row := p.db.QueryRowContext(ctx, query, user.Username, user.Password, user.LastActivity)
	if err := row.Err(); err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return nil, cstErrors.UserAlreadyExistsError
		}
		return nil, fmt.Errorf("%s: %w", op, row.Err())
	}
	if err := row.Scan(&user.ID,
		&user.CreatedAt); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return user, nil
}

func (p *PostgresRepository) GetUserByLogin(ctx context.Context, login string) (*models.User, error) {
	const op = "postgres.GetUserByLogin"
	const query = `SELECT id, login, password, last_activity, created_at
       FROM users
       WHERE login=$1`

	var user models.User

	row := p.db.QueryRowContext(ctx, query, login)
	if row.Err() != nil {
		return nil, fmt.Errorf("%s: %w", op, row.Err())
	}
	if err := row.Scan(&user.ID,
		&user.Username,
		&user.Password,
		&user.LastActivity,
		&user.CreatedAt); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &user, nil
}

func (p *PostgresRepository) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	const op = "postgres.GetUserByID"
	const query = `SELECT id, login, password, last_activity, created_at
       FROM users
       WHERE id=$1`

	var user models.User

	row := p.db.QueryRowContext(ctx, query, id)
	if row.Err() != nil {
		return nil, fmt.Errorf("%s: %w", op, row.Err())
	}
	if err := row.Scan(&user.ID,
		&user.Username,
		&user.Password,
		&user.LastActivity,
		&user.CreatedAt); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &user, nil
}

func (p *PostgresRepository) FixLastActivity(ctx context.Context, id string) error {
	const op = "postgres.FixLastActivity"
	const query = `UPDATE users
					SET last_activity=now()
					WHERE id=$1`

	stmt, err := p.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	if _, err = stmt.ExecContext(ctx, id); err != nil {
		return err
	}
	return nil
}

func (p *PostgresRepository) GetUsersByIDs(ctx context.Context, ids []string) ([]*models.User, error) {
	const op = "postgres.GetUsersByIDs"

	params := p.GetParamsString(ids)

	query := fmt.Sprintf(`SELECT id, login, password, last_activity, created_at
       FROM users
       WHERE id IN (%s)
       ORDER BY created_at DESC`, params)

	var err error
	rows, err := p.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var user models.User
		if err = rows.Scan(&user.ID,
			&user.Username,
			&user.Password,
			&user.LastActivity,
			&user.CreatedAt); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return users, nil
}

func (p *PostgresRepository) GetPostsByIDs(ctx context.Context, ids []string) ([]*models.Post, error) {
	const op = "postgres.GetPostsByIDs"

	params := p.GetParamsString(ids)

	query := fmt.Sprintf(`SELECT id, author_id, title, content, allowed_comments, edited_at, created_at
       FROM posts
       WHERE id IN (%s)
       ORDER BY created_at DESC`, params)

	var err error
	rows, err := p.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		var post models.Post
		if err = rows.Scan(&post.ID,
			&post.AuthorID,
			&post.Title,
			&post.Content,
			&post.AllowComments,
			&post.EditedAt,
			&post.CreatedAt); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		posts = append(posts, &post)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return posts, nil
}

func (p *PostgresRepository) GetCommentsByIDs(ctx context.Context, ids []string) ([]*models.Comment, error) {
	const op = "postgres.GetCommentsByIDs"

	params := p.GetParamsString(ids)

	query := fmt.Sprintf(`SELECT id, post_id, parent_comment_id, author_id, text, created_at
       FROM comments
       WHERE id IN (%s)
       ORDER BY created_at DESC`, params)

	var err error
	rows, err := p.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var comments []*models.Comment
	for rows.Next() {
		var comment models.Comment
		if err = rows.Scan(&comment.ID,
			&comment.PostID,
			&comment.ParentID,
			&comment.AuthorID,
			&comment.Text,
			&comment.CreatedAt); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		comments = append(comments, &comment)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return comments, nil
}

func (p *PostgresRepository) GetCommentsByPostIDs(ctx context.Context, ids []string) ([]*models.Comment, error) {
	const op = "postgres.GetCommentsByPostIDs"

	params := p.GetParamsString(ids)

	query := fmt.Sprintf(`SELECT id, post_id, parent_comment_id, author_id, text, created_at
       FROM comments
       WHERE post_id IN (%s) AND parent_comment_id IS NULL
       ORDER BY created_at DESC`, params)

	var err error
	rows, err := p.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var comments []*models.Comment
	for rows.Next() {
		var comment models.Comment
		if err = rows.Scan(&comment.ID,
			&comment.PostID,
			&comment.ParentID,
			&comment.AuthorID,
			&comment.Text,
			&comment.CreatedAt); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		comments = append(comments, &comment)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return comments, nil
}

func (p *PostgresRepository) GetCommentsByAuthorIDs(ctx context.Context, ids []string) ([]*models.Comment, error) {
	const op = "postgres.GetCommentsByAuthorIDs"

	params := p.GetParamsString(ids)

	query := fmt.Sprintf(`SELECT id, post_id, parent_comment_id, author_id, text, created_at
       FROM comments
       WHERE author_id IN (%s)
       ORDER BY created_at DESC`, params)

	var err error
	rows, err := p.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var comments []*models.Comment
	for rows.Next() {
		var comment models.Comment
		if err = rows.Scan(&comment.ID,
			&comment.PostID,
			&comment.ParentID,
			&comment.AuthorID,
			&comment.Text,
			&comment.CreatedAt); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		comments = append(comments, &comment)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return comments, nil
}

func (p *PostgresRepository) GetCommentsByParentIDs(ctx context.Context, ids []string) ([]*models.Comment, error) {
	const op = "postgres.GetCommentsByParentIDs"

	params := p.GetParamsString(ids)

	query := fmt.Sprintf(`SELECT id, post_id, parent_comment_id, author_id, text, created_at
       FROM comments
       WHERE parent_comment_id IN (%s)
       ORDER BY created_at DESC`, params)

	var err error
	rows, err := p.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var comments []*models.Comment
	for rows.Next() {
		var comment models.Comment
		if err = rows.Scan(&comment.ID,
			&comment.PostID,
			&comment.ParentID,
			&comment.AuthorID,
			&comment.Text,
			&comment.CreatedAt); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		comments = append(comments, &comment)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return comments, nil
}
