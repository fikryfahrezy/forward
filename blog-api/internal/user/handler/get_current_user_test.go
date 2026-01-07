package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fikryfahrezy/forward/blog-api/internal/server"
	"github.com/fikryfahrezy/forward/blog-api/internal/user"
)

func TestGetCurrentUser_Success(t *testing.T) {
	cleanupUsers(t)

	// First register a user and get token
	registerBody := user.RegisterRequest{
		Username: "meuser",
		Email:    "me@example.com",
		Password: "password123",
	}
	body, _ := json.Marshal(registerBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	testServer.Mux().ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("Registration failed: %s", rec.Body.String())
	}

	var registerResponse server.APIResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &registerResponse); err != nil {
		t.Fatalf("Failed to parse register response: %v", err)
	}

	result := registerResponse.Result.(map[string]any)
	token := result["token"].(string)

	// Now get current user
	req = httptest.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec = httptest.NewRecorder()
	testServer.Mux().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d. Body: %s", http.StatusOK, rec.Code, rec.Body.String())
	}

	var response server.APIResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response.Message != "User retrieved successfully" {
		t.Errorf("Expected message 'User retrieved successfully', got '%s'", response.Message)
	}

	userResult, ok := response.Result.(map[string]any)
	if !ok {
		t.Fatal("Expected result to be a user map")
	}

	if userResult["username"] != "meuser" {
		t.Errorf("Expected username 'meuser', got '%v'", userResult["username"])
	}
	if userResult["email"] != "me@example.com" {
		t.Errorf("Expected email 'me@example.com', got '%v'", userResult["email"])
	}
}

func TestGetCurrentUser_NoToken(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
	rec := httptest.NewRecorder()
	testServer.Mux().ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d. Body: %s", http.StatusUnauthorized, rec.Code, rec.Body.String())
	}
}

func TestGetCurrentUser_InvalidToken(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	rec := httptest.NewRecorder()
	testServer.Mux().ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d. Body: %s", http.StatusUnauthorized, rec.Code, rec.Body.String())
	}
}

func TestGetCurrentUser_MalformedAuthHeader(t *testing.T) {
	tests := []struct {
		name       string
		authHeader string
	}{
		{name: "No Bearer prefix", authHeader: "some-token"},
		{name: "Wrong prefix", authHeader: "Basic some-token"},
		{name: "Empty token", authHeader: "Bearer "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
			req.Header.Set("Authorization", tt.authHeader)
			rec := httptest.NewRecorder()
			testServer.Mux().ServeHTTP(rec, req)

			if rec.Code != http.StatusUnauthorized {
				t.Errorf("Expected status %d, got %d. Body: %s", http.StatusUnauthorized, rec.Code, rec.Body.String())
			}
		})
	}
}
