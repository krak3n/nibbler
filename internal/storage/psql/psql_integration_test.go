//+build integration

package psql

import (
	"context"
	"testing"

	"github.com/krak3n/nibbler/internal/config"
	"github.com/krak3n/nibbler/internal/storage"
	"github.com/stretchr/testify/require"
)

// TestWriteURL_Idempotent ensures we can write the same ID/URL combination
// multiple times and not get a new row
func TestWriteURL_Idempotent(t *testing.T) {
	dsn := storage.DSN{
		Name:    config.DBName,
		User:    config.DBUser,
		Pass:    config.DBPassword,
		Host:    config.DBHost,
		SSLMode: config.DBSSLMode,
	}

	store, err := New(dsn)
	require.NoError(t, err)

	err = store.WriteURL(context.Background(), &storage.URL{
		ID:  "foo",
		URL: "http://foo.com",
	})
	require.NoError(t, err)

	err = store.WriteURL(context.Background(), &storage.URL{
		ID:  "foo",
		URL: "http://foo.com",
	})
	require.NoError(t, err)
}
