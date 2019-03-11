package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
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
	router := mux.NewRouter()
	router.HandleFunc("/shorten", handler.Shorten)
	router.HandleFunc("/{id}", handler.Reverse)
	srv := &http.Server{
		Handler: router,
	}
	for _, opt := range opts {
		opt(srv)
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

	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()

	// h.store.WriteURL(ctx, &storage.URL{})

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("shorten")))
}

// Reverse queries the store for a URL and if present reditects the user to the
// original URL
func (h *Handler) Reverse(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()
	// h.store.ReadURL(ctx, "")

	vars := mux.Vars(r)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("reverse %s", vars["id"])))
}
