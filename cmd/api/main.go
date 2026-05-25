package main

import (
	"project/internal/db"
	"project/internal/store"
)

func main() {
	// get env!
	db, err := db.New("")
	if err != nil {
		return
	}
	app := application{
		config: config{
			addr: ":8080",
		},
		store: store.Store{
			Posts: &store.PostsStore{
				Pool: db,
			},
		},
	}

	app.run()
}
