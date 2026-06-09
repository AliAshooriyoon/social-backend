package main

import (
	"fmt"
	"net/http"
	"time"

	"project/internal/auth"
	"project/internal/store"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

type application struct {
	config config
	store  store.Store
	logger *zap.SugaredLogger
}

type config struct {
	addr string
	env  string
	auth auth.Authenticator
}

func (app *application) mount() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hi"))
	})
	docsURL := fmt.Sprintf("localhost%s/swagger/doc.json", app.config.addr)
	router.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(docsURL)))
	router.Route("/posts", func(r chi.Router) {
		r.Get("/", app.getPostsHandler)
		r.Post("/new-post", app.createPostHandler)
		r.Put("/update-post", app.updatePostHandler)
		r.Delete("/remove-post", app.deletePostHandler)
	})
	router.Route("/user", func(r chi.Router) {
		r.Post("/", app.createUserHandler)
		r.Route("/{user_id}", func(r chi.Router) {
			r.Get("/", app.getUserHandler)
			r.Put("/", app.updateUserHandler)
			r.Delete("/", app.deleteUserHandler)
		})
	})
	return router
}

func (app *application) run() error {
	mn := app.mount()
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mn,
		WriteTimeout: time.Second * 10,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}
	app.logger.Infow("Server started on ", "addr", srv.Addr, "env", app.config.env)
	return srv.ListenAndServe()
}
