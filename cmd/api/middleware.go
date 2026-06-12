package main

import (
	"context"
	"fmt"
	"net/http"

	"project/internal/store"
)

func (app *application) AuthTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("refreshToken")
		if err != nil {
			app.unauthorizedError(w, r, err)
			return
		}
		token := cookie.Value
		jwtToken, err := app.config.auth.ValidateToken(token)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}
		claim, _ := jwtToken.Claims.(*ClaimType)
		user, err := app.store.Users.GetByID(r.Context(), claim.UserID)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) checkPostOwnership(requiredRole string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := app.getUserFromCTX(r)
		if err != nil {
			app.unauthorizedError(w, r, err)
			return
		}
		post, err := app.getPostFromCTX(r)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}
		if post.UserID == user.ID {
			next.ServeHTTP(w, r)
		}
		if user.Role > 2 {
			next.ServeHTTP(w, r)
		}
	}
}

func (app application) getUserFromCTX(r *http.Request) (*store.User, error) {
	user, ok := r.Context().Value(userContextKey).(*store.User)
	if !ok {
		return nil, fmt.Errorf("invalid type of user")
	}
	return user, nil
}

func (app application) getPostFromCTX(r *http.Request) (*store.Post, error) {
	user, ok := r.Context().Value(userContextKey).(*store.Post)
	if !ok {
		return nil, fmt.Errorf("invalid type of post")
	}
	return user, nil
}
