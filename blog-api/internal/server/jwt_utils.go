package server

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTGenerator struct {
	secretKey     []byte
	tokenDuration time.Duration
}

func NewJWTGenerator(secretKey string, tokenDuration time.Duration) *JWTGenerator {
	return &JWTGenerator{
		secretKey:     []byte(secretKey),
		tokenDuration: tokenDuration,
	}
}

func (g *JWTGenerator) GenerateToken(userID, username, email string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"email":    email,
		"exp":      time.Now().Add(g.tokenDuration).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(g.secretKey)
}
