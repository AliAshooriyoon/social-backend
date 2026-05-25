// Package store: storage of project
package store

type Store struct {
	Posts interface {
		GetAll(int) error
	}
}
