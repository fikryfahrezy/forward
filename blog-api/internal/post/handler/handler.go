package handler

import (
	"net/http"

	"github.com/fikryfahrezy/forward/blog-api/internal/post"
	"github.com/fikryfahrezy/forward/blog-api/internal/post/service"
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
	server.HandleFunc("GET /api/v1/posts", h.ListPosts)
	server.HandleFunc("GET /api/v1/posts/{slug}", h.GetPostBySlug)

	// Protected routes
	server.HandleFuncWithAuth("POST /api/v1/posts", h.CreatePost)
	server.HandleFuncWithAuth("PUT /api/v1/posts/{postId}", h.UpdatePost)
	server.HandleFuncWithAuth("DELETE /api/v1/posts/{postId}", h.DeletePost)
}

func (h *Handler) handleError(w http.ResponseWriter, err error) {
	switch err {
	case post.ErrInvalidInput:
		server.ErrorResponse(w, http.StatusUnprocessableEntity, "", err)
	case post.ErrPostNotFound:
		server.ErrorResponse(w, http.StatusNotFound, "", err)
	case post.ErrUnauthorized:
		server.ErrorResponse(w, http.StatusForbidden, "", err)
	case post.ErrSlugExists:
		server.ErrorResponse(w, http.StatusConflict, "", err)
	default:
		server.ErrorResponse(w, http.StatusInternalServerError, "Internal server error", nil)
	}
}
