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

func TestGetPostBySlug_Success(t *testing.T) {
	cleanup(t)

	token := registerAndGetToken(t, "postauthor", "author@example.com", "password123")

	// Create a post
	reqBody := post.CreatePostRequest{
		Title:   "Test Post for Slug",
		Content: "This is the content.",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/posts", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	testServer.Mux().ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("Failed to create post: %s", rec.Body.String())
	}

	// Get the slug from list
	req = httptest.NewRequest(http.MethodGet, "/api/v1/posts", nil)
	rec = httptest.NewRecorder()
	testServer.Mux().ServeHTTP(rec, req)

	var listResponse server.APIResponse
	//nolint:errcheck
	json.Unmarshal(rec.Body.Bytes(), &listResponse)
	result := listResponse.Result.(map[string]any)
	posts := result["posts"].([]any)
	postData := posts[0].(map[string]any)
	slug := postData["slug"].(string)

	// Get post by slug (public endpoint)
	req = httptest.NewRequest(http.MethodGet, "/api/v1/posts/"+slug, nil)
	rec = httptest.NewRecorder()
	testServer.Mux().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d. Body: %s", http.StatusOK, rec.Code, rec.Body.String())
	}

	var response server.APIResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response.Message != "Post retrieved successfully" {
		t.Errorf("Expected message 'Post retrieved successfully', got '%s'", response.Message)
	}

	postResult := response.Result.(map[string]any)
	if postResult["title"] != "Test Post for Slug" {
		t.Errorf("Expected title 'Test Post for Slug', got '%v'", postResult["title"])
	}
}

func TestGetPostBySlug_NotFound(t *testing.T) {
	cleanup(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/posts/non-existent-slug", nil)
	rec := httptest.NewRecorder()
	testServer.Mux().ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d. Body: %s", http.StatusNotFound, rec.Code, rec.Body.String())
	}
}
