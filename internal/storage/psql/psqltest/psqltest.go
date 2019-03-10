package psqltest

import (
	"testing"

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
