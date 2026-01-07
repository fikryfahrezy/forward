package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fikryfahrezy/forward/blog-api/internal/post"
	"github.com/fikryfahrezy/forward/blog-api/internal/server"
)

func TestListPosts_Success(t *testing.T) {
	cleanup(t)

	token := registerAndGetToken(t, "postauthor", "author@example.com", "password123")

	// Create a few posts
	for i := 1; i <= 3; i++ {
		reqBody := post.CreatePostRequest{
			Title:   fmt.Sprintf("Post %d", i),
			Content: fmt.Sprintf("Content for post %d", i),
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/posts", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		testServer.Mux().ServeHTTP(rec, req)

		if rec.Code != http.StatusCreated {
			t.Fatalf("Failed to create post %d: %s", i, rec.Body.String())
		}
	}

	// List posts (public endpoint)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/posts", nil)
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
	posts := result["posts"].([]any)
	if len(posts) != 3 {
		t.Errorf("Expected 3 posts, got %d", len(posts))
	}
}

func TestListPosts_Pagination(t *testing.T) {
	cleanup(t)

	token := registerAndGetToken(t, "postauthor", "author@example.com", "password123")

	// Create 5 posts
	for i := 1; i <= 5; i++ {
		reqBody := post.CreatePostRequest{
			Title:   fmt.Sprintf("Post %d", i),
			Content: fmt.Sprintf("Content for post %d", i),
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/posts", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		testServer.Mux().ServeHTTP(rec, req)
	}

	// Request page 1 with page_size 2
	req := httptest.NewRequest(http.MethodGet, "/api/v1/posts?page=1&page_size=2", nil)
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
	posts := result["posts"].([]any)
	if len(posts) != 2 {
		t.Errorf("Expected 2 posts, got %d", len(posts))
	}
	if result["total_count"].(float64) != 5 {
		t.Errorf("Expected total_count 5, got %v", result["total_count"])
	}
}
