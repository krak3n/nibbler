package psql

import (
	"context"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/krak3n/nibbler/internal/storage"

	// Initialise Postgres Driver
	_ "github.com/lib/pq"
)

// Store holds a database connection and fulfills the storage.Store interface
type Store struct {
	db *sqlx.DB
}

// New connects to the database and constructs a new Store
func New(dsn storage.DSN) (*Store, error) {
	db, err := sqlx.Connect("postgres", dsn.String())
	if err != nil {
		return nil, err
	}
	return &Store{
		db: db,
	}, nil
}

// WriteURL writes a URL to the database, in conflict of the URL the id of the
// already shortened URL will be returned
func (s *Store) WriteURL(ctx context.Context, url *storage.URL) error {
	return s.InTx(func(tx *sqlx.Tx) error {
		// Generate a new ID if one has not been provided
		if url.ID == "" {
			id, err := s.GenerateID(tx)
			if err != nil {
				return err
			}
			url.ID = id
		}

		// Insert Query that uses PSQL On Conflict to return the row if there is
		// a conflict on the URL - here the ID is not important, only the URL
		qry := `INSERT INTO urls (id, url)
		VALUES (:id, :url)
		ON CONFLICT (url) DO UPDATE SET url=EXCLUDED.url
		RETURNING id, url`

		// Preapre the query
		stmt, err := tx.PrepareNamed(qry)
		if err != nil {
			return err
		}

		// Defer closing statement
		defer func() {
			if err := stmt.Close(); err != nil {
				log.Println(err)
			}
		}()

		// Execute the query, expect a row to be returned
		res := stmt.QueryRow(url)
		if res.Err() != nil {
			return err
		}

		// Scan the row into the struct
		return res.StructScan(url)
	})
}

// ReadURL reads a URL from the Database by ID
func (s *Store) ReadURL(ctx context.Context, id string) (*storage.URL, error) {
	var url = new(storage.URL)

	err := s.InTx(func(tx *sqlx.Tx) error {
		// Prepare the query
		stmt, err := tx.PrepareNamed(`SELECT * FROM urls WHERE id=:id LIMIT 1`)
		if err != nil {
			return err
		}

		// Ensure we clean up the statement
		defer func() {
			if err := stmt.Close(); err != nil {
				log.Println(err)
			}
		}()

		// Execute the query, expect a row to be returned
		res := stmt.QueryRow(&storage.URL{ID: id})
		if res.Err() != nil {
			return err
		}

		// Scan the row into the struct
		return res.StructScan(url)
	})

	return url, err
}

// InTx runs a given function within a database transaction
func (s *Store) InTx(fn func(tx *sqlx.Tx) error) error {
	// Create a transaction
	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}

	// Execute the given function, handle rollback on error
	if err = fn(tx); err != nil {
		if rberr := tx.Rollback(); rberr != nil {
			log.Println("rollback error:", rberr)
		}
		return err
	}

	// Commit the transaction
	return tx.Commit()
}

// GenerateID generates a new ID
func (s *Store) GenerateID(tx *sqlx.Tx) (string, error) {
	for {
		stmt, err := tx.PrepareNamed(`SELECT COUNT(*) FROM urls WHERE id=:id`)
		if err != nil {
			return "", err
		}

		// Ensure we clean up the statement
		defer func() {
			if err := stmt.Close(); err != nil {
				log.Println(err)
			}
		}()

		id, err := storage.GenerateID(6)
		if err != nil {
			return "", err
		}

		var n int
		if err := stmt.Get(&n, map[string]interface{}{"id": id}); err != nil {
			return "", err
		}

		if n == 0 {
			return id, nil
		}
	}
}
