package http

import (
	"net/http"

	appError "github.com/fikryfahrezy/forward/blog-api/internal/error"
)

// APIResponse represents a standard API response
type APIResponse struct {
	Message string `json:"message"`
	Error   string `json:"error"`
	Result  any    `json:"result"`
}

func ErrorResponse(w http.ResponseWriter, statusCode int, message string, err error) {
	errorCode := http.StatusText(statusCode)
	errorMessage := message

	if err != nil {
		// Extract error code if it's an AppError
		code := appError.GetCode(err)
		if code != "" {
			errorCode = code
		}

		// Extract error message if it's an AppError
		message := appError.GetMessage(err)
		if message != "" {
			errorMessage = message
		}

	}

	JSON(w, statusCode, APIResponse{
		Message: errorMessage,
		Error:   errorCode,
	})
}
