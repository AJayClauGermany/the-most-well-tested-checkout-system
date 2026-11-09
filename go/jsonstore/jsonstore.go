// Package jsonstore is an order repository that writes one JSON file per order into a
// folder: <folder>/ORDER-1.json, <folder>/ORDER-2.json, … Each file has the shape of
// spec/schema/order.schema.json.
//
// It is the simplest "real" store: no database, just files you can open in an editor.
// Its tests are INTEGRATION tests (they touch the real file system in a temporary
// folder), so they sit behind the build tag `integration`.
package jsonstore

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/AJayClauGermany/the-most-well-tested-checkout-system/go/checkout"
)

// Repository keeps orders as JSON files in one folder. Ids continue where the files in
// the folder left off, so reopening the same folder later keeps numbering correctly.
type Repository struct {
	mu     sync.Mutex
	folder string
	clock  checkout.Clock
}

// Open creates the folder if needed and returns a repository for it. A nil clock means the
// real one.
func Open(folder string, clock checkout.Clock) (*Repository, error) {
	if clock == nil {
		clock = checkout.SystemClock()
	}
	if err := os.MkdirAll(folder, 0o755); err != nil {
		return nil, err
	}
	return &Repository{folder: folder, clock: clock}, nil
}

// Save implements checkout.OrderRepository: it writes <folder>/ORDER-<n>.json.
func (r *Repository) Save(receipt checkout.Receipt) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	highest, err := r.highestOrderNumber()
	if err != nil {
		return "", err
	}
	order := checkout.Order{
		ID:        checkout.OrderID(highest + 1),
		CreatedAt: checkout.FormatTimestamp(r.clock.Now()),
		Receipt:   receipt,
	}
	// Indented, because humans read these files. An order holds only strings and ints, so
	// encoding it cannot fail.
	data, _ := json.MarshalIndent(order, "", "  ")
	if err := os.WriteFile(r.pathOf(order.ID), append(data, '\n'), 0o644); err != nil {
		return "", err
	}
	return order.ID, nil
}

// Find implements checkout.OrderRepository: nil, nil when there is no such order.
func (r *Repository) Find(id string) (*checkout.Order, error) {
	// Only accept real order ids: "../secret" must not read files outside the folder.
	if _, ok := checkout.OrderNumber(id); !ok {
		return nil, nil
	}
	path := r.pathOf(id)
	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var order checkout.Order
	if err := json.Unmarshal(data, &order); err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	return &order, nil
}

// pathOf returns the file an order id is stored in.
func (r *Repository) pathOf(id string) string {
	return filepath.Join(r.folder, id+".json")
}

// highestOrderNumber returns the highest order number among the ORDER-<n>.json files in the
// folder, or 0 when there are none. Other files, and folders, are ignored.
func (r *Repository) highestOrderNumber() (int, error) {
	entries, err := os.ReadDir(r.folder)
	if err != nil {
		return 0, err
	}
	highest := 0
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !strings.HasSuffix(name, ".json") {
			continue
		}
		if n, ok := checkout.OrderNumber(strings.TrimSuffix(name, ".json")); ok && n > highest {
			highest = n
		}
	}
	return highest, nil
}
