package store

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	ID        int       `json:"id"`
	UserName  string    `json:"user_name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Role      int       `json:"role"` // default 1 in db
	CreatedAt time.Time `json:"created_at"`
}

type UsersStore struct {
	Pool *pgxpool.Pool
}

func (u *UsersStore) GetByID(ctx context.Context, id int) (*User, error) {
	query := `
					SELECT email, user_name, password, created_at 
					FROM users 
					WHERE id = $1;
`
	var user User
	err := u.Pool.QueryRow(ctx, query, id).Scan(&user.Email, &user.UserName, &user.Password, user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
