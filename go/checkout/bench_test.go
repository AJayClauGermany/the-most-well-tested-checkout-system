package checkout_test

// A BENCHMARK: how fast is Checkout? Run with
//
//	go test ./checkout -run '^$' -bench . -benchmem
//
// (-run '^$' skips the tests so only benchmarks run.) b.Loop (Go 1.24) runs the body
// as many times as needed for a stable measurement and keeps the compiler from
// optimising the call away. The checkout is tiny, so one benchmark is enough to show the
// tool – this is not a performance project (see docs/testing-matrix.md).

import (
	"testing"

	"github.com/AJayClauGermany/the-most-well-tested-checkout-system/go/checkout"
	"github.com/AJayClauGermany/the-most-well-tested-checkout-system/go/internal/fixtures"
)

func BenchmarkCheckout(b *testing.B) {
	catalog := loadCatalog(b)
	oneApple := checkout.Basket{Items: []checkout.BasketItem{{SKU: "APPLE", Quantity: 1}}}

	b.Run("one apple", func(b *testing.B) {
		for b.Loop() {
			_, _ = checkout.Checkout(oneApple, catalog)
		}
	})

	b.Run("every product, both promotions, a coupon and express shipping", func(b *testing.B) {
		for b.Loop() {
			_, _ = checkout.Checkout(fixtures.EverythingBasket, catalog)
		}
	})
}
