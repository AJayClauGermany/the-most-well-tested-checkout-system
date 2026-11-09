package checkout

// Saving orders (rules §6): the order type, the clock, the repository interface and
// the simplest repository there is – a list in memory.
//
// The two interfaces below are what tests replace with TEST DOUBLES. The mocks in
// internal/mocks are generated from this file with mockgen (`go generate ./...`).

//go:generate go tool mockgen -source=orders.go -destination=../internal/mocks/orders.go -package=mocks

import "time"

// Order is a saved order, exactly as spec/schema/order.schema.json describes it.
type Order struct {
	ID string `json:"id"`
	// CreatedAt is UTC, ISO 8601, whole seconds, trailing Z: "2026-01-15T10:30:00Z".
	CreatedAt string  `json:"created_at"`
	Receipt   Receipt `json:"receipt"`
}

// Clock answers "what time is it?" (rules §6: the time source must be replaceable).
type Clock interface {
	Now() time.Time
}

// ClockFunc lets a plain function be used as a Clock.
type ClockFunc func() time.Time

// Now calls the function.
func (f ClockFunc) Now() time.Time { return f() }

// SystemClock returns the real clock. TODO: not implemented yet.
func SystemClock() Clock { return ClockFunc(func() time.Time { return time.Time{} }) }

// FixedClock returns a clock that always answers t. TODO: not implemented yet.
func FixedClock(t time.Time) Clock { return ClockFunc(func() time.Time { return time.Time{} }) }

// FormatTimestamp formats a time the way created_at must look. TODO: not implemented yet.
func FormatTimestamp(t time.Time) string { return "" }

// OrderID returns the id of the n-th order: "ORDER-1", … TODO: not implemented yet.
func OrderID(n int) string { return "" }

// OrderNumber is the inverse of OrderID. TODO: not implemented yet.
func OrderNumber(id string) (int, bool) { return 0, false }

// OrderRepository is where orders are kept (rules §6).
type OrderRepository interface {
	// Save stores a receipt and returns the new order's id.
	Save(receipt Receipt) (string, error)
	// Find looks an order up. A nil *Order with a nil error means "not found".
	Find(id string) (*Order, error)
}

// InMemoryOrderRepository keeps orders in memory. It doubles as the hand-written FAKE.
type InMemoryOrderRepository struct{}

// NewInMemoryOrderRepository creates an empty repository. TODO: not implemented yet.
func NewInMemoryOrderRepository(clock Clock) *InMemoryOrderRepository {
	return &InMemoryOrderRepository{}
}

// Save implements OrderRepository. TODO: not implemented yet.
func (r *InMemoryOrderRepository) Save(receipt Receipt) (string, error) {
	return "", ErrNotImplemented
}

// Find implements OrderRepository. TODO: not implemented yet.
func (r *InMemoryOrderRepository) Find(id string) (*Order, error) {
	return nil, ErrNotImplemented
}

// Count returns how many orders are saved. TODO: not implemented yet.
func (r *InMemoryOrderRepository) Count() int { return -1 }

// PlacedOrder is what PlaceOrder gives back: the receipt and the id it was saved under.
type PlacedOrder struct {
	OrderID string
	Receipt Receipt
}

// PlaceOrder checks out, then saves the receipt as an order. TODO: not implemented yet.
func PlaceOrder(basket Basket, catalog Catalog, repository OrderRepository) (PlacedOrder, error) {
	return PlacedOrder{}, ErrNotImplemented
}
