package auth

import "github.com/golang-jwt/jwt/v5"

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
