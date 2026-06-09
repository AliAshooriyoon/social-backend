package auth

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

type JWTAuth struct {
	secret string
	iss    string
}

func NewAuth(secret, iss string) *JWTAuth {
	return &JWTAuth{secret: secret, iss: iss}
}

func (j *JWTAuth) GenerateToken(claim jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenStr, err := token.SignedString(j.secret)
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

func (j *JWTAuth) ValidateToken(tokenStr string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, errors.New("invalid alg")
		}
		return []byte(j.secret), nil
	}, jwt.WithExpirationRequired())
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	return token, nil
}
