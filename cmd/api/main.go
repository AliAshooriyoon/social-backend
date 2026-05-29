package main

import (
	"os"

	"project/internal/db"
	"project/internal/store"

	"github.com/subosito/gotenv"
	"go.uber.org/zap"
)

func main() {
	gotenv.Load()
	dbAddr := os.Getenv("DATABASE")
	logger := zap.Must(zap.NewProduction()).Sugar()

	defer logger.Sync()

	db, err := db.New(dbAddr)
	if err != nil {
		logger.Fatal(err)
		return
	}
	cfg := config{
		addr: ":8080",
	}

	app := application{
		config: cfg,
		store: store.Store{
			Posts: &store.PostsStore{
				Pool: db,
			},
			Users: &store.UsersStore{
				Pool: db,
			},
		},
		logger: logger,
	}
	err = app.run()
	if err != nil {
		logger.Fatal(err)
	}
}
