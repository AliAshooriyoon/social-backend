package main

import (
	"context"
	"net/http"
)

type keyContextType string

var contextKey keyContextType = "user_auth"

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

		ctx := context.WithValue(r.Context(), contextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
