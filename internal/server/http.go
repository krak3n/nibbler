package server

import (
	"context"
	"net/http"

	"github.com/krak3n/nibbler/internal/storage"
)

// An Option can configure a http server
type Option func(*http.Server)

// WithAddress configures a http server address
func WithAddress(addr string) Option {
	return func(s *http.Server) {
		s.Addr = addr
	}
}

// New constructs a new http server
func New(store storage.Store, opts ...Option) *http.Server {
	handler := NewHandler(store)
	mux := http.NewServeMux()
	mux.HandleFunc("/shorten", handler.Shorten)
	mux.HandleFunc("/reverse", handler.Reverse)
	srv := &http.Server{
		Handler: mux,
	}
	return srv
}

// Handler implents http Handler functions
type Handler struct {
	store storage.Store
}

// NewHandler constructs a new Handler
func NewHandler(store storage.Store) *Handler {
	return &Handler{
		store: store,
	}
}

// Shorten generates a short URL
func (h *Handler) Shorten(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	h.store.WriteURL(ctx, &storage.URL{})
}

// Reverse queries the store for a URL and if present reditects the user to the
// original URL
func (h *Handler) Reverse(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	h.store.ReadURL(ctx, "")
}
