// Package sqlitestore is an order repository backed by SQLite, through Go's standard
// database/sql package and the pure-Go driver modernc.org/sqlite.
package sqlitestore

import (
	"errors"

	// The blank import registers the driver under the name "sqlite" with database/sql.
	_ "modernc.org/sqlite"

	"github.com/AJayClauGermany/the-most-well-tested-checkout-system/go/checkout"
)

var errNotImplemented = errors.New("sqlitestore: not implemented yet")

// Repository keeps orders in one table of an SQLite database.
type Repository struct{}

// Open opens (or creates) the database at path. TODO: not implemented yet.
func Open(path string, clock checkout.Clock) (*Repository, error) {
	return nil, errNotImplemented
}

// Save implements checkout.OrderRepository. TODO: not implemented yet.
func (r *Repository) Save(receipt checkout.Receipt) (string, error) {
	return "", errNotImplemented
}

// Find implements checkout.OrderRepository. TODO: not implemented yet.
func (r *Repository) Find(id string) (*checkout.Order, error) {
	return nil, errNotImplemented
}

// Close closes the database. TODO: not implemented yet.
func (r *Repository) Close() error {
	return errNotImplemented
}
