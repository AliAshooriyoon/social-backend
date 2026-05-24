package main

import "net/http"

type PostsPayloadType struct {
	countOfPosts int
}

func (app *application) getPostsHandler(w http.ResponseWriter, r *http.Request) {
	var postsPayload PostsPayloadType
	err := readJSON(w, r, postsPayload)
	if err != nil {
		return
	}
}
