//+build integration

package server

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gorilla/mux"
	"github.com/krak3n/nibbler/internal/storage"
	"github.com/krak3n/nibbler/internal/storage/psql/psqltest"
	"github.com/stretchr/testify/require"
)

func TestShorten(t *testing.T) {
	query := url.Values{}
	query.Add("url", "http://google.com")

	url := url.URL{
		Path:     "/shorten",
		RawQuery: query.Encode(),
	}
	t.Log(url.String())

	r := httptest.NewRequest(http.MethodGet, url.String(), nil)
	w := httptest.NewRecorder()

	store := psqltest.NewStore(t)
	defer psqltest.Truncate(t, store)

	handler := NewHandler(store)
	handler.Shorten(w, r)
	require.Equal(t, http.StatusOK, w.Code)

	b, err := ioutil.ReadAll(w.Body)
	require.NoError(t, err)

	u, err := store.ReadURL(context.Background(), string(b))
	require.NoError(t, err)
	require.Equal(t, string(b), u.ID)
	require.Equal(t, "http://google.com", u.URL)
}

func TestReverse(t *testing.T) {
	query := url.Values{}
	query.Add("url", "http://google.com")

	store := psqltest.NewStore(t)
	defer psqltest.Truncate(t, store)

	u := &storage.URL{
		URL: "http://google.com",
	}
	require.NoError(t, store.WriteURL(context.Background(), u))

	url := url.URL{
		Path: fmt.Sprintf("/%s", u.ID),
	}
	t.Log(url.String())

	r := mux.SetURLVars(httptest.NewRequest(http.MethodGet, url.String(), nil), map[string]string{
		"id": u.ID,
	})
	w := httptest.NewRecorder()

	handler := NewHandler(store)
	handler.Reverse(w, r)
	require.Equal(t, http.StatusPermanentRedirect, w.Code)
	require.Equal(t, w.Header().Get("Location"), "http://google.com")
}
