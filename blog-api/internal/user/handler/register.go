package handler

import (
	"encoding/json"
	"net/http"

	"github.com/fikryfahrezy/forward/blog-api/internal/server"
	"github.com/fikryfahrezy/forward/blog-api/internal/user"
)

// Register godoc
// @Summary      Register a new user
// @Description  Create a new user account and return JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      user.RegisterRequest  true  "Registration details"
// @Success      201      {object}  user.AuthResponse  "User registered successfully"
// @Failure      400      {object}  object{message=string,error=string}  "Invalid input"
// @Failure      409      {object}  object{message=string,error=string}  "Email or username already exists"
// @Failure      500      {object}  object{message=string,error=string}  "Internal server error"
// @Router       /api/v1/auth/register [post]
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req user.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.ErrorResponse(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	u, err := h.service.Register(r.Context(), req)
	if err != nil {
		h.handleError(w, err)
		return
	}

	token, err := h.jwtGenerator.GenerateToken(u.ID.String(), u.Username, u.Email)
	if err != nil {
		server.ErrorResponse(w, http.StatusInternalServerError, "Failed to generate token", nil)
		return
	}

	server.JSON(w, http.StatusCreated, server.APIResponse{
		Message: "User registered successfully",
		Result: user.AuthResponse{
			Token: token,
			User:  u,
		},
	})
}
