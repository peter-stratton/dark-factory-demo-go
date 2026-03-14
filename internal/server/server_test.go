package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/peter-stratton/dark-factory-demo-go/internal/model"
	"github.com/peter-stratton/dark-factory-demo-go/internal/store"
)

func setupTestServer() (*Server, http.Handler) {
	s := store.New()
	srv := New(s)
	return srv, srv.Router()
}

func TestHealthEndpoint(t *testing.T) {
	_, handler := setupTestServer()

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}

func TestCreateBookmark(t *testing.T) {
	_, handler := setupTestServer()

	body := `{"url":"https://example.com","title":"Example"}`
	req := httptest.NewRequest("POST", "/bookmarks", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", w.Code)
	}

	var b model.Bookmark
	if err := json.NewDecoder(w.Body).Decode(&b); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if b.URL != "https://example.com" {
		t.Fatalf("expected URL https://example.com, got %s", b.URL)
	}
	if b.ID == "" {
		t.Fatal("expected non-empty ID")
	}
}

func TestCreateBookmarkValidation(t *testing.T) {
	_, handler := setupTestServer()

	body := `{"url":"","title":"Example"}`
	req := httptest.NewRequest("POST", "/bookmarks", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestListBookmarks(t *testing.T) {
	_, handler := setupTestServer()

	// Create two bookmarks
	for _, title := range []string{"A", "B"} {
		body := `{"url":"https://example.com","title":"` + title + `"}`
		req := httptest.NewRequest("POST", "/bookmarks", bytes.NewBufferString(body))
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
	}

	req := httptest.NewRequest("GET", "/bookmarks", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var bookmarks []model.Bookmark
	if err := json.NewDecoder(w.Body).Decode(&bookmarks); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(bookmarks) != 2 {
		t.Fatalf("expected 2 bookmarks, got %d", len(bookmarks))
	}
}

func TestGetBookmark(t *testing.T) {
	_, handler := setupTestServer()

	// Create a bookmark
	body := `{"url":"https://example.com","title":"Example"}`
	req := httptest.NewRequest("POST", "/bookmarks", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	var created model.Bookmark
	if err := json.NewDecoder(w.Body).Decode(&created); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	// Get it
	req = httptest.NewRequest("GET", "/bookmarks/"+created.ID, nil)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var got model.Bookmark
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if got.ID != created.ID {
		t.Fatalf("expected ID %s, got %s", created.ID, got.ID)
	}
}

func TestGetBookmarkNotFound(t *testing.T) {
	_, handler := setupTestServer()

	req := httptest.NewRequest("GET", "/bookmarks/nonexistent", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", w.Code)
	}
}

func TestDeleteBookmark(t *testing.T) {
	_, handler := setupTestServer()

	// Create a bookmark
	body := `{"url":"https://example.com","title":"Example"}`
	req := httptest.NewRequest("POST", "/bookmarks", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	var created model.Bookmark
	if err := json.NewDecoder(w.Body).Decode(&created); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	// Delete it
	req = httptest.NewRequest("DELETE", "/bookmarks/"+created.ID, nil)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", w.Code)
	}

	// Verify it's gone
	req = httptest.NewRequest("GET", "/bookmarks/"+created.ID, nil)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status 404 after delete, got %d", w.Code)
	}
}
