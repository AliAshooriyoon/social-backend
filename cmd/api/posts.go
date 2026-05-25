package main

import (
	"net/http"
	"strconv"
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

type postPayload struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	UserID      int    `json:"user_id"`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var postPayload postPayload
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
