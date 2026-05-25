package store

import "github.com/jackc/pgx/v5/pgxpool"

type Post struct {
	Title       string
	Description string
	UserID      int
	ID          int
}

type PostsStore struct {
	Pool *pgxpool.Pool
}

func (p *PostsStore) GetAll(count int) error {
	return nil
}
