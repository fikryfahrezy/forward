package handler

import (
	"net/http"

	"github.com/fikryfahrezy/forward/blog-api/internal/comment"
	"github.com/fikryfahrezy/forward/blog-api/internal/comment/service"
	"github.com/fikryfahrezy/forward/blog-api/internal/server"
)

type Handler struct {
	service *service.Service
}

func New(svc *service.Service) *Handler {
	return &Handler{
		service: svc,
	}
}

func (h *Handler) SetupRoutes(server *server.Server) {
	// Public routes
	server.HandleFunc("GET /api/v1/posts/{postId}/comments", h.ListComments)

	// Protected routes
	server.HandleFuncWithAuth("POST /api/v1/posts/{postId}/comments", h.CreateComment)
	server.HandleFuncWithAuth("PUT /api/v1/comments/{commentId}", h.UpdateComment)
	server.HandleFuncWithAuth("DELETE /api/v1/comments/{commentId}", h.DeleteComment)
}

func (h *Handler) handleError(w http.ResponseWriter, err error) {
	switch err {
	case comment.ErrInvalidInput:
		server.ErrorResponse(w, http.StatusUnprocessableEntity, "", err)
	case comment.ErrCommentNotFound, comment.ErrPostNotFound:
		server.ErrorResponse(w, http.StatusNotFound, "", err)
	case comment.ErrUnauthorized:
		server.ErrorResponse(w, http.StatusForbidden, "", err)
	default:
		server.ErrorResponse(w, http.StatusInternalServerError, "Internal server error", nil)
	}
}
