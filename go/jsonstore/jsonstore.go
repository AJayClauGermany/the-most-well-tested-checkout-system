// Package jsonstore is an order repository that writes one JSON file per order into a
// folder: <folder>/ORDER-1.json, <folder>/ORDER-2.json, … in the shape of
// spec/schema/order.schema.json.
package jsonstore

import (
	"errors"

	"github.com/AJayClauGermany/the-most-well-tested-checkout-system/go/checkout"
)

var errNotImplemented = errors.New("jsonstore: not implemented yet")

// Repository keeps orders as JSON files in one folder.
type Repository struct{}

// Open returns a repository for folder. TODO: not implemented yet.
func Open(folder string, clock checkout.Clock) (*Repository, error) {
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
