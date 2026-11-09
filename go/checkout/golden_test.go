package checkout_test

// SNAPSHOT (golden file) TESTS: did the output change since last time?
//
// A golden file holds the approved output. The test produces the output again and
// compares it with the file, byte for byte. Go has no built-in snapshot tool, so this is
// the idiomatic hand-rolled version: files in testdata/ (a folder name the go tool
// ignores) and an -update flag that rewrites them when a change is intended:
//
//	go test ./checkout -update
//
// Always read the diff before updating: a changed golden file is a changed behaviour.
// The hand-checked receipt in spec/examples/receipt.json plays the role of an "approved"
// output shared by all six languages.

import (
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/AJayClauGermany/the-most-well-tested-checkout-system/go/checkout"
	"github.com/AJayClauGermany/the-most-well-tested-checkout-system/go/internal/fixtures"
	"github.com/AJayClauGermany/the-most-well-tested-checkout-system/go/internal/spec"
)

var update = flag.Bool("update", false, "rewrite the golden files in testdata/ with the current output")

// assertMatchesGolden compares got with testdata/<name>, or rewrites the file with -update.
func assertMatchesGolden(t *testing.T, name string, got []byte) {
	t.Helper()
	path := filepath.Join("testdata", name)
	if *update {
		if err := os.WriteFile(path, got, 0o644); err != nil {
			t.Fatalf("writing the golden file: %v", err)
		}
		return
	}
	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading the golden file: %v (run `go test ./checkout -update` to create it)", err)
	}
	if diff := cmp.Diff(string(want), string(got)); diff != "" {
		t.Errorf("%s does not match the golden file (-golden +got):\n%s\n"+
			"If the change is intended, run: go test ./checkout -update", name, diff)
	}
}

// prettyJSON renders a receipt the way the golden files store it: indented, newline at the end.
func prettyJSON(t *testing.T, value any) []byte {
	t.Helper()
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		t.Fatalf("encoding JSON: %v", err)
	}
	return append(data, '\n')
}

func TestTheExampleReceiptMatchesTheHandCheckedSpecExample(t *testing.T) {
	t.Parallel()
	dir, err := spec.Dir()
	if err != nil {
		t.Fatal(err)
	}
	var approved checkout.Receipt
	if err := spec.ReadJSON(filepath.Join(dir, "examples", "receipt.json"), &approved); err != nil {
		t.Fatal(err)
	}

	receipt, err := checkout.Checkout(fixtures.ExampleBasket, loadCatalog(t))

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if diff := cmp.Diff(approved, receipt); diff != "" {
		t.Errorf("Checkout() differs from spec/examples/receipt.json (-approved +got):\n%s", diff)
	}
}

func TestTheExampleReceiptMatchesItsGoldenFile(t *testing.T) {
	t.Parallel()
	receipt, err := checkout.Checkout(fixtures.ExampleBasket, loadCatalog(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertMatchesGolden(t, "example-receipt.golden.json", prettyJSON(t, receipt))
}

func TestTheBiggestSharedCaseMatchesItsGoldenFile(t *testing.T) {
	t.Parallel()
	// Every product, both promotions, WELCOME5 and express shipping.
	receipt, err := checkout.Checkout(fixtures.EverythingBasket, loadCatalog(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertMatchesGolden(t, "full-everything.golden.json", prettyJSON(t, receipt))
}

func TestOnePenMatchesItsGoldenFile(t *testing.T) {
	t.Parallel()
	receipt, err := checkout.Checkout(items(item("PEN", 1)), loadCatalog(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertMatchesGolden(t, "one-pen.golden.json", prettyJSON(t, receipt))
}
