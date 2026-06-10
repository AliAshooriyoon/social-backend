package main

import (
	"os"

	"project/internal/auth"
	"project/internal/db"
	"project/internal/store"

	"github.com/subosito/gotenv"
	"go.uber.org/zap"
)

type keyContextType string

var (
	userContextKey keyContextType = "user_auth"
	postContextKey keyContextType = "post"
)

//	@title			Swagger Example API
//	@version		1.0
//	@description	This is a sample server Petstore server.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @host		petstore.swagger.io
// @BasePath	/v2
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
	secretToken := os.Getenv("SECRET_TOKEN")
	authenticator := auth.NewAuth(secretToken, "aliash")
	cfg := config{
		addr: ":8080",
		auth: authenticator,
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
