package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/peter-stratton/dark-factory-demo-go/internal/model"
	"github.com/peter-stratton/dark-factory-demo-go/internal/store"
)

// Server handles HTTP requests for the bookmarks API.
type Server struct {
	store *store.Store
}

// New creates a new Server.
func New(s *store.Store) *Server {
	return &Server{store: s}
}

// Router returns an http.Handler with all routes registered.
func (s *Server) Router() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /bookmarks", s.handleList)
	mux.HandleFunc("POST /bookmarks", s.handleCreate)
	mux.HandleFunc("GET /bookmarks/{id}", s.handleGet)
	mux.HandleFunc("PATCH /bookmarks/{id}", s.handleUpdate)
	mux.HandleFunc("GET /health", s.handleHealth)

	return mux
}

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleList(w http.ResponseWriter, _ *http.Request) {
	bookmarks := s.store.List()
	writeJSON(w, http.StatusOK, bookmarks)
}

func (s *Server) handleCreate(w http.ResponseWriter, r *http.Request) {
	var req model.CreateBookmarkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := req.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	bookmark := s.store.Create(req)
	writeJSON(w, http.StatusCreated, bookmark)
}

func (s *Server) handleGet(w http.ResponseWriter, r *http.Request) {
	id := extractID(r)

	bookmark, err := s.store.Get(id)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			writeError(w, http.StatusNotFound, "bookmark not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeJSON(w, http.StatusOK, bookmark)
}

func (s *Server) handleUpdate(w http.ResponseWriter, r *http.Request) {
	id := extractID(r)

	var req model.UpdateBookmarkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	bookmark, err := s.store.Update(id, req)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			writeError(w, http.StatusNotFound, "bookmark not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeJSON(w, http.StatusOK, bookmark)
}

func extractID(r *http.Request) string {
	// Go 1.22+ path value
	if id := r.PathValue("id"); id != "" {
		return id
	}
	// Fallback: extract from path
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) >= 3 {
		return parts[2]
	}
	return ""
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
