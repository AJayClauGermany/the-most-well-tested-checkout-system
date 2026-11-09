package checkout_test

// UNIT TESTS for the two kinds of promotion (rules §4.2), written as TABLE-DRIVEN tests –
// Go's native way of running the same check for many examples. Each row becomes a named
// subtest, so `go test -v` shows exactly which example failed, and `-run` can pick one:
//
//	go test ./checkout -run 'TestBuyXGetYFree/5_coffees' -v

import (
	"fmt"
	"testing"

	"github.com/AJayClauGermany/the-most-well-tested-checkout-system/go/checkout"
)

var (
	buyTwoGetOne     = &checkout.Promotion{Type: checkout.BuyXGetYFree, Buy: 2, Free: 1}
	tenPercentOffTen = &checkout.Promotion{Type: checkout.BulkDiscount, MinQuantity: 10, Percent: 10}
)

func TestPromotionDiscountIsZeroWithoutAPromotion(t *testing.T) {
	t.Parallel()
	if got := checkout.PromotionDiscount(nil, 50, 899); got != 0 {
		t.Errorf("PromotionDiscount(nil, 50, 899) = %d, want 0", got)
	}
}

func TestBuyXGetYFree(t *testing.T) {
	t.Parallel()
	// Coffee is buy 2 get 1 free at $8.99: for every complete group of 3, one is free.
	tests := []struct {
		quantity  int
		freeItems int
	}{
		{quantity: 1, freeItems: 0},
		{quantity: 2, freeItems: 0}, // one short of a group
		{quantity: 3, freeItems: 1}, // exactly one group
		{quantity: 4, freeItems: 1},
		{quantity: 5, freeItems: 1}, // one short of the second group
		{quantity: 6, freeItems: 2},
		{quantity: 99, freeItems: 33},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d_coffees", tt.quantity), func(t *testing.T) {
			t.Parallel()
			got := checkout.PromotionDiscount(buyTwoGetOne, tt.quantity, 899)

			if want := tt.freeItems * 899; got != want {
				t.Errorf("discount for %d coffees = %d, want %d (%d free)", tt.quantity, got, want, tt.freeItems)
			}
		})
	}
}

func TestBuyXGetYFreeUsesThePromotionsOwnNumbers(t *testing.T) {
	t.Parallel()
	buyOneGetOne := &checkout.Promotion{Type: checkout.BuyXGetYFree, Buy: 1, Free: 1}
	buyThreeGetTwo := &checkout.Promotion{Type: checkout.BuyXGetYFree, Buy: 3, Free: 2}

	tests := []struct {
		name      string
		promotion *checkout.Promotion
		quantity  int
		want      int
	}{
		{name: "buy 1 get 1: every second of 4 is free", promotion: buyOneGetOne, quantity: 4, want: 200},
		{name: "buy 1 get 1: the fifth item is not in a group", promotion: buyOneGetOne, quantity: 5, want: 200},
		{name: "buy 3 get 2: 4 is not a group of 5", promotion: buyThreeGetTwo, quantity: 4, want: 0},
		{name: "buy 3 get 2: 5 is one group, 2 free", promotion: buyThreeGetTwo, quantity: 5, want: 200},
		{name: "buy 3 get 2: 10 is two groups, 4 free", promotion: buyThreeGetTwo, quantity: 10, want: 400},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := checkout.PromotionDiscount(tt.promotion, tt.quantity, 100); got != tt.want {
				t.Errorf("got %d, want %d", got, tt.want)
			}
		})
	}
}

func TestBuyXGetYFreeMultipliesByTheUnitPrice(t *testing.T) {
	t.Parallel()
	if got := checkout.PromotionDiscount(buyTwoGetOne, 3, 100); got != 100 {
		t.Errorf("3 items at 100 = %d off, want 100", got)
	}
	if got := checkout.PromotionDiscount(buyTwoGetOne, 3, 1); got != 1 {
		t.Errorf("3 items at 1 = %d off, want 1", got)
	}
}

func TestBuyXGetYFreeWithABrokenPromotionGivesNoDiscountInsteadOfCrashing(t *testing.T) {
	t.Parallel()
	// The catalog schema does not allow buy: 0, free: 0 – but a program should not divide
	// by zero just because its data is odd.
	broken := &checkout.Promotion{Type: checkout.BuyXGetYFree}

	if got := checkout.PromotionDiscount(broken, 5, 100); got != 0 {
		t.Errorf("got %d, want 0", got)
	}
}

func TestBulkDiscount(t *testing.T) {
	t.Parallel()
	// Notebooks are 10% off from 10 at $3.49.
	tests := []struct {
		quantity int
		want     int
		why      string
	}{
		{quantity: 9, want: 0, why: "one short of the minimum"},
		{quantity: 10, want: 349, why: "exactly the minimum: 10% of $34.90"},
		{quantity: 11, want: 384, why: "383.9 cents rounds up"},
		{quantity: 15, want: 524, why: "exactly 523.5 cents rounds UP"},
		{quantity: 99, want: 3455, why: "10% of $345.51 is 3455.1"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d_notebooks", tt.quantity), func(t *testing.T) {
			t.Parallel()
			if got := checkout.PromotionDiscount(tenPercentOffTen, tt.quantity, 349); got != tt.want {
				t.Errorf("discount for %d notebooks = %d, want %d (%s)", tt.quantity, got, tt.want, tt.why)
			}
		})
	}
}

func TestBulkDiscountUsesThePromotionsOwnMinimumAndPercentage(t *testing.T) {
	t.Parallel()
	quarterOffFromThree := &checkout.Promotion{Type: checkout.BulkDiscount, MinQuantity: 3, Percent: 25}

	tests := []struct {
		quantity  int
		unitPrice int
		want      int
	}{
		{quantity: 2, unitPrice: 100, want: 0},  // below the minimum
		{quantity: 3, unitPrice: 100, want: 75}, // 25% of 300
		{quantity: 3, unitPrice: 101, want: 76}, // 75.75 rounds up
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d_at_%d", tt.quantity, tt.unitPrice), func(t *testing.T) {
			t.Parallel()
			if got := checkout.PromotionDiscount(quarterOffFromThree, tt.quantity, tt.unitPrice); got != tt.want {
				t.Errorf("got %d, want %d", got, tt.want)
			}
		})
	}
}

func TestAnUnknownPromotionTypeGivesNoDiscount(t *testing.T) {
	t.Parallel()
	mystery := &checkout.Promotion{Type: "mystery_deal", Percent: 50}

	if got := checkout.PromotionDiscount(mystery, 10, 100); got != 0 {
		t.Errorf("got %d, want 0", got)
	}
}
