package store

type Post struct {
	Title       string
	Description string
	UserID      int
	ID          int
}

type PostsStore struct{}

func (p *PostsStore) GetAll() error {
	return nil
}
