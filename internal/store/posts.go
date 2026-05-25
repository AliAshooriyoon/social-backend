package store

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Post struct {
	Title       string
	Description string
	UserID      int
	ID          int
	CreatedAt   time.Time
}

type PostsStore struct {
	Pool *pgxpool.Pool
}

func (p *PostsStore) GetLasts(ctx context.Context, count int) (*[]Post, error) {
	query := `SELECT * FROM posts 
ORDER BY created_at DESC 
LIMIT $1;`
	rows, err := p.Pool.Query(ctx, query, count)
	if err != nil {
		return nil, err
	}
	var posts []Post
	for rows.Next() {
		var post Post
		err := rows.Scan(&post.ID, &post.Title, &post.Description, &post.UserID, &post.CreatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return &posts, nil
}

func (p *PostsStore) Create(ctx context.Context, title, description string, userID int) error {
	query := `INSERT INTO posts(title,description,user_id) VALUES ($1,$2,$3)`
	_, err := p.Pool.Exec(ctx, query, title, description, userID)
	if err != nil {
		return err
	}
	return nil
}
