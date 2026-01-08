package handler

import (
	"encoding/json"
	"net/http"

	"github.com/fikryfahrezy/forward/blog-api/internal/server"
	"github.com/fikryfahrezy/forward/blog-api/internal/user"
)

// Login godoc
// @Summary      Login user
// @Description  Authenticate user and return JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      user.LoginRequest  true  "Login credentials"
// @Success      200      {object}  server.APIResponse{message=string,result=user.AuthResponse}  "Login successful"
// @Failure      400      {object}  server.APIResponse{message=string,error=string}              "Invalid input"
// @Failure      401      {object}  server.APIResponse{message=string,error=string}              "Invalid credentials"
// @Failure      500      {object}  server.APIResponse{message=string,error=string}              "Internal server error"
// @Router       /api/v1/auth/login [post]
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req user.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.ErrorResponse(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	u, err := h.service.Login(r.Context(), req)
	if err != nil {
		h.handleError(w, err)
		return
	}

	server.JSON(w, http.StatusOK, server.APIResponse{
		Message: "Login successful",
		Result:  u,
	})
}
