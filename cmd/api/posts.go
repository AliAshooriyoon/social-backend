package main

import (
	"errors"
	"net/http"
	"strconv"

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
