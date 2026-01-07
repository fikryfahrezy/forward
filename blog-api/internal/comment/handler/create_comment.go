package handler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"github.com/fikryfahrezy/forward/blog-api/internal/comment"
	"github.com/fikryfahrezy/forward/blog-api/internal/server"
)

// CreateComment godoc
// @Summary      Create a new comment
// @Description  Create a new comment on a post
// @Tags         comments
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        postId   path      string                        true  "Post ID"
// @Param        request  body      comment.CreateCommentRequest  true  "Comment details"
// @Success      201      {object}  server.APIResponse{message=string,result=comment.CommentID}  "Comment created successfully"
// @Failure      400      {object}  server.APIResponse{message=string,error=string}              "Invalid input"
// @Failure      401      {object}  server.APIResponse{message=string,error=string}              "Unauthorized"
// @Failure      404      {object}  server.APIResponse{message=string,error=string}              "Post not found"
// @Failure      500      {object}  server.APIResponse{message=string,error=string}              "Internal server error"
// @Router       /api/v1/posts/{postId}/comments [post]
func (h *Handler) CreateComment(w http.ResponseWriter, r *http.Request) {
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

	postIDStr := r.PathValue("postId")
	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		server.ErrorResponse(w, http.StatusBadRequest, "Invalid post ID", nil)
		return
	}

	var req comment.CreateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.ErrorResponse(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	c, err := h.service.Create(r.Context(), postID, authorID, req)
	if err != nil {
		h.handleError(w, err)
		return
	}

	server.JSON(w, http.StatusCreated, server.APIResponse{
		Message: "Comment created successfully",
		Result:  c,
	})
}
