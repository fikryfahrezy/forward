package handler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"github.com/fikryfahrezy/forward/blog-api/internal/post"
	"github.com/fikryfahrezy/forward/blog-api/internal/server"
)

// CreatePost godoc
// @Summary      Create a new post
// @Description  Create a new blog post
// @Tags         posts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      post.CreatePostRequest  true  "Post details"
// @Success      201      {object}  server.APIResponse{message=string,result=post.PostID}  "Post created successfully"
// @Failure      400      {object}  server.APIResponse{message=string,error=string}        "Invalid input"
// @Failure      401      {object}  server.APIResponse{message=string,error=string}        "Unauthorized"
// @Failure      500      {object}  server.APIResponse{message=string,error=string}        "Internal server error"
// @Router       /api/v1/posts [post]
func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	claims, ok := server.GetUserClaims(r.Context())
	if !ok {
		server.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	authorID, err := uuid.Parse(claims.UserID)
	if err != nil {
		server.ErrorResponse(w, http.StatusUnauthorized, "Invalid user ID", nil)
		return
	}

	var req post.CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.ErrorResponse(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	p, err := h.service.Create(r.Context(), authorID, req)
	if err != nil {
		h.handleError(w, err)
		return
	}

	server.JSON(w, http.StatusCreated, server.APIResponse{
		Message: "Post created successfully",
		Result:  p,
	})
}
