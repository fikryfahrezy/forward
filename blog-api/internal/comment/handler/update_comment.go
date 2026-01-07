package handler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"github.com/fikryfahrezy/forward/blog-api/internal/comment"
	"github.com/fikryfahrezy/forward/blog-api/internal/server"
)

// UpdateComment godoc
// @Summary      Update a comment
// @Description  Update an existing comment (only the author can update)
// @Tags         comments
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        commentId  path      string                        true  "Comment ID"
// @Param        request    body      comment.UpdateCommentRequest  true  "Updated comment details"
// @Success      200        {object}  server.APIResponse{message=string,result=comment.CommentID}  "Comment updated successfully"
// @Failure      400        {object}  server.APIResponse{message=string,error=string}              "Invalid input"
// @Failure      401        {object}  server.APIResponse{message=string,error=string}              "Unauthorized"
// @Failure      403        {object}  server.APIResponse{message=string,error=string}              "Forbidden - not the author"
// @Failure      404        {object}  server.APIResponse{message=string,error=string}              "Comment not found"
// @Failure      500        {object}  server.APIResponse{message=string,error=string}              "Internal server error"
// @Router       /api/v1/comments/{commentId} [put]
func (h *Handler) UpdateComment(w http.ResponseWriter, r *http.Request) {
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

	var req comment.UpdateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.ErrorResponse(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	c, err := h.service.Update(r.Context(), commentID, authorID, req)
	if err != nil {
		h.handleError(w, err)
		return
	}

	server.JSON(w, http.StatusOK, server.APIResponse{
		Message: "Comment updated successfully",
		Result:  c,
	})
}
