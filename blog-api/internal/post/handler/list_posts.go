package handler

import (
	"net/http"
	"strconv"

	"github.com/fikryfahrezy/forward/blog-api/internal/server"
)

// ListPosts godoc
// @Summary      List all posts
// @Description  Get a paginated list of all blog posts
// @Tags         posts
// @Produce      json
// @Param        page      query     int  false  "Page number"  default(1)
// @Param        page_size query     int  false  "Page size"    default(10)
// @Success      200       {object}  server.APIResponse{message=string,result=post.PostListResponse}  "Posts retrieved successfully"
// @Failure      500       {object}  server.APIResponse{message=string,error=string}                  "Internal server error"
// @Router       /api/v1/posts [get]
func (h *Handler) ListPosts(w http.ResponseWriter, r *http.Request) {
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

	result, err := h.service.List(r.Context(), page, pageSize)
	if err != nil {
		h.handleError(w, err)
		return
	}

	server.JSON(w, http.StatusOK, server.APIResponse{
		Message: "Posts retrieved successfully",
		Result:  result,
	})
}
