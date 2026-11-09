//go:build integration

// INTEGRATION TESTS for the SQLite order repository, with the pure-Go driver
// modernc.org/sqlite. Only compiled with the build tag:  go test -tags integration ./...
//
// A real database, but no server to install: SQLite is one file. Each test gets its own
// file in a folder from t.TempDir(), and t.Cleanup closes the database afterwards.
package sqlitestore_test

import (
	"database/sql"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/AJayClauGermany/the-most-well-tested-checkout-system/go/checkout"
	"github.com/AJayClauGermany/the-most-well-tested-checkout-system/go/internal/fixtures"
	"github.com/AJayClauGermany/the-most-well-tested-checkout-system/go/internal/spec"
	"github.com/AJayClauGermany/the-most-well-tested-checkout-system/go/sqlitestore"
)

// Compiles only if *sqlitestore.Repository fits the OrderRepository interface.
var _ checkout.OrderRepository = (*sqlitestore.Repository)(nil)

// open opens a repository at path with the given clock and closes it when the test ends.
func open(t *testing.T, path string, clock checkout.Clock) *sqlitestore.Repository {
	t.Helper()
	repository, err := sqlitestore.Open(path, clock)
	if err != nil {
		t.Fatalf("Open(%s) returned an error: %v", path, err)
	}
	t.Cleanup(func() { repository.Close() }) // closing twice is harmless
	return repository
}

// newDatabase returns the path of a fresh database file in a temporary folder.
func newDatabase(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "orders.db")
}

// save saves the one-apple receipt, or stops the test.
func save(t *testing.T, repository *sqlitestore.Repository) string {
	t.Helper()
	id, err := repository.Save(fixtures.OneAppleReceipt())
	if err != nil {
		t.Fatalf("Save() returned an error: %v", err)
	}
	return id
}

func TestSavesAndFindsAnOrder(t *testing.T) {
	t.Parallel()
	repository := open(t, newDatabase(t), fixtures.FixedClock())
	receipt, err := checkout.Checkout(checkout.Basket{
		Items: []checkout.BasketItem{{SKU: "COFFEE", Quantity: 3}}, Coupon: "SAVE10",
	}, fixtures.TinyCatalog)
	if err != nil {
		t.Fatal(err)
	}

	id, err := repository.Save(receipt)
	if err != nil {
		t.Fatal(err)
	}
	order, err := repository.Find(id)

	if err != nil {
		t.Fatalf("Find() returned an error: %v", err)
	}
	want := &checkout.Order{ID: "ORDER-1", CreatedAt: fixtures.FixedTimestamp, Receipt: receipt}
	if diff := cmp.Diff(want, order); diff != "" {
		t.Errorf("Find() mismatch (-want +got):\n%s", diff)
	}
}

func TestGivesBackAnOrderThatFollowsOrderSchemaJSON(t *testing.T) {
	t.Parallel()
	repository := open(t, newDatabase(t), fixtures.FixedClock())
	order, _ := repository.Find(save(t, repository))

	problems, err := spec.SchemaErrors("order.schema.json", order)

	if err != nil {
		t.Fatal(err)
	}
	if len(problems) != 0 {
		t.Errorf("the order does not follow order.schema.json: %v", problems)
	}
}

func TestNumbersOrdersSequentially(t *testing.T) {
	t.Parallel()
	repository := open(t, newDatabase(t), fixtures.FixedClock())

	ids := []string{save(t, repository), save(t, repository), save(t, repository)}

	if diff := cmp.Diff([]string{"ORDER-1", "ORDER-2", "ORDER-3"}, ids); diff != "" {
		t.Errorf("ids mismatch (-want +got):\n%s", diff)
	}
}

func TestNumberingContinuesAfterTheDatabaseIsClosedAndOpenedAgain(t *testing.T) {
	t.Parallel()
	path := newDatabase(t)
	first := open(t, path, fixtures.FixedClock())
	save(t, first)
	save(t, first)
	first.Close()

	second := open(t, path, fixtures.FixedClock())

	if id := save(t, second); id != "ORDER-3" {
		t.Errorf("id = %q, want ORDER-3", id)
	}
	if order, _ := second.Find("ORDER-2"); order == nil || order.ID != "ORDER-2" {
		t.Errorf("Find(ORDER-2) = %v, want the second order", order)
	}
}

func TestReallyWritesToTheDatabaseFileCheckedWithASeparateConnection(t *testing.T) {
	t.Parallel()
	path := newDatabase(t)
	repository := open(t, path, fixtures.FixedClock())
	save(t, repository)
	repository.Close()

	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM orders").Scan(&count); err != nil {
		t.Fatal(err)
	}

	if count != 1 {
		t.Errorf("the orders table has %d rows, want 1", count)
	}
}

func TestRefusesToSaveOrFindOnceItIsClosed(t *testing.T) {
	t.Parallel()
	repository := open(t, newDatabase(t), fixtures.FixedClock())
	repository.Close()

	if _, err := repository.Save(fixtures.OneAppleReceipt()); err == nil {
		t.Error("Save() after Close() returned no error")
	}
	if _, err := repository.Find("ORDER-1"); err == nil {
		t.Error("Find() after Close() returned no error")
	}
}

func TestWorksInMemoryToo(t *testing.T) {
	t.Parallel()
	// ":memory:" is SQLite's name for a throw-away database that only lives in memory.
	repository := open(t, ":memory:", fixtures.FixedClock())

	id := save(t, repository)

	order, err := repository.Find(id)
	if err != nil || order == nil {
		t.Fatalf("Find(%q) = %v, %v", id, order, err)
	}
	if diff := cmp.Diff(fixtures.OneAppleReceipt(), order.Receipt); diff != "" {
		t.Errorf("receipt mismatch (-want +got):\n%s", diff)
	}
}

func TestOpenFailsForAPathThatCannotBeADatabase(t *testing.T) {
	t.Parallel()
	// The folder itself is not a file SQLite can open.
	if _, err := sqlitestore.Open(t.TempDir(), nil); err == nil {
		t.Error("Open() of a folder returned no error")
	}
}

func TestFindReportsAReceiptThatIsNotValidJSON(t *testing.T) {
	t.Parallel()
	path := newDatabase(t)
	repository := open(t, path, fixtures.FixedClock())
	save(t, repository)
	repository.Close()
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec("UPDATE orders SET receipt = '{ broken'"); err != nil {
		t.Fatal(err)
	}
	db.Close()

	_, err = open(t, path, fixtures.FixedClock()).Find("ORDER-1")

	if err == nil {
		t.Error("Find() of a broken receipt returned no error")
	}
}

func TestAnswersNotFoundWithoutThrowing(t *testing.T) {
	t.Parallel()
	repository := open(t, newDatabase(t), fixtures.FixedClock())
	save(t, repository)

	for _, id := range []string{"ORDER-2", "ORDER-0", "ORDER-01", "ORDER-x", "1", "", "ORDER-1; DROP TABLE orders"} {
		order, err := repository.Find(id)
		if order != nil || err != nil {
			t.Errorf("Find(%q) = %v, %v; want nil, nil", id, order, err)
		}
	}
	if order, _ := repository.Find("ORDER-1"); order == nil { // and the table is still there
		t.Error("ORDER-1 is gone")
	}
}

func TestEachOrderIsStampedWithTheClocksTimeAtTheMomentOfSaving(t *testing.T) {
	t.Parallel()
	now := fixtures.FixedTime
	repository := open(t, newDatabase(t), checkout.ClockFunc(func() time.Time { return now }))

	first := save(t, repository)
	now = time.Date(2026, 1, 16, 9, 0, 0, 0, time.UTC)
	second := save(t, repository)

	firstOrder, _ := repository.Find(first)
	secondOrder, _ := repository.Find(second)
	if firstOrder.CreatedAt != "2026-01-15T10:30:00Z" || secondOrder.CreatedAt != "2026-01-16T09:00:00Z" {
		t.Errorf("created_at = %q and %q", firstOrder.CreatedAt, secondOrder.CreatedAt)
	}
}

func TestANilClockMeansTheRealClock(t *testing.T) {
	t.Parallel()
	repository := open(t, newDatabase(t), nil)
	before := time.Now().UTC().Truncate(time.Second)

	order, _ := repository.Find(save(t, repository))

	createdAt, err := time.Parse(time.RFC3339, order.CreatedAt)
	if err != nil || createdAt.Before(before) || createdAt.After(time.Now()) {
		t.Errorf("created_at = %q, want the current time", order.CreatedAt)
	}
}

func TestSavingFromManyGoroutinesGivesEveryOrderItsOwnNumber(t *testing.T) {
	t.Parallel()
	// A CONCURRENCY test (run it with -race): database/sql is safe to share between
	// goroutines, and SQLite's AUTOINCREMENT hands out every number exactly once.
	repository := open(t, newDatabase(t), fixtures.FixedClock())
	const savers = 20

	var wg sync.WaitGroup
	ids := make(chan string, savers)
	for range savers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			id, err := repository.Save(fixtures.OneAppleReceipt())
			if err != nil {
				t.Errorf("Save() returned an error: %v", err)
			}
			ids <- id
		}()
	}
	wg.Wait()
	close(ids)

	seen := map[string]bool{}
	for id := range ids {
		seen[id] = true
	}
	for n := 1; n <= savers; n++ {
		if !seen[checkout.OrderID(n)] {
			t.Errorf("id %s was never handed out", checkout.OrderID(n))
		}
	}
}
