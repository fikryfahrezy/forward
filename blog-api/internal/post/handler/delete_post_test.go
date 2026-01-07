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

func TestDeletePost_Success(t *testing.T) {
	cleanup(t)

	token := registerAndGetToken(t, "postauthor", "author@example.com", "password123")

	// Create a post
	createReq := post.CreatePostRequest{
		Title:   "Post to Delete",
		Content: "This will be deleted.",
	}
	body, _ := json.Marshal(createReq)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/posts", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	testServer.Mux().ServeHTTP(rec, req)

	var createResponse server.APIResponse
	//nolint:errcheck
	json.Unmarshal(rec.Body.Bytes(), &createResponse)
	result := createResponse.Result.(map[string]any)
	postID := result["id"].(string)

	// Delete the post
	req = httptest.NewRequest(http.MethodDelete, "/api/v1/posts/"+postID, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec = httptest.NewRecorder()
	testServer.Mux().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d. Body: %s", http.StatusOK, rec.Code, rec.Body.String())
	}

	// Verify post is deleted (list should be empty)
	req = httptest.NewRequest(http.MethodGet, "/api/v1/posts", nil)
	rec = httptest.NewRecorder()
	testServer.Mux().ServeHTTP(rec, req)

	var listResponse server.APIResponse
	//nolint:errcheck
	json.Unmarshal(rec.Body.Bytes(), &listResponse)
	listResult := listResponse.Result.(map[string]any)
	posts := listResult["posts"].([]any)
	if len(posts) != 0 {
		t.Errorf("Expected 0 posts after deletion, got %d", len(posts))
	}
}

func TestDeletePost_NotAuthor(t *testing.T) {
	cleanup(t)

	// User 1 creates a post
	token1 := registerAndGetToken(t, "author1", "author1@example.com", "password123")

	createReq := post.CreatePostRequest{
		Title:   "Author1 Post",
		Content: "Content by author1.",
	}
	body, _ := json.Marshal(createReq)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/posts", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token1)
	rec := httptest.NewRecorder()
	testServer.Mux().ServeHTTP(rec, req)

	var createResponse server.APIResponse
	//nolint:errcheck
	json.Unmarshal(rec.Body.Bytes(), &createResponse)
	result := createResponse.Result.(map[string]any)
	postID := result["id"].(string)

	// User 2 tries to delete
	token2 := registerAndGetToken(t, "author2", "author2@example.com", "password123")

	req = httptest.NewRequest(http.MethodDelete, "/api/v1/posts/"+postID, nil)
	req.Header.Set("Authorization", "Bearer "+token2)
	rec = httptest.NewRecorder()
	testServer.Mux().ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Errorf("Expected status %d, got %d. Body: %s", http.StatusForbidden, rec.Code, rec.Body.String())
	}
}

func TestDeletePost_Unauthorized(t *testing.T) {
	cleanup(t)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/posts/some-id", nil)
	rec := httptest.NewRecorder()
	testServer.Mux().ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d. Body: %s", http.StatusUnauthorized, rec.Code, rec.Body.String())
	}
}

func TestDeletePost_NotFound(t *testing.T) {
	cleanup(t)

	token := registerAndGetToken(t, "postauthor", "author@example.com", "password123")

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/posts/550e8400-e29b-41d4-a716-446655440000", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	testServer.Mux().ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d. Body: %s", http.StatusNotFound, rec.Code, rec.Body.String())
	}
}
