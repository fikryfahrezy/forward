package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fikryfahrezy/forward/blog-api/internal/comment"
	"github.com/fikryfahrezy/forward/blog-api/internal/server"
)

func TestListComments_Success(t *testing.T) {
	cleanup(t)

	token := registerAndGetToken(t, "commenter", "commenter@example.com", "password123")
	postID := createPost(t, token, "Test Post", "Test content")

	for i := 1; i <= 3; i++ {
		reqBody := comment.CreateCommentRequest{
			Content: fmt.Sprintf("Comment %d", i),
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/posts/"+postID+"/comments", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		testServer.Mux().ServeHTTP(rec, req)

		if rec.Code != http.StatusCreated {
			t.Fatalf("Failed to create comment %d: %s", i, rec.Body.String())
		}
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/posts/"+postID+"/comments", nil)
	rec := httptest.NewRecorder()
	testServer.Mux().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d. Body: %s", http.StatusOK, rec.Code, rec.Body.String())
	}

	var response server.APIResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	result := response.Result.(map[string]any)
	comments := result["comments"].([]any)
	if len(comments) != 3 {
		t.Errorf("Expected 3 comments, got %d", len(comments))
	}
}

func TestListComments_Pagination(t *testing.T) {
	cleanup(t)

	token := registerAndGetToken(t, "commenter", "commenter@example.com", "password123")
	postID := createPost(t, token, "Test Post", "Test content")

	for i := 1; i <= 5; i++ {
		reqBody := comment.CreateCommentRequest{
			Content: fmt.Sprintf("Comment %d", i),
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/posts/"+postID+"/comments", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		testServer.Mux().ServeHTTP(rec, req)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/posts/"+postID+"/comments?page=1&page_size=2", nil)
	rec := httptest.NewRecorder()
	testServer.Mux().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d. Body: %s", http.StatusOK, rec.Code, rec.Body.String())
	}

	var response server.APIResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	result := response.Result.(map[string]any)
	comments := result["comments"].([]any)
	if len(comments) != 2 {
		t.Errorf("Expected 2 comments, got %d", len(comments))
	}
	if result["total_count"].(float64) != 5 {
		t.Errorf("Expected total_count 5, got %v", result["total_count"])
	}
}

func TestListComments_PostNotFound(t *testing.T) {
	cleanup(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/posts/550e8400-e29b-41d4-a716-446655440000/comments", nil)
	rec := httptest.NewRecorder()
	testServer.Mux().ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d. Body: %s", http.StatusNotFound, rec.Code, rec.Body.String())
	}
}
