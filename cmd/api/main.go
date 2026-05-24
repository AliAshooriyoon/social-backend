package main

import "project/internal/store"

func main() {
	app := application{
		config: config{
			addr: ":8080",
		},
		store: store.Store{
			Posts: &store.PostsStore{},
		},
	}
	app.run()
}
