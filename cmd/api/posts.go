package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

func (app *application) getPostsHandler(w http.ResponseWriter, r *http.Request) {
	count, err := strconv.Atoi(r.URL.Query().Get("count"))
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	posts, err := app.store.Posts.GetLasts(r.Context(), count)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, posts)
}

type postCreatePayload struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	UserID      int    `json:"user_id"`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var postPayload postCreatePayload
	if err := readJSON(w, r, &postPayload); err != nil {
		app.badRequestError(w, r, err)
		return
	}
	err := app.store.Posts.Create(r.Context(), postPayload.Title, postPayload.Description, postPayload.UserID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, postPayload)
}

type postUpdatePayload struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	ID          int    `json:"id"`
}

func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	var postPayload postUpdatePayload
	if err := readJSON(w, r, &postPayload); err != nil {
		app.badRequestError(w, r, err)
		return
	}
	err := app.store.Posts.Update(r.Context(), postPayload.Title, postPayload.Description, postPayload.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeJSONError(w, http.StatusNotFound, err.Error())
			return
		}
		app.internalServerError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

type postDeletePayload struct {
	ID int `json:"id"`
}

func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	var postID postDeletePayload
	if err := readJSON(w, r, &postID); err != nil {
		app.badRequestError(w, r, err)
		return
	}
	err := app.store.Posts.Delete(r.Context(), postID.ID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, nil)
}

func (app *application) postContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		paramStr := chi.URLParam(r, "post_id")
		param, err := strconv.Atoi(paramStr)
		if err != nil {
			app.badRequestError(w, r, err)
			return
		}
		post, err := app.store.Posts.GetByID(r.Context(), param)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}
		ctx := context.WithValue(r.Context(), postContextKey, post)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
