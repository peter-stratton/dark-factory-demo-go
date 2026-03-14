package store

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"

	"github.com/peter-stratton/dark-factory-demo-go/internal/model"
)

// Store is an in-memory bookmark store.
type Store struct {
	mu        sync.RWMutex
	bookmarks map[string]model.Bookmark
}

// New creates a new empty Store.
func New() *Store {
	return &Store{
		bookmarks: make(map[string]model.Bookmark),
	}
}

// List returns all bookmarks.
func (s *Store) List() []model.Bookmark {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]model.Bookmark, 0, len(s.bookmarks))
	for _, b := range s.bookmarks {
		result = append(result, b)
	}
	return result
}

// Get returns a bookmark by ID.
func (s *Store) Get(id string) (model.Bookmark, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	b, ok := s.bookmarks[id]
	if !ok {
		return model.Bookmark{}, model.ErrNotFound
	}
	return b, nil
}

// Create adds a new bookmark and returns it.
func (s *Store) Create(req model.CreateBookmarkRequest) model.Bookmark {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UTC()
	b := model.Bookmark{
		ID:        generateID(),
		URL:       req.URL,
		Title:     req.Title,
		Tags:      req.Tags,
		CreatedAt: now,
		UpdatedAt: now,
	}
	s.bookmarks[b.ID] = b
	return b
}

// Update modifies an existing bookmark.
func (s *Store) Update(id string, req model.UpdateBookmarkRequest) (model.Bookmark, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	b, ok := s.bookmarks[id]
	if !ok {
		return model.Bookmark{}, model.ErrNotFound
	}

	if req.URL != nil {
		b.URL = *req.URL
	}
	if req.Title != nil {
		b.Title = *req.Title
	}
	if req.Tags != nil {
		b.Tags = req.Tags
	}
	b.UpdatedAt = time.Now().UTC()
	s.bookmarks[id] = b
	return b, nil
}

// Delete removes a bookmark by ID.
func (s *Store) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.bookmarks[id]; !ok {
		return model.ErrNotFound
	}
	delete(s.bookmarks, id)
	return nil
}

func generateID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}
