package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fikryfahrezy/forward/blog-api/internal/post"
	"github.com/fikryfahrezy/forward/blog-api/internal/server"
)

func TestCreatePost_Success(t *testing.T) {
	cleanup(t)

	token := registerAndGetToken(t, "postauthor", "author@example.com", "password123")

	reqBody := post.CreatePostRequest{
		Title:   "My First Post",
		Content: "This is the content of my first post.",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/posts", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	testServer.Mux().ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d. Body: %s", http.StatusCreated, rec.Code, rec.Body.String())
	}

	var response server.APIResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response.Message != "Post created successfully" {
		t.Errorf("Expected message 'Post created successfully', got '%s'", response.Message)
	}
}

func TestCreatePost_Unauthorized(t *testing.T) {
	cleanup(t)

	reqBody := post.CreatePostRequest{
		Title:   "My First Post",
		Content: "This is the content of my first post.",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/posts", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	testServer.Mux().ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d. Body: %s", http.StatusUnauthorized, rec.Code, rec.Body.String())
	}
}

func TestCreatePost_InvalidInput(t *testing.T) {
	cleanup(t)

	token := registerAndGetToken(t, "postauthor", "author@example.com", "password123")

	tests := []struct {
		name    string
		request post.CreatePostRequest
	}{
		{
			name:    "Empty title",
			request: post.CreatePostRequest{Title: "", Content: "Some content"},
		},
		{
			name:    "Empty content",
			request: post.CreatePostRequest{Title: "Some title", Content: ""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.request)

			req := httptest.NewRequest(http.MethodPost, "/api/v1/posts", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)
			rec := httptest.NewRecorder()
			testServer.Mux().ServeHTTP(rec, req)

			if rec.Code != http.StatusUnprocessableEntity {
				t.Errorf("Expected status %d, got %d. Body: %s", http.StatusUnprocessableEntity, rec.Code, rec.Body.String())
			}
		})
	}
}
