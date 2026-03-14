package model

import "time"

// Bookmark represents a saved URL with metadata.
type Bookmark struct {
	ID        string    `json:"id"`
	URL       string    `json:"url"`
	Title     string    `json:"title"`
	Tags      []string  `json:"tags,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateBookmarkRequest is the payload for creating a bookmark.
type CreateBookmarkRequest struct {
	URL   string   `json:"url"`
	Title string   `json:"title"`
	Tags  []string `json:"tags,omitempty"`
}

// UpdateBookmarkRequest is the payload for updating a bookmark.
type UpdateBookmarkRequest struct {
	URL   *string  `json:"url,omitempty"`
	Title *string  `json:"title,omitempty"`
	Tags  []string `json:"tags,omitempty"`
}

// Validate checks that required fields are present on a create request.
func (r CreateBookmarkRequest) Validate() error {
	if r.URL == "" {
		return ErrURLRequired
	}
	if r.Title == "" {
		return ErrTitleRequired
	}
	return nil
}
