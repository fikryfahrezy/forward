package handler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"github.com/fikryfahrezy/forward/blog-api/internal/post"
	"github.com/fikryfahrezy/forward/blog-api/internal/server"
)

// UpdatePost godoc
// @Summary      Update a post
// @Description  Update an existing blog post (only the author can update)
// @Tags         posts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        postId       path      string                  true  "Post ID"
// @Param        request      body      post.UpdatePostRequest  true  "Updated post details"
// @Success      200          {object}  server.APIResponse{message=string,result=post.PostID}  "Post updated successfully"
// @Failure      400          {object}  server.APIResponse{message=string,error=string}        "Invalid input"
// @Failure      401          {object}  server.APIResponse{message=string,error=string}        "Unauthorized"
// @Failure      403          {object}  server.APIResponse{message=string,error=string}        "Forbidden - not the author"
// @Failure      404          {object}  server.APIResponse{message=string,error=string}        "Post not found"
// @Failure      500          {object}  server.APIResponse{message=string,error=string}        "Internal server error"
// @Router       /api/v1/posts/{postId} [put]
func (h *Handler) UpdatePost(w http.ResponseWriter, r *http.Request) {
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

	var req post.UpdatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.ErrorResponse(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	p, err := h.service.Update(r.Context(), postID, authorID, req)
	if err != nil {
		h.handleError(w, err)
		return
	}

	server.JSON(w, http.StatusOK, server.APIResponse{
		Message: "Post updated successfully",
		Result:  p,
	})
}
