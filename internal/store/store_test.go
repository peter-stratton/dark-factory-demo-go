package store

import (
	"testing"

	"github.com/peter-stratton/dark-factory-demo-go/internal/model"
)

func TestCreateAndGet(t *testing.T) {
	s := New()

	b := s.Create(model.CreateBookmarkRequest{
		URL:   "https://example.com",
		Title: "Example",
		Tags:  []string{"test"},
	})

	if b.ID == "" {
		t.Fatal("expected non-empty ID")
	}
	if b.URL != "https://example.com" {
		t.Fatalf("expected URL https://example.com, got %s", b.URL)
	}

	got, err := s.Get(b.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Title != "Example" {
		t.Fatalf("expected title Example, got %s", got.Title)
	}
}

func TestGetNotFound(t *testing.T) {
	s := New()

	_, err := s.Get("nonexistent")
	if err != model.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestList(t *testing.T) {
	s := New()

	s.Create(model.CreateBookmarkRequest{URL: "https://a.com", Title: "A"})
	s.Create(model.CreateBookmarkRequest{URL: "https://b.com", Title: "B"})

	items := s.List()
	if len(items) != 2 {
		t.Fatalf("expected 2 bookmarks, got %d", len(items))
	}
}

func TestUpdate(t *testing.T) {
	s := New()

	b := s.Create(model.CreateBookmarkRequest{
		URL:   "https://example.com",
		Title: "Example",
	})

	newTitle := "Updated"
	updated, err := s.Update(b.ID, model.UpdateBookmarkRequest{
		Title: &newTitle,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.Title != "Updated" {
		t.Fatalf("expected title Updated, got %s", updated.Title)
	}
	if updated.URL != "https://example.com" {
		t.Fatal("URL should not have changed")
	}
}

func TestUpdateNotFound(t *testing.T) {
	s := New()

	_, err := s.Update("nonexistent", model.UpdateBookmarkRequest{})
	if err != model.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

