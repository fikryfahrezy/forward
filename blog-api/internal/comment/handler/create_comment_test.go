package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fikryfahrezy/forward/blog-api/internal/comment"
	"github.com/fikryfahrezy/forward/blog-api/internal/server"
)

func TestCreateComment_Success(t *testing.T) {
	cleanup(t)

	token := registerAndGetToken(t, "commenter", "commenter@example.com", "password123")
	postID := createPost(t, token, "Test Post", "Test content for comments")

	reqBody := comment.CreateCommentRequest{
		Content: "This is a great post!",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/posts/"+postID+"/comments", bytes.NewReader(body))
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

	if response.Message != "Comment created successfully" {
		t.Errorf("Expected message 'Comment created successfully', got '%s'", response.Message)
	}
}

func TestCreateComment_Unauthorized(t *testing.T) {
	cleanup(t)

	token := registerAndGetToken(t, "author", "author@example.com", "password123")
	postID := createPost(t, token, "Test Post", "Test content")

	reqBody := comment.CreateCommentRequest{
		Content: "This is a comment!",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/posts/"+postID+"/comments", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	testServer.Mux().ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d. Body: %s", http.StatusUnauthorized, rec.Code, rec.Body.String())
	}
}

func TestCreateComment_InvalidInput(t *testing.T) {
	cleanup(t)

	token := registerAndGetToken(t, "commenter", "commenter@example.com", "password123")
	postID := createPost(t, token, "Test Post", "Test content")

	reqBody := comment.CreateCommentRequest{
		Content: "",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/posts/"+postID+"/comments", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	testServer.Mux().ServeHTTP(rec, req)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Errorf("Expected status %d, got %d. Body: %s", http.StatusUnprocessableEntity, rec.Code, rec.Body.String())
	}
}

func TestCreateComment_PostNotFound(t *testing.T) {
	cleanup(t)

	token := registerAndGetToken(t, "commenter", "commenter@example.com", "password123")

	reqBody := comment.CreateCommentRequest{
		Content: "Comment on non-existent post",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/posts/550e8400-e29b-41d4-a716-446655440000/comments", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	testServer.Mux().ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d. Body: %s", http.StatusNotFound, rec.Code, rec.Body.String())
	}
}
