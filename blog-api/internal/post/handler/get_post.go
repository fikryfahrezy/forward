package handler

import (
	"net/http"

	"github.com/fikryfahrezy/forward/blog-api/internal/server"
)

// GetPostBySlug godoc
// @Summary      Get post by slug
// @Description  Get a single blog post by its slug
// @Tags         posts
// @Produce      json
// @Param        slug  path      string  true  "Post slug"
// @Success      200   {object}  server.APIResponse{message=string,result=post.PostItem}  "Post retrieved successfully"
// @Failure      404   {object}  server.APIResponse{message=string,error=string}          "Post not found"
// @Failure      500   {object}  server.APIResponse{message=string,error=string}          "Internal server error"
// @Router       /api/v1/posts/{slug} [get]
func (h *Handler) GetPostBySlug(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	if slug == "" {
		server.ErrorResponse(w, http.StatusBadRequest, "Post slug is required", nil)
		return
	}

	p, err := h.service.GetBySlug(r.Context(), slug)
	if err != nil {
		h.handleError(w, err)
		return
	}

	server.JSON(w, http.StatusOK, server.APIResponse{
		Message: "Post retrieved successfully",
		Result:  p,
	})
}
