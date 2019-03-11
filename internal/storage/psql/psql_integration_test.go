//+build integration

package psql_test

import (
	"context"
	"testing"

	"github.com/krak3n/nibbler/internal/storage"
	"github.com/krak3n/nibbler/internal/storage/psql/psqltest"
	"github.com/stretchr/testify/require"
)

// TestWriteURL_CreateThenRead creates a new transaction and then reads it from
// the database
func TestWriteURL_CreateThenRead(t *testing.T) {
	ctx := context.Background()

	store := psqltest.NewStore(t)
	defer psqltest.Truncate(t, store)

	insertURL := &storage.URL{
		ID:  psqltest.GenerateID(t),
		URL: "http://foo.com",
	}

	require.NoError(t, store.WriteURL(ctx, insertURL))

	url, err := store.ReadURL(ctx, insertURL.ID)
	require.NoError(t, err)

	require.Equal(t, insertURL, url)
}

// TestWriteURL_Idempotent ensures we can write the same ID/URL combination
// multiple times and not get a new row
func TestWriteURL_IdempotentURL(t *testing.T) {
	ctx := context.Background()

	store := psqltest.NewStore(t)
	defer psqltest.Truncate(t, store)

	id := psqltest.GenerateID(t)
	expected := &storage.URL{
		ID:  id,
		URL: "http://foo.com",
	}

	insert1 := &storage.URL{
		ID:  id,
		URL: "http://foo.com",
	}

	require.NoError(t, store.WriteURL(ctx, insert1))
	require.Equal(t, expected, insert1)

	insert2 := &storage.URL{
		ID:  psqltest.GenerateID(t),
		URL: "http://foo.com",
	}

	require.NoError(t, store.WriteURL(ctx, insert2))
	require.Equal(t, expected, insert2)
}

func TestWriteURL_NewID(t *testing.T) {
	ctx := context.Background()

	store := psqltest.NewStore(t)
	defer psqltest.Truncate(t, store)

	url := &storage.URL{
		URL: "http://foo.com",
	}

	require.NoError(t, store.WriteURL(ctx, url))

	t.Log("url inserted with id:", url.ID)

	require.NotEmpty(t, url.ID)
	require.Equal(t, "http://foo.com", url.URL)
}
