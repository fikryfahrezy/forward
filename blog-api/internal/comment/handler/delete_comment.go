package handler

import (
	"net/http"

	"github.com/google/uuid"

	"github.com/fikryfahrezy/forward/blog-api/internal/server"
)

// DeleteComment godoc
// @Summary      Delete a comment
// @Description  Delete an existing comment (only the author can delete)
// @Tags         comments
// @Produce      json
// @Security     BearerAuth
// @Param        commentId  path      string  true  "Comment ID"
// @Success      200        {object}  server.APIResponse{message=string,result=comment.CommentID}  "Comment deleted successfully"
// @Failure      401        {object}  server.APIResponse{message=string,error=string}              "Unauthorized"
// @Failure      403        {object}  server.APIResponse{message=string,error=string}              "Forbidden - not the author"
// @Failure      404        {object}  server.APIResponse{message=string,error=string}              "Comment not found"
// @Failure      500        {object}  server.APIResponse{message=string,error=string}              "Internal server error"
// @Router       /api/v1/comments/{commentId} [delete]
func (h *Handler) DeleteComment(w http.ResponseWriter, r *http.Request) {
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

	commentIDStr := r.PathValue("commentId")
	commentID, err := uuid.Parse(commentIDStr)
	if err != nil {
		server.ErrorResponse(w, http.StatusBadRequest, "Invalid comment ID", nil)
		return
	}

	c, err := h.service.Delete(r.Context(), commentID, authorID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	server.JSON(w, http.StatusOK, server.APIResponse{
		Message: "Comment deleted successfully",
		Result:  c,
	})
}
