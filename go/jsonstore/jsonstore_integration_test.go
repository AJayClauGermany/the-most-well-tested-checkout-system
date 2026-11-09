//go:build integration

// INTEGRATION TESTS for the JSON-file order repository: real files in a real temporary
// folder. Only compiled with the build tag:  go test -tags integration ./...
//
// Unlike a unit test, this touches the file system, so each test gets its own fresh folder
// from t.TempDir(), which Go deletes when the test ends. The stored files must follow
// spec/schema/order.schema.json – a contract test inside an integration test.
package jsonstore_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"sync"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/AJayClauGermany/the-most-well-tested-checkout-system/go/checkout"
	"github.com/AJayClauGermany/the-most-well-tested-checkout-system/go/internal/fixtures"
	"github.com/AJayClauGermany/the-most-well-tested-checkout-system/go/internal/spec"
	"github.com/AJayClauGermany/the-most-well-tested-checkout-system/go/jsonstore"
)

// The repository must fit the interface that PlaceOrder and the tests use. This line does
// not run anything: it only compiles if *jsonstore.Repository has Save and Find.
var _ checkout.OrderRepository = (*jsonstore.Repository)(nil)

// open opens a repository in folder with the fixed clock, or stops the test.
func open(t *testing.T, folder string) *jsonstore.Repository {
	t.Helper()
	repository, err := jsonstore.Open(folder, fixtures.FixedClock())
	if err != nil {
		t.Fatalf("Open(%s) returned an error: %v", folder, err)
	}
	return repository
}

// save saves the one-apple receipt, or stops the test.
func save(t *testing.T, repository *jsonstore.Repository) string {
	t.Helper()
	id, err := repository.Save(fixtures.OneAppleReceipt())
	if err != nil {
		t.Fatalf("Save() returned an error: %v", err)
	}
	return id
}

// writeFile puts a file into the folder, or stops the test.
func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestSaveWritesORDER1JSONOnTheFirstSave(t *testing.T) {
	t.Parallel()
	folder := t.TempDir()
	repository := open(t, folder)

	id := save(t, repository)

	if id != "ORDER-1" {
		t.Errorf("id = %q, want ORDER-1", id)
	}
	if _, err := os.Stat(filepath.Join(folder, "ORDER-1.json")); err != nil {
		t.Errorf("ORDER-1.json was not written: %v", err)
	}
}

func TestSaveWritesAFileThatFollowsOrderSchemaJSONExactly(t *testing.T) {
	t.Parallel()
	folder := t.TempDir()
	save(t, open(t, folder))

	var stored any
	if err := spec.ReadJSON(filepath.Join(folder, "ORDER-1.json"), &stored); err != nil {
		t.Fatal(err)
	}

	problems, err := spec.SchemaErrors("order.schema.json", stored)
	if err != nil {
		t.Fatal(err)
	}
	if len(problems) != 0 {
		t.Errorf("ORDER-1.json does not follow order.schema.json: %v", problems)
	}
	var order checkout.Order
	if err := spec.ReadJSON(filepath.Join(folder, "ORDER-1.json"), &order); err != nil {
		t.Fatal(err)
	}
	want := checkout.Order{ID: "ORDER-1", CreatedAt: fixtures.FixedTimestamp, Receipt: fixtures.OneAppleReceipt()}
	if diff := cmp.Diff(want, order); diff != "" {
		t.Errorf("stored order mismatch (-want +got):\n%s", diff)
	}
}

func TestSaveWritesHumanReadableJSON(t *testing.T) {
	t.Parallel()
	folder := t.TempDir()
	save(t, open(t, folder))

	text, err := os.ReadFile(filepath.Join(folder, "ORDER-1.json"))

	if err != nil {
		t.Fatal(err)
	}
	if !json.Valid(text) || !slices.Contains([]byte(text), '\n') || text[len(text)-1] != '\n' {
		t.Errorf("want indented JSON ending in a newline, got:\n%s", text)
	}
}

func TestFindReadsASavedOrderBackExactlyAsSaved(t *testing.T) {
	t.Parallel()
	repository := open(t, t.TempDir())
	receipt, err := checkout.Checkout(checkout.Basket{
		Items: []checkout.BasketItem{{SKU: "COFFEE", Quantity: 3}}, Coupon: "SAVE10",
	}, fixtures.TinyCatalog)
	if err != nil {
		t.Fatal(err)
	}
	id, _ := repository.Save(receipt)

	order, err := repository.Find(id)

	if err != nil {
		t.Fatalf("Find() returned an error: %v", err)
	}
	want := &checkout.Order{ID: "ORDER-1", CreatedAt: fixtures.FixedTimestamp, Receipt: receipt}
	if diff := cmp.Diff(want, order); diff != "" {
		t.Errorf("Find() mismatch (-want +got):\n%s", diff)
	}
}

func TestSaveNumbersOrdersSequentially(t *testing.T) {
	t.Parallel()
	folder := t.TempDir()
	repository := open(t, folder)

	ids := []string{save(t, repository), save(t, repository), save(t, repository)}

	if diff := cmp.Diff([]string{"ORDER-1", "ORDER-2", "ORDER-3"}, ids); diff != "" {
		t.Errorf("ids mismatch (-want +got):\n%s", diff)
	}
	entries, _ := os.ReadDir(folder)
	var names []string
	for _, entry := range entries {
		names = append(names, entry.Name())
	}
	if diff := cmp.Diff([]string{"ORDER-1.json", "ORDER-2.json", "ORDER-3.json"}, names); diff != "" {
		t.Errorf("files mismatch (-want +got):\n%s", diff)
	}
}

func TestNumberingContinuesWhenTheSameFolderIsOpenedAgain(t *testing.T) {
	t.Parallel()
	folder := t.TempDir()
	first := open(t, folder)
	save(t, first)
	save(t, first)

	second := open(t, folder)

	if id := save(t, second); id != "ORDER-3" {
		t.Errorf("id = %q, want ORDER-3", id)
	}
	if order, _ := second.Find("ORDER-1"); order == nil || order.ID != "ORDER-1" {
		t.Errorf("Find(ORDER-1) = %v, want the first order", order)
	}
}

func TestNumbersFollowTheHighestExistingOrderEvenWithGapsAndTwoDigitNumbers(t *testing.T) {
	t.Parallel()
	folder := t.TempDir()
	writeFile(t, filepath.Join(folder, "ORDER-3.json"), `{}`)
	writeFile(t, filepath.Join(folder, "ORDER-12.json"), `{}`)
	repository := open(t, folder)

	// Not ORDER-4: comparing the names alphabetically would wrongly say 3 > 12.
	if id := save(t, repository); id != "ORDER-13" {
		t.Errorf("id = %q, want ORDER-13", id)
	}
}

func TestFilesThatAreNotOrdersAreIgnored(t *testing.T) {
	t.Parallel()
	folder := t.TempDir()
	writeFile(t, filepath.Join(folder, "README.txt"), "not an order")
	writeFile(t, filepath.Join(folder, "ORDER-9.json.bak"), "not an order either")
	writeFile(t, filepath.Join(folder, "ORDER-9"), "no .json extension, so not an order")
	writeFile(t, filepath.Join(folder, "ORDER-09.json"), "a leading zero is not an order id")
	writeFile(t, filepath.Join(folder, "notes.json"), "{}")
	repository := open(t, folder)

	if id := save(t, repository); id != "ORDER-1" {
		t.Errorf("id = %q, want ORDER-1", id)
	}
}

func TestOpenCreatesTheFolderIfItDoesNotExistYet(t *testing.T) {
	t.Parallel()
	nested := filepath.Join(t.TempDir(), "orders", "2026")

	save(t, open(t, nested))

	if _, err := os.Stat(filepath.Join(nested, "ORDER-1.json")); err != nil {
		t.Errorf("ORDER-1.json was not written: %v", err)
	}
}

func TestOpenFailsWhenTheFolderCannotBeCreated(t *testing.T) {
	t.Parallel()
	file := filepath.Join(t.TempDir(), "a-file")
	writeFile(t, file, "in the way")

	if _, err := jsonstore.Open(filepath.Join(file, "orders"), nil); err == nil {
		t.Error("Open() inside a file returned no error")
	}
}

func TestFindAnswersNotFoundForIDsThatHaveNoFile(t *testing.T) {
	t.Parallel()
	repository := open(t, t.TempDir())
	save(t, repository)

	for _, id := range []string{"ORDER-2", "ORDER-0", "nonsense", ""} {
		order, err := repository.Find(id)
		if order != nil || err != nil {
			t.Errorf("Find(%q) = %v, %v; want nil, nil", id, order, err)
		}
	}
}

func TestFindDoesNotLetAnIDEscapeTheFolder(t *testing.T) {
	t.Parallel()
	folder := t.TempDir()
	writeFile(t, filepath.Join(folder, "secret.json"), `{"id":"x","created_at":"x","receipt":{}}`)
	repository := open(t, filepath.Join(folder, "orders"))

	for _, id := range []string{"../secret", "../secret.json"} {
		if order, err := repository.Find(id); order != nil || err != nil {
			t.Errorf("Find(%q) = %v, %v; want nil, nil", id, order, err)
		}
	}
}

func TestFindReportsAFileThatIsNotValidJSON(t *testing.T) {
	t.Parallel()
	folder := t.TempDir()
	writeFile(t, filepath.Join(folder, "ORDER-1.json"), "{ broken")
	repository := open(t, folder)

	if _, err := repository.Find("ORDER-1"); err == nil {
		t.Error("Find() of a broken file returned no error")
	}
}

func TestFindReportsAnOrderThatCannotBeRead(t *testing.T) {
	t.Parallel()
	folder := t.TempDir()
	// A folder where a file should be: it exists, but cannot be read as a file.
	if err := os.Mkdir(filepath.Join(folder, "ORDER-1.json"), 0o755); err != nil {
		t.Fatal(err)
	}
	repository := open(t, folder)

	if _, err := repository.Find("ORDER-1"); err == nil {
		t.Error("Find() of an unreadable order returned no error")
	}
}

func TestSaveFailsWhenTheFolderIsGone(t *testing.T) {
	t.Parallel()
	folder := filepath.Join(t.TempDir(), "orders")
	repository := open(t, folder)
	if err := os.RemoveAll(folder); err != nil {
		t.Fatal(err)
	}

	if _, err := repository.Save(fixtures.OneAppleReceipt()); err == nil {
		t.Error("Save() into a deleted folder returned no error")
	}
}

func TestSaveFailsWhenTheFileCannotBeWritten(t *testing.T) {
	t.Parallel()
	folder := t.TempDir()
	repository := open(t, folder)
	// A folder named like the next order file is in the way of writing it. (Folders are
	// not orders, so the next id is still ORDER-1.)
	if err := os.Mkdir(filepath.Join(folder, "ORDER-1.json"), 0o755); err != nil {
		t.Fatal(err)
	}

	if _, err := repository.Save(fixtures.OneAppleReceipt()); err == nil {
		t.Error("Save() over a folder returned no error")
	}
}

func TestEachOrderIsStampedWithTheClocksTimeAtTheMomentOfSaving(t *testing.T) {
	t.Parallel()
	now := fixtures.FixedTime
	repository, err := jsonstore.Open(t.TempDir(), checkout.ClockFunc(func() time.Time { return now }))
	if err != nil {
		t.Fatal(err)
	}

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
	repository, err := jsonstore.Open(t.TempDir(), nil)
	if err != nil {
		t.Fatal(err)
	}
	before := time.Now().UTC().Truncate(time.Second)

	id := save(t, repository)

	order, _ := repository.Find(id)
	createdAt, err := time.Parse(time.RFC3339, order.CreatedAt)
	if err != nil || createdAt.Before(before) || createdAt.After(time.Now()) {
		t.Errorf("created_at = %q, want the current time", order.CreatedAt)
	}
}

func TestSavingFromManyGoroutinesGivesEveryOrderItsOwnFile(t *testing.T) {
	t.Parallel()
	// A CONCURRENCY test (run it with -race): without the mutex, two goroutines could
	// both see ORDER-3 as the highest file and both write ORDER-4.
	repository := open(t, t.TempDir())
	const savers = 20

	var wg sync.WaitGroup
	for range savers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := repository.Save(fixtures.OneAppleReceipt()); err != nil {
				t.Errorf("Save() returned an error: %v", err)
			}
		}()
	}
	wg.Wait()

	if id := save(t, repository); id != checkout.OrderID(savers+1) {
		t.Errorf("the next id is %q, want %q", id, checkout.OrderID(savers+1))
	}
}
