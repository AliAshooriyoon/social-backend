package main

import (
	"os"

	"project/internal/db"
	"project/internal/store"

	"github.com/subosito/gotenv"
)

func main() {
	gotenv.Load()
	dbAddr := os.Getenv("DATABASE")
	db, err := db.New(dbAddr)
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
