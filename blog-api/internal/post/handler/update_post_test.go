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

func TestUpdatePost_Success(t *testing.T) {
	cleanup(t)

	token := registerAndGetToken(t, "postauthor", "author@example.com", "password123")

	// Create a post
	createReq := post.CreatePostRequest{
		Title:   "Original Title",
		Content: "Original content.",
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

	// Update the post
	updateReq := post.UpdatePostRequest{
		Title:   "Updated Title",
		Content: "Updated content.",
	}
	body, _ = json.Marshal(updateReq)

	req = httptest.NewRequest(http.MethodPut, "/api/v1/posts/"+postID, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
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

	if response.Message != "Post updated successfully" {
		t.Errorf("Expected message 'Post updated successfully', got '%s'", response.Message)
	}
}

func TestUpdatePost_NotAuthor(t *testing.T) {
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

	// User 2 tries to update
	token2 := registerAndGetToken(t, "author2", "author2@example.com", "password123")

	updateReq := post.UpdatePostRequest{
		Title:   "Hacked Title",
		Content: "Hacked content.",
	}
	body, _ = json.Marshal(updateReq)

	req = httptest.NewRequest(http.MethodPut, "/api/v1/posts/"+postID, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token2)
	rec = httptest.NewRecorder()
	testServer.Mux().ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Errorf("Expected status %d, got %d. Body: %s", http.StatusForbidden, rec.Code, rec.Body.String())
	}
}

func TestUpdatePost_Unauthorized(t *testing.T) {
	cleanup(t)

	updateReq := post.UpdatePostRequest{
		Title:   "Updated Title",
		Content: "Updated content.",
	}
	body, _ := json.Marshal(updateReq)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/posts/some-id", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	testServer.Mux().ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d. Body: %s", http.StatusUnauthorized, rec.Code, rec.Body.String())
	}
}
