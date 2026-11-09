package checkout_test

// A FUZZ TEST: property P10 ("validation is total") attacked with random bytes.
//
// `go test` runs the seed corpus below like ordinary tests. `go test -fuzz=FuzzCheckout`
// keeps mutating the seeds for as long as you let it (CI uses -fuzztime=10s), looking for
// an input that makes the code panic or fail an assertion. A failing input is saved under
// testdata/fuzz/FuzzCheckout/ and replayed by every later `go test` run.

import (
	"encoding/json"
	"errors"
	"slices"
	"testing"

	"github.com/AJayClauGermany/the-most-well-tested-checkout-system/go/checkout"
)

func FuzzCheckout(f *testing.F) {
	catalog := loadCatalog(f)

	// The seed corpus: every shared basket, plus a few odd ones. The fuzzer mutates these.
	for _, c := range loadCases(f) {
		data, err := json.Marshal(c.Basket)
		if err != nil {
			f.Fatal(err)
		}
		f.Add(data)
	}
	f.Add([]byte(`{"items":[{"sku":"APPLE","quantity":1}],"coupon":null}`))
	f.Add([]byte(`{"items":[{"sku":"\u0000","quantity":-99999999}],"shipping":"\n"}`))
	f.Add([]byte(`{"items":[]}`))
	f.Add([]byte(`not json at all`))

	f.Fuzz(func(t *testing.T, data []byte) {
		var basket checkout.Basket
		if err := json.Unmarshal(data, &basket); err != nil {
			return // not a basket at all: nothing to check (the CLI handles that case)
		}

		receipt, err := checkout.Checkout(basket, catalog) // must never panic

		if err != nil {
			var checkoutErr *checkout.CheckoutError
			if !errors.As(err, &checkoutErr) {
				t.Fatalf("err = %v, want a *CheckoutError", err)
			}
			if !slices.Contains(checkout.ErrorCodes, checkoutErr.Code) {
				t.Errorf("code %q is not one of the six error codes", checkoutErr.Code)
			}
			return
		}
		if len(receipt.Lines) == 0 {
			t.Error("a successful receipt has no lines")
		}
		if want := receipt.Subtotal - receipt.PromotionDiscount - receipt.CouponDiscount + receipt.Shipping + receipt.Tax; receipt.Total != want {
			t.Errorf("total = %d, want %d", receipt.Total, want)
		}
	})
}
