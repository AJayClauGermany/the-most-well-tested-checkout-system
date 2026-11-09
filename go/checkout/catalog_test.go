package checkout

// This is an example UNIT TEST: it checks one small piece of code (Catalog.Find)
// on its own, with a tiny made-up catalog, and needs no files or network.
//
// Go testing basics, all from the standard library ("testing" package):
//   - Test files end in _test.go. `go build` ignores them; `go test` runs them.
//   - A test is a function named TestXxx that takes t *testing.T.
//   - There is no assert library: you compare with a normal `if` and call
//     t.Errorf (report and keep going) or t.Fatalf (report and stop this test).
//   - Messages follow the Go habit "got X, want Y" so failures are easy to read.
//
// This file says `package checkout` (not checkout_test), so it is a
// "white-box" test: it could also see unexported (lower-case) names.

import (
	"errors"
	"testing"
)

// A tiny catalog just for these tests. Small, made-up data keeps unit tests
// easy to read: you can see every product the test depends on.
var testCatalog = Catalog{Products: []Product{
	{SKU: "APPLE", Name: "Apple", UnitPrice: 50},
	{SKU: "PEN", Name: "Pen", UnitPrice: 120},
}}

func TestCatalogFindReturnsTheProductWithThatSKU(t *testing.T) {
	// Arrange: get everything ready (here: testCatalog above).

	// Act: call the code we are testing.
	product, found := testCatalog.Find("PEN")

	// Assert: check the result.
	if !found {
		t.Fatal("Find(\"PEN\") did not find the product")
	}
	if product.Name != "Pen" {
		t.Errorf("Find(\"PEN\") returned name %q, want %q", product.Name, "Pen")
	}
}

// A TABLE-DRIVEN TEST: the most common style of Go test. Each row of the table
// is one case, and t.Run turns every row into its own named "subtest", so the
// output shows exactly which row failed. Run a single row with:
//
//	go test ./checkout -run 'TestCatalogFindDoesNotFindUnknownSKUs/lower-case'
func TestCatalogFindDoesNotFindUnknownSKUs(t *testing.T) {
	tests := []struct {
		name string
		sku  string
	}{
		{name: "not-in-catalog", sku: "BANANA"},
		{name: "lower-case", sku: "apple"}, // SKUs are case sensitive (rules §4.1)
		{name: "empty", sku: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, found := testCatalog.Find(tt.sku)

			if found {
				t.Errorf("Find(%q) found a product, want not found", tt.sku)
			}
		})
	}
}

func TestRequireProductReturnsTheProductWhenItExists(t *testing.T) {
	product, err := RequireProduct(testCatalog, "APPLE")

	if err != nil {
		t.Fatalf("RequireProduct(\"APPLE\") returned an error: %v", err)
	}
	if product != testCatalog.Products[0] {
		t.Errorf("RequireProduct(\"APPLE\") = %+v, want %+v", product, testCatalog.Products[0])
	}
}

func TestRequireProductFailsWithUnknownProductWhenItDoesNot(t *testing.T) {
	_, err := RequireProduct(testCatalog, "BANANA")

	var checkoutErr *CheckoutError
	if !errors.As(err, &checkoutErr) {
		t.Fatalf("got err = %v, want a *CheckoutError", err)
	}
	if checkoutErr.Code != CodeUnknownProduct {
		t.Errorf("code = %s, want %s", checkoutErr.Code, CodeUnknownProduct)
	}
	// The message names the SKU, so it is useful to a human reading a log.
	if checkoutErr.Message != "Unknown product: BANANA" {
		t.Errorf("message = %q, want %q", checkoutErr.Message, "Unknown product: BANANA")
	}
}
