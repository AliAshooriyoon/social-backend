package main

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type loginUserPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type claimType struct {
	jwt.RegisteredClaims
	UserID int `json:"user_id"`
}

func (app *application) GenerateToken(w http.ResponseWriter, r *http.Request) {
	var payload loginUserPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.internalServerError(w, r, err)
		return
	}
	user, err := app.store.Users.GetByEmail(r.Context(), payload.Email)
	if err != nil {
		return
	}
	claim := claimType{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
		UserID: user.ID,
	}
	token, err := app.config.auth.GenerateToken(claim)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "refreshToken",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		MaxAge:   60 * 60 * 24,
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	if err := writeJSON(w, http.StatusOK, token); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
