// Package store: storage of project
package store

import "context"

type Store struct {
	Posts interface {
		GetLasts(context.Context, int) (*[]Post, error)
		Create(context.Context, string, string, int) error
		Update(context.Context, string, string, int) error
	}
}
