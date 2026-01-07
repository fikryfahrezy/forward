package handler

import (
	"net/http"

	"github.com/google/uuid"

	"github.com/fikryfahrezy/forward/blog-api/internal/server"
)

// GetCurrentUser godoc
// @Summary      Get current user
// @Description  Get the currently authenticated user's profile
// @Tags         users
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  server.APIResponse{message=string,result=user.UserProfile}  "User retrieved successfully"
// @Failure      401  {object}  server.APIResponse{message=string,error=string}             "Unauthorized"
// @Failure      404  {object}  server.APIResponse{message=string,error=string}             "User not found"
// @Failure      500  {object}  server.APIResponse{message=string,error=string}             "Internal server error"
// @Router       /api/v1/users/me [get]
func (h *Handler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	claims, ok := server.GetUserClaims(r.Context())
	if !ok {
		server.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		server.ErrorResponse(w, http.StatusUnauthorized, "Invalid user ID", nil)
		return
	}

	u, err := h.service.GetByID(r.Context(), userID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	server.JSON(w, http.StatusOK, server.APIResponse{
		Message: "User retrieved successfully",
		Result:  u,
	})
}
