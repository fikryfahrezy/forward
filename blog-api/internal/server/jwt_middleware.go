package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const (
	UserClaimsKey contextKey = "user_claims"
	UserTokenKey  contextKey = "user_token"
)

type JWTConfig struct {
	SecretKey string
}

type JWTMiddleware struct {
	secretKey []byte
}

type UserClaims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func NewJWTMiddleware(config JWTConfig) *JWTMiddleware {
	return &JWTMiddleware{
		secretKey: []byte(config.SecretKey),
	}
}

func (m *JWTMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			ErrorResponse(w, http.StatusUnauthorized, "Authorization header is required", nil)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			ErrorResponse(w, http.StatusUnauthorized, "Invalid authorization header format", nil)
			return
		}

		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return m.secretKey, nil
		})

		if err != nil || !token.Valid {
			ErrorResponse(w, http.StatusUnauthorized, "Invalid or expired token", nil)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			ErrorResponse(w, http.StatusUnauthorized, "Invalid token claims", nil)
			return
		}

		userClaims := extractUserClaims(claims)

		ctx := context.WithValue(r.Context(), UserClaimsKey, userClaims)
		ctx = context.WithValue(ctx, UserTokenKey, tokenString)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func extractUserClaims(claims jwt.MapClaims) UserClaims {
	userClaims := UserClaims{}

	if userID, ok := claims["user_id"].(string); ok {
		userClaims.UserID = userID
	}
	if username, ok := claims["username"].(string); ok {
		userClaims.Username = username
	}
	if email, ok := claims["email"].(string); ok {
		userClaims.Email = email
	}

	return userClaims
}

func GetUserClaims(ctx context.Context) (UserClaims, bool) {
	claims, ok := ctx.Value(UserClaimsKey).(UserClaims)
	return claims, ok
}

func GetUserToken(ctx context.Context) (string, bool) {
	token, ok := ctx.Value(UserTokenKey).(string)
	return token, ok
}
