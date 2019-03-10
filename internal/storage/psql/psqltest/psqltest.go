package psqltest

import (
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/krak3n/nibbler/internal/config"
	"github.com/krak3n/nibbler/internal/storage"
	"github.com/krak3n/nibbler/internal/storage/psql"
	"github.com/stretchr/testify/require"
)

// NewStore creates a new postgres storage db
func NewStore(t *testing.T) *psql.Store {
	dsn := storage.DSN{
		Name:    config.DBName,
		User:    config.DBUser,
		Pass:    config.DBPassword,
		Host:    config.DBHost,
		SSLMode: config.DBSSLMode,
	}

	store, err := psql.New(dsn)
	require.NoError(t, err)

	return store
}

// GenerateID generates a New ID
func GenerateID(t *testing.T) string {
	id, err := storage.GenerateID(6)
	require.NoError(t, err)
	t.Log("Generated ID:", id)
	return id
}

// Truncate empties the database tables
func Truncate(t *testing.T, store *psql.Store) {
	require.NoError(t, store.InTx(func(tx *sqlx.Tx) error {
		_, err := tx.Exec(`TRUNCATE urls;`)
		require.NoError(t, err)

		return nil
	}))
}
