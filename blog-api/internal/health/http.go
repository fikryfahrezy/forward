package health

import (
	"log/slog"
	"net/http"

	"github.com/fikryfahrezy/forward/blog-api/internal/database"
	"github.com/fikryfahrezy/forward/blog-api/internal/server"
)

type HealthHandler struct {
	db *database.DB
}

// DatabaseCheck represents database health status
type DatabaseCheck struct {
	Status string `json:"status" example:"ok"`
	Error  string `json:"error,omitempty" example:"connection refused"`
}

// HealthChecks contains all service dependency checks
type HealthChecks struct {
	Database DatabaseCheck `json:"database"`
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status  string       `json:"status" example:"ok"`
	Message string       `json:"message" example:"Service is healthy"`
	Checks  HealthChecks `json:"checks"`
}

func NewHealthHandler(db *database.DB) *HealthHandler {
	return &HealthHandler{
		db: db,
	}
}

// HealthCheck godoc
// @Summary      Check service health
// @Description  Returns the health status of the service and its dependencies
// @Tags         health
// @Produce      json
// @Success      200  {object}  HealthResponse  "Service is healthy"
// @Failure      503  {object}  HealthResponse  "Service is unhealthy"
// @Router       /api/health [get]
func (h *HealthHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	status := "ok"
	message := "Service is healthy"
	httpStatus := http.StatusOK

	dbCheck := map[string]any{"status": "unknown"}

	// Check database connection
	if err := h.db.Health(r.Context()); err != nil {
		slog.Error("Database health check failed",
			slog.String("error", err.Error()),
		)
		status = "unhealthy"
		message = "Database connection failed"
		httpStatus = http.StatusServiceUnavailable
		dbCheck = map[string]any{
			"status": "unhealthy",
			"error":  err.Error(),
		}
	}

	response := map[string]any{
		"status":  status,
		"message": message,
		"checks": map[string]any{
			"database": dbCheck,
		},
	}

	server.JSON(w, httpStatus, response)
}

// SetupRoutes configures health check routes
func (h *HealthHandler) SetupRoutes(server *server.Server) {
	// Health check endpoint (no versioning needed)
	server.HandleFunc("GET /api/health", h.HealthCheck)
}
