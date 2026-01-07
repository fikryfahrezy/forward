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

func TestDeleteComment_Success(t *testing.T) {
	cleanup(t)

	token := registerAndGetToken(t, "commenter", "commenter@example.com", "password123")
	postID := createPost(t, token, "Test Post", "Test content")

	createReq := comment.CreateCommentRequest{
		Content: "Comment to delete",
	}
	body, _ := json.Marshal(createReq)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/posts/"+postID+"/comments", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	testServer.Mux().ServeHTTP(rec, req)

	var createResponse server.APIResponse
	//nolint:errcheck
	json.Unmarshal(rec.Body.Bytes(), &createResponse)
	result := createResponse.Result.(map[string]any)
	commentID := result["id"].(string)

	req = httptest.NewRequest(http.MethodDelete, "/api/v1/comments/"+commentID, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec = httptest.NewRecorder()
	testServer.Mux().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d. Body: %s", http.StatusOK, rec.Code, rec.Body.String())
	}

	req = httptest.NewRequest(http.MethodGet, "/api/v1/posts/"+postID+"/comments", nil)
	rec = httptest.NewRecorder()
	testServer.Mux().ServeHTTP(rec, req)

	var listResponse server.APIResponse
	//nolint:errcheck
	json.Unmarshal(rec.Body.Bytes(), &listResponse)
	listResult := listResponse.Result.(map[string]any)
	comments := listResult["comments"].([]any)
	if len(comments) != 0 {
		t.Errorf("Expected 0 comments after deletion, got %d", len(comments))
	}
}

func TestDeleteComment_NotAuthor(t *testing.T) {
	cleanup(t)

	token1 := registerAndGetToken(t, "author1", "author1@example.com", "password123")
	postID := createPost(t, token1, "Test Post", "Test content")

	createReq := comment.CreateCommentRequest{
		Content: "Author1's comment",
	}
	body, _ := json.Marshal(createReq)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/posts/"+postID+"/comments", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token1)
	rec := httptest.NewRecorder()
	testServer.Mux().ServeHTTP(rec, req)

	var createResponse server.APIResponse
	//nolint:errcheck
	json.Unmarshal(rec.Body.Bytes(), &createResponse)
	result := createResponse.Result.(map[string]any)
	commentID := result["id"].(string)

	token2 := registerAndGetToken(t, "author2", "author2@example.com", "password123")

	req = httptest.NewRequest(http.MethodDelete, "/api/v1/comments/"+commentID, nil)
	req.Header.Set("Authorization", "Bearer "+token2)
	rec = httptest.NewRecorder()
	testServer.Mux().ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Errorf("Expected status %d, got %d. Body: %s", http.StatusForbidden, rec.Code, rec.Body.String())
	}
}

func TestDeleteComment_Unauthorized(t *testing.T) {
	cleanup(t)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/comments/some-id", nil)
	rec := httptest.NewRecorder()
	testServer.Mux().ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d. Body: %s", http.StatusUnauthorized, rec.Code, rec.Body.String())
	}
}

func TestDeleteComment_NotFound(t *testing.T) {
	cleanup(t)

	token := registerAndGetToken(t, "commenter", "commenter@example.com", "password123")

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/comments/550e8400-e29b-41d4-a716-446655440000", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	testServer.Mux().ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d. Body: %s", http.StatusNotFound, rec.Code, rec.Body.String())
	}
}
