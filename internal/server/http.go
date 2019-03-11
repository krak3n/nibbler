package server

import (
	"context"
	"log"
	"net/http"

	"github.com/asaskevich/govalidator"
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
	router.HandleFunc("/shorten", handler.Shorten).Queries("url", "{url}")
	router.HandleFunc("/{id}", handler.Reverse)
	srv := &http.Server{
		Handler: router,
		// TODO: timeout options
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if !govalidator.IsURL(r.URL.Query().Get("url")) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid url"))
		return
	}

	url := &storage.URL{URL: r.URL.Query().Get("url")}
	if err := h.store.WriteURL(ctx, url); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error shortening url"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(url.ID))
}

// Reverse queries the store for a URL and if present reditects the user to the
// original URL
func (h *Handler) Reverse(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	vars := mux.Vars(r)

	url, err := h.store.ReadURL(ctx, vars["id"])
	if err != nil {
		// TODO: switch type of error, e.g no rows should 404
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error reversing url"))
	}

	http.Redirect(w, r, url.URL, http.StatusPermanentRedirect)
}
