package storage

import (
	"context"
	"fmt"
)

// DSN is a database data source name, e.g it's connection URL
type DSN struct {
	Host    string
	User    string
	Pass    string
	Name    string
	SSLMode string
}

func (dsn DSN) String() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=%s",
		dsn.Host,
		dsn.User,
		dsn.Pass,
		dsn.Name,
		dsn.SSLMode)
}

// A Store can handle database operations
type Store interface {
	URLReadWriter
}

// A URLWriter can write new URLs into a Store
type URLWriter interface {
	WriteURL(context.Context, *URL) error
}

// A URLReader can read urls from the store based on their short id
type URLReader interface {
	ReadURL(context.Context, string) error
}

// A URLReadWriter can read and write URLs from the store
type URLReadWriter interface {
	URLWriter
	URLReader
}

// URL represents a client row
type URL struct {
	ID  string `db:"id"`
	URL string `db:"url"`
}
