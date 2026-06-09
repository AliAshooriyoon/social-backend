// Package auth: authentication
package auth

import "github.com/golang-jwt/jwt/v5"

func GenerateToken() {}

type Authenticator interface {
	GenerateToken(claim jwt.Claims) (string, error)
	ValidateToken(string) (*jwt.Token, error)
}
