package handler

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"

	"github.com/fikryfahrezy/forward/blog-api/internal/server"
)

// ListComments godoc
// @Summary      List comments for a post
// @Description  Get a paginated list of comments for a specific post
// @Tags         comments
// @Produce      json
// @Param        postId    path      string  true   "Post ID"
// @Param        page      query     int     false  "Page number"  default(1)
// @Param        page_size query     int     false  "Page size"    default(10)
// @Success      200       {object}  server.APIResponse{message=string,error=string,result=comment.CommentListResponse}  "Comments retrieved successfully"
// @Failure      404       {object}  server.APIResponse{message=string,error=string}                                     "Post not found"
// @Failure      500       {object}  server.APIResponse{message=string,error=string}                                     "Internal server error"
// @Router       /api/v1/posts/{postId}/comments [get]
func (h *Handler) ListComments(w http.ResponseWriter, r *http.Request) {
	postIDStr := r.PathValue("postId")
	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		server.ErrorResponse(w, http.StatusBadRequest, "Invalid post ID", nil)
		return
	}

	page := 1
	pageSize := 10

	if p := r.URL.Query().Get("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if ps := r.URL.Query().Get("page_size"); ps != "" {
		if parsed, err := strconv.Atoi(ps); err == nil && parsed > 0 {
			pageSize = parsed
		}
	}

	result, err := h.service.ListByPostID(r.Context(), postID, page, pageSize)
	if err != nil {
		h.handleError(w, err)
		return
	}

	server.JSON(w, http.StatusOK, server.APIResponse{
		Message: "Comments retrieved successfully",
		Result:  result,
	})
}
