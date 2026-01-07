package handler

import (
	"net/http"

	"github.com/google/uuid"

	"github.com/fikryfahrezy/forward/blog-api/internal/server"
)

// DeletePost godoc
// @Summary      Delete a post
// @Description  Delete an existing blog post (only the author can delete)
// @Tags         posts
// @Produce      json
// @Security     BearerAuth
// @Param        postId  path      string  true  "Post ID"
// @Success      200     {object}  server.APIResponse{message=string,result=post.PostID}  "Post deleted successfully"
// @Failure      401     {object}  server.APIResponse{message=string,error=string}        "Unauthorized"
// @Failure      403     {object}  server.APIResponse{message=string,error=string}        "Forbidden - not the author"
// @Failure      404     {object}  server.APIResponse{message=string,error=string}        "Post not found"
// @Failure      500     {object}  server.APIResponse{message=string,error=string}        "Internal server error"
// @Router       /api/v1/posts/{postId} [delete]
func (h *Handler) DeletePost(w http.ResponseWriter, r *http.Request) {
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

	p, err := h.service.Delete(r.Context(), postID, authorID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	server.JSON(w, http.StatusOK, server.APIResponse{
		Message: "Post deleted successfully",
		Result:  p,
	})
}
