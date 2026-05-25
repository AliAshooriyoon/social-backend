// Package db: manage database
package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func New(addr string) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	pool, err := pgxpool.New(ctx, addr)
	if err != nil {
		return nil, err
	}
	pool.Begin(ctx)
	defer pool.Close()
	err = pool.Ping(ctx)
	if err != nil {
		return nil, err
	}
	pool.Config().MaxConns = 30
	pool.Config().MaxConnIdleTime = 30
	return pool, nil
}
