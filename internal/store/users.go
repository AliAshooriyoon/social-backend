package store

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
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

func (u *UsersStore) Create(ctx context.Context, username, email, password string) error {
	query := `INSERT INTO users (email, userName, password)
VALUES (
    $1,
    $2,
    $3
);`
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = u.Pool.Exec(ctx, query, username, email, hashed)
	if err != nil {
		return err
	}
	return nil
}

func (u *UsersStore) Update(ctx context.Context, username, email, password string, userID int) error {
	query := `UPDATE users
SET
    email = $1,
    userName = $2,
    password = $3
		WHERE id = $4;`
	_, err := u.Pool.Exec(ctx, query, username, email, password, userID)
	if err != nil {
		return err
	}
	return nil
}

func (u *UsersStore) Delete(ctx context.Context, userID int) error {
	query := `DELETE FROM users
		WHERE id = $1;`
	_, err := u.Pool.Exec(ctx, query, userID)
	if err != nil {
		return err
	}

	return nil
}
