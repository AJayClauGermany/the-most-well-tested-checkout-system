package checkout_test

// UNIT TESTS for merging basket items into receipt lines and pricing them (rules §4.2).
// Whole slices and structs are compared with go-cmp's cmp.Diff, which prints exactly
// which field differs instead of two long "%+v" dumps.

import (
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/AJayClauGermany/the-most-well-tested-checkout-system/go/checkout"
	"github.com/AJayClauGermany/the-most-well-tested-checkout-system/go/internal/fixtures"
)

// itoa turns a number into a string, for subtest names.
func itoa(n int) string { return strconv.Itoa(n) }

func TestMergeItemsKeepsDifferentSKUsApartInBasketOrder(t *testing.T) {
	t.Parallel()
	merged := checkout.MergeItems([]checkout.BasketItem{item("PEN", 2), item("APPLE", 1)})

	want := []checkout.BasketItem{item("PEN", 2), item("APPLE", 1)}
	if diff := cmp.Diff(want, merged); diff != "" {
		t.Errorf("MergeItems() mismatch (-want +got):\n%s", diff)
	}
}

func TestMergeItemsAddsUpTheQuantitiesOfTheSameSKU(t *testing.T) {
	t.Parallel()
	// The example from rules §4.2: [APPLE x1, PEN x2, APPLE x3] gives [APPLE x4, PEN x2].
	// APPLE keeps the place of its first appearance.
	merged := checkout.MergeItems([]checkout.BasketItem{item("APPLE", 1), item("PEN", 2), item("APPLE", 3)})

	want := []checkout.BasketItem{item("APPLE", 4), item("PEN", 2)}
	if diff := cmp.Diff(want, merged); diff != "" {
		t.Errorf("MergeItems() mismatch (-want +got):\n%s", diff)
	}
}

func TestMergeItemsMergesThreeLinesOfOneSKUIntoOne(t *testing.T) {
	t.Parallel()
	merged := checkout.MergeItems([]checkout.BasketItem{item("COFFEE", 1), item("COFFEE", 1), item("COFFEE", 1)})

	want := []checkout.BasketItem{item("COFFEE", 3)}
	if diff := cmp.Diff(want, merged); diff != "" {
		t.Errorf("MergeItems() mismatch (-want +got):\n%s", diff)
	}
}

func TestMergeItemsGivesAnEmptyListForAnEmptyList(t *testing.T) {
	t.Parallel()
	if got := checkout.MergeItems(nil); len(got) != 0 {
		t.Errorf("MergeItems(nil) = %v, want an empty list", got)
	}
}

func TestMergeItemsDoesNotChangeTheBasketsOwnItems(t *testing.T) {
	t.Parallel()
	original := []checkout.BasketItem{item("APPLE", 1), item("APPLE", 1)}

	checkout.MergeItems(original)

	want := []checkout.BasketItem{item("APPLE", 1), item("APPLE", 1)}
	if diff := cmp.Diff(want, original); diff != "" {
		t.Errorf("the basket's items were changed (-want +got):\n%s", diff)
	}
}

func TestPriceLineCopiesTheProductAndMultipliesForTheLineTotal(t *testing.T) {
	t.Parallel()
	line := checkout.PriceLine(fixtures.Apple, 3)

	want := checkout.ReceiptLine{SKU: "APPLE", Name: "Apple", Quantity: 3, UnitPrice: 50, LineTotal: 150, LineDiscount: 0}
	if diff := cmp.Diff(want, line); diff != "" {
		t.Errorf("PriceLine() mismatch (-want +got):\n%s", diff)
	}
}

func TestPriceLineIncludesThePromotionDiscount(t *testing.T) {
	t.Parallel()
	tests := []struct {
		product          checkout.Product
		quantity         int
		wantTotal        int
		wantLineDiscount int
	}{
		{product: fixtures.Coffee, quantity: 3, wantTotal: 2697, wantLineDiscount: 899},    // one of three free
		{product: fixtures.Notebook, quantity: 10, wantTotal: 3490, wantLineDiscount: 349}, // 10% off
		{product: fixtures.Coffee, quantity: 2, wantTotal: 1798, wantLineDiscount: 0},      // below the promotion
		{product: fixtures.Notebook, quantity: 9, wantTotal: 3141, wantLineDiscount: 0},    // one short of the minimum
	}

	for _, tt := range tests {
		t.Run(tt.product.SKU+"x"+itoa(tt.quantity), func(t *testing.T) {
			t.Parallel()
			line := checkout.PriceLine(tt.product, tt.quantity)

			if line.LineTotal != tt.wantTotal {
				t.Errorf("line_total = %d, want %d", line.LineTotal, tt.wantTotal)
			}
			if line.LineDiscount != tt.wantLineDiscount {
				t.Errorf("line_discount = %d, want %d", line.LineDiscount, tt.wantLineDiscount)
			}
		})
	}
}
