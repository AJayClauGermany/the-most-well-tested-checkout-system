// Package sqlitestore is an order repository backed by SQLite, through Go's standard
// database/sql package and the pure-Go driver modernc.org/sqlite (no C compiler needed).
//
// A real database, but no server to install: SQLite is one file. Its tests are
// INTEGRATION tests (a real database file in a temporary folder), so they sit behind the
// build tag `integration`.
package sqlitestore

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	// The blank import registers the driver under the name "sqlite" with database/sql.
	_ "modernc.org/sqlite"

	"github.com/AJayClauGermany/the-most-well-tested-checkout-system/go/checkout"
)

// One table. The row number becomes the order number, so ids are sequential and survive
// closing and reopening the file.
const createTable = `
CREATE TABLE IF NOT EXISTS orders (
	number     INTEGER PRIMARY KEY AUTOINCREMENT,
	created_at TEXT NOT NULL,
	receipt    TEXT NOT NULL
)`

// Repository keeps orders in one table of an SQLite database file (or ":memory:").
type Repository struct {
	db    *sql.DB
	clock checkout.Clock
}

// Open opens (or creates) the database at path and makes sure the table exists. A nil
// clock means the real one. Call Close when you are done.
func Open(path string, clock checkout.Clock) (*Repository, error) {
	if clock == nil {
		clock = checkout.SystemClock()
	}
	// sql.Open only checks the driver name ("sqlite" is registered by the import above);
	// the file itself is opened by the first query, so its errors show up below.
	db, _ := sql.Open("sqlite", path)
	// One connection at a time: an SQLite file only takes one writer anyway, and a
	// ":memory:" database would otherwise be a different, empty database per connection.
	db.SetMaxOpenConns(1)
	if _, err := db.Exec(createTable); err != nil {
		db.Close()
		return nil, err
	}
	return &Repository{db: db, clock: clock}, nil
}

// Save implements checkout.OrderRepository. The receipt is stored as JSON text.
func (r *Repository) Save(receipt checkout.Receipt) (string, error) {
	// A receipt holds only strings and ints, so encoding it cannot fail.
	data, _ := json.Marshal(receipt)

	// "?" placeholders keep the values apart from the SQL, so a receipt can never be
	// mistaken for a command (no SQL injection). RETURNING hands back the new row number.
	var number int
	err := r.db.QueryRow("INSERT INTO orders (created_at, receipt) VALUES (?, ?) RETURNING number",
		checkout.FormatTimestamp(r.clock.Now()), string(data)).Scan(&number)
	if err != nil {
		return "", err
	}
	return checkout.OrderID(number), nil
}

// Find implements checkout.OrderRepository: nil, nil when there is no such order.
func (r *Repository) Find(id string) (*checkout.Order, error) {
	number, ok := checkout.OrderNumber(id)
	if !ok {
		return nil, nil
	}
	var createdAt, data string
	err := r.db.QueryRow("SELECT created_at, receipt FROM orders WHERE number = ?", number).
		Scan(&createdAt, &data)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var receipt checkout.Receipt
	if err := json.Unmarshal([]byte(data), &receipt); err != nil {
		return nil, fmt.Errorf("order %s: %w", id, err)
	}
	return &checkout.Order{ID: id, CreatedAt: createdAt, Receipt: receipt}, nil
}

// Close closes the database. Save and Find fail afterwards instead of silently doing nothing.
func (r *Repository) Close() error {
	return r.db.Close()
}
