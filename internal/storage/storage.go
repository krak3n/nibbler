package storage

import (
	"context"
	"fmt"
)

// UrlsTable stores shortened URLs
var UrlsTable = "urls"

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
	URLWriter
}

// A URLWriter can write new URLs into a Store
type URLWriter interface {
	WirteUrl(context.Context, *URL) error
}

// URL represents a client row
type URL struct {
	ID  string `db:"id"`
	URL string `db:"url"`
}
