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

func TestRegister_Success(t *testing.T) {
	cleanupUsers(t)

	reqBody := user.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	testServer.Mux().ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d. Body: %s", http.StatusCreated, rec.Code, rec.Body.String())
	}

	var response server.APIResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response.Message != "User registered successfully" {
		t.Errorf("Expected message 'User registered successfully', got '%s'", response.Message)
	}

	// Verify token is returned
	result, ok := response.Result.(map[string]any)
	if !ok {
		t.Fatal("Expected result to be a map")
	}

	if _, exists := result["token"]; !exists {
		t.Error("Expected token in response")
	}

	if userMap, exists := result["user"].(map[string]any); exists {
		if userMap["username"] != "testuser" {
			t.Errorf("Expected username 'testuser', got '%v'", userMap["username"])
		}
		if userMap["email"] != "test@example.com" {
			t.Errorf("Expected email 'test@example.com', got '%v'", userMap["email"])
		}
	} else {
		t.Error("Expected user in response")
	}
}

func TestRegister_DuplicateEmail(t *testing.T) {
	cleanupUsers(t)

	// Register first user
	reqBody := user.RegisterRequest{
		Username: "testuser1",
		Email:    "duplicate@example.com",
		Password: "password123",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	testServer.Mux().ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("First registration failed: %s", rec.Body.String())
	}

	// Try to register with same email
	reqBody2 := user.RegisterRequest{
		Username: "testuser2",
		Email:    "duplicate@example.com",
		Password: "password456",
	}
	body2, _ := json.Marshal(reqBody2)

	req2 := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(body2))
	req2.Header.Set("Content-Type", "application/json")
	rec2 := httptest.NewRecorder()
	testServer.Mux().ServeHTTP(rec2, req2)

	if rec2.Code != http.StatusConflict {
		t.Errorf("Expected status %d, got %d. Body: %s", http.StatusConflict, rec2.Code, rec2.Body.String())
	}
}

func TestRegister_DuplicateUsername(t *testing.T) {
	cleanupUsers(t)

	// Register first user
	reqBody := user.RegisterRequest{
		Username: "duplicateuser",
		Email:    "test1@example.com",
		Password: "password123",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	testServer.Mux().ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("First registration failed: %s", rec.Body.String())
	}

	// Try to register with same username
	reqBody2 := user.RegisterRequest{
		Username: "duplicateuser",
		Email:    "test2@example.com",
		Password: "password456",
	}
	body2, _ := json.Marshal(reqBody2)

	req2 := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(body2))
	req2.Header.Set("Content-Type", "application/json")
	rec2 := httptest.NewRecorder()
	testServer.Mux().ServeHTTP(rec2, req2)

	if rec2.Code != http.StatusConflict {
		t.Errorf("Expected status %d, got %d. Body: %s", http.StatusConflict, rec2.Code, rec2.Body.String())
	}
}

func TestRegister_InvalidInput(t *testing.T) {
	tests := []struct {
		name    string
		request user.RegisterRequest
	}{
		{
			name:    "Empty username",
			request: user.RegisterRequest{Username: "", Email: "test@example.com", Password: "password123"},
		},
		{
			name:    "Empty email",
			request: user.RegisterRequest{Username: "testuser", Email: "", Password: "password123"},
		},
		{
			name:    "Empty password",
			request: user.RegisterRequest{Username: "testuser", Email: "test@example.com", Password: ""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.request)

			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			testServer.Mux().ServeHTTP(rec, req)

			if rec.Code != http.StatusUnprocessableEntity {
				t.Errorf("Expected status %d, got %d. Body: %s", http.StatusUnprocessableEntity, rec.Code, rec.Body.String())
			}
		})
	}
}
