package handler

import (
	"net/http"
	"time"

	"github.com/fikryfahrezy/forward/blog-api/internal/server"
	"github.com/fikryfahrezy/forward/blog-api/internal/user"
	"github.com/fikryfahrezy/forward/blog-api/internal/user/service"
)

type Handler struct {
	service      *service.Service
	jwtGenerator *server.JWTGenerator
}

func New(svc *service.Service, jwtSecretKey string, tokenDuration time.Duration) *Handler {
	return &Handler{
		service:      svc,
		jwtGenerator: server.NewJWTGenerator(jwtSecretKey, tokenDuration),
	}
}

func (h *Handler) SetupRoutes(server *server.Server) {
	// Public routes (no auth required)
	server.HandleFunc("POST /api/v1/auth/register", h.Register)
	server.HandleFunc("POST /api/v1/auth/login", h.Login)

	// Protected routes (auth required)
	server.HandleFuncWithAuth("GET /api/v1/users/me", h.GetCurrentUser)
}

func (h *Handler) handleError(w http.ResponseWriter, err error) {
	switch err {
	case user.ErrInvalidInput:
		server.ErrorResponse(w, http.StatusUnprocessableEntity, "", err)
	case user.ErrEmailExists, user.ErrUsernameExists:
		server.ErrorResponse(w, http.StatusConflict, "", err)
	case user.ErrInvalidCredentials, user.ErrUserNotFound:
		server.ErrorResponse(w, http.StatusUnauthorized, "", err)
	default:
		server.ErrorResponse(w, http.StatusInternalServerError, "Internal server error", nil)
	}
}
