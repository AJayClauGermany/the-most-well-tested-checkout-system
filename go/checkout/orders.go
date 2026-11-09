package checkout

// Saving orders (rules §6): the order type, the clock, the repository interface and
// the simplest repository there is – a list in memory.
//
// The two interfaces below are what tests replace with TEST DOUBLES: a fixed clock
// instead of the real one, a fake or a mock repository instead of a real store.
// The mocks in internal/mocks are generated from this file with mockgen
// (`go generate ./...` re-creates them – see the go:generate line below).

//go:generate go tool mockgen -source=orders.go -destination=../internal/mocks/orders.go -package=mocks

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"sync"
	"time"
)

// Order is a saved order, exactly as spec/schema/order.schema.json describes it.
type Order struct {
	ID string `json:"id"`
	// CreatedAt is UTC, ISO 8601, whole seconds, trailing Z: "2026-01-15T10:30:00Z".
	CreatedAt string  `json:"created_at"`
	Receipt   Receipt `json:"receipt"`
}

// Clock answers "what time is it?". Production uses SystemClock; tests pass one that
// always says the same time – the rules (§6) demand that the time source can be
// replaced, because created_at is the one value that changes from run to run.
type Clock interface {
	Now() time.Time
}

// ClockFunc lets a plain function be used as a Clock (the same trick as http.HandlerFunc):
//
//	clock := checkout.ClockFunc(func() time.Time { return someTime })
type ClockFunc func() time.Time

// Now calls the function.
func (f ClockFunc) Now() time.Time { return f() }

// SystemClock returns the real clock.
func SystemClock() Clock { return ClockFunc(time.Now) }

// FixedClock returns a clock that always answers t – a STUB, in test-double terms.
func FixedClock(t time.Time) Clock {
	return ClockFunc(func() time.Time { return t })
}

// timestampLayout is Go's way of spelling "2026-01-15T10:30:00Z" (it is an example date,
// not a format string with letters). Format drops anything below whole seconds.
const timestampLayout = "2006-01-02T15:04:05Z"

// FormatTimestamp formats a time the way created_at must look: whole seconds, UTC, trailing Z.
func FormatTimestamp(t time.Time) string {
	return t.UTC().Format(timestampLayout)
}

const orderIDPrefix = "ORDER-"

// OrderID returns the id of the n-th order saved into a repository: "ORDER-1", "ORDER-2", …
func OrderID(n int) string {
	return orderIDPrefix + strconv.Itoa(n)
}

// orderIDPattern is a positive whole number with no leading zero, exactly as
// order.schema.json spells it.
var orderIDPattern = regexp.MustCompile(`^ORDER-([1-9][0-9]*)$`)

// OrderNumber is the inverse of OrderID: "ORDER-7" gives 7, true. Anything that is not an
// order id gives 0, false.
func OrderNumber(id string) (int, bool) {
	match := orderIDPattern.FindStringSubmatch(id)
	if match == nil {
		return 0, false
	}
	n, err := strconv.Atoi(match[1])
	if err != nil { // only for numbers too big for an int
		return 0, false
	}
	return n, true
}

// OrderRepository is where orders are kept (rules §6). It is implemented in memory
// (below), as JSON files (package jsonstore) and in SQLite (package sqlitestore).
type OrderRepository interface {
	// Save stores a receipt and returns the new order's id.
	Save(receipt Receipt) (string, error)
	// Find looks an order up. A nil *Order with a nil error means "not found" – a normal
	// answer, not an error. The error is for real trouble (a file that cannot be read, …).
	Find(id string) (*Order, error)
}

// InMemoryOrderRepository keeps orders in a plain slice. It doubles as the hand-written
// FAKE for tests. A mutex makes it safe to use from several goroutines at once.
type InMemoryOrderRepository struct {
	mu     sync.Mutex
	clock  Clock
	orders []Order
}

// NewInMemoryOrderRepository creates an empty repository. A nil clock means the real one.
func NewInMemoryOrderRepository(clock Clock) *InMemoryOrderRepository {
	if clock == nil {
		clock = SystemClock()
	}
	return &InMemoryOrderRepository{clock: clock}
}

// Save implements OrderRepository. It never fails.
func (r *InMemoryOrderRepository) Save(receipt Receipt) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := OrderID(len(r.orders) + 1)
	r.orders = append(r.orders, Order{
		ID:        id,
		CreatedAt: FormatTimestamp(r.clock.Now()),
		// A copy, so that whoever gave us the receipt cannot change the saved order afterwards.
		Receipt: cloneReceipt(receipt),
	})
	return id, nil
}

// Find implements OrderRepository.
func (r *InMemoryOrderRepository) Find(id string) (*Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, order := range r.orders {
		if order.ID == id {
			order.Receipt = cloneReceipt(order.Receipt) // the caller gets its own copy too
			return &order, nil
		}
	}
	return nil, nil
}

// Count returns how many orders are saved. Handy for tests ("the order repository is empty").
func (r *InMemoryOrderRepository) Count() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.orders)
}

// cloneReceipt copies a receipt deeply enough that changing one does not change the other
// (the Lines slice is the only part that would otherwise be shared).
func cloneReceipt(receipt Receipt) Receipt {
	receipt.Lines = slices.Clone(receipt.Lines)
	return receipt
}

// PlacedOrder is what PlaceOrder gives back: the receipt and the id it was saved under.
type PlacedOrder struct {
	OrderID string
	Receipt Receipt
}

// PlaceOrder checks out, then saves the receipt as an order. A failed checkout saves
// nothing: the *CheckoutError is returned before the repository is touched.
func PlaceOrder(basket Basket, catalog Catalog, repository OrderRepository) (PlacedOrder, error) {
	receipt, err := Checkout(basket, catalog)
	if err != nil {
		return PlacedOrder{}, err
	}
	id, err := repository.Save(receipt)
	if err != nil {
		return PlacedOrder{}, fmt.Errorf("saving the order: %w", err)
	}
	return PlacedOrder{OrderID: id, Receipt: receipt}, nil
}
