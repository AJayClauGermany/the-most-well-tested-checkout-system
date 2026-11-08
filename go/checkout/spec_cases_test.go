// This file says `package checkout_test` (with _test), so it is a "black-box"
// test: it can only use the exported names, exactly like a real user of the
// package. It has to be, too: the spec loader imports checkout, so a test
// inside package checkout importing the loader would be an import cycle.
package checkout_test

import (
	"errors"
	"testing"

	"github.com/AJayClauGermany/the-most-well-tested-checkout-system/go/checkout"
	"github.com/AJayClauGermany/the-most-well-tested-checkout-system/go/internal/spec"
)

// TestSpecCases runs every shared case from spec/test-cases/*.json against
// Checkout. Each case becomes its own subtest named after its id, e.g.
//
//	go test ./checkout -run 'TestSpecCases/promo-coffee-three-bags-one-free' -v
//
// These are the "red tests waiting" of test-driven development: they describe
// what Checkout must do before the code exists, so for now they are skipped.
func TestSpecCases(t *testing.T) {
	catalog, err := spec.LoadCatalog()
	if err != nil {
		t.Fatalf("loading the catalog: %v", err)
	}
	cases, err := spec.LoadCases()
	if err != nil {
		t.Fatalf("loading the spec cases: %v", err)
	}

	for _, c := range cases {
		t.Run(c.ID, func(t *testing.T) {
			t.Skip("Checkout() is not implemented yet – remove this skip when you start implementing")

			receipt, err := checkout.Checkout(c.Basket, catalog)

			if c.ExpectedError != "" {
				checkError(t, c, err)
			} else {
				checkTotals(t, c, receipt, err)
			}
		})
	}
}

// checkError asserts that Checkout failed with the expected error code.
func checkError(t *testing.T, c spec.Case, err error) {
	t.Helper() // report failures at the caller's line, not in here

	var checkoutErr *checkout.CheckoutError
	if !errors.As(err, &checkoutErr) {
		t.Fatalf("%s\nwant error %s, got err = %v", c.Description, c.ExpectedError, err)
	}
	if checkoutErr.Code != c.ExpectedError {
		t.Errorf("%s\nerror code = %s, want %s", c.Description, checkoutErr.Code, c.ExpectedError)
	}
}

// checkTotals asserts that Checkout succeeded with the six expected totals.
func checkTotals(t *testing.T, c spec.Case, receipt checkout.Receipt, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("%s\nunexpected error: %v", c.Description, err)
	}
	got := spec.Totals{
		Subtotal:          receipt.Subtotal,
		PromotionDiscount: receipt.PromotionDiscount,
		CouponDiscount:    receipt.CouponDiscount,
		Shipping:          receipt.Shipping,
		Tax:               receipt.Tax,
		Total:             receipt.Total,
	}
	if got != *c.Expected {
		t.Errorf("%s\ngot  %+v\nwant %+v", c.Description, got, *c.Expected)
	}
}
