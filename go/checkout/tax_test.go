package checkout_test

// UNIT TESTS for tax (rules §4.6).

import (
	"testing"

	"github.com/AJayClauGermany/the-most-well-tested-checkout-system/go/checkout"
)

func TestTaxIsEightPercent(t *testing.T) {
	t.Parallel()
	if checkout.TaxPercent != 8 {
		t.Errorf("TaxPercent = %d, want 8", checkout.TaxPercent)
	}
}

func TestTaxOn(t *testing.T) {
	t.Parallel()
	tests := []struct {
		amount int
		want   int
		why    string
	}{
		{amount: 50, want: 4, why: "8% of $0.50 is exactly 4 cents"},
		{amount: 120, want: 10, why: "9.6 cents rounds up"},
		{amount: 240, want: 19, why: "19.2 cents rounds down"},
		{amount: 625, want: 50, why: "exactly 50 cents"},
		{amount: 1706, want: 136, why: "136.48 rounds down"},
		{amount: 11880, want: 950, why: "950.4 rounds down"},
		{amount: 0, want: 0, why: "no spend, no tax"},
		// 8% of a whole number of cents is never exactly half a cent (x × 8 never ends in
		// 50), so the half-up rule is pinned down by the 10% coupon tests. Here: plain rounding.
		{amount: 1875, want: 150, why: "exactly 150.00"},
		{amount: 1881, want: 150, why: "150.48 rounds down"},
		{amount: 1888, want: 151, why: "151.04 rounds down to 151"},
	}

	for _, tt := range tests {
		t.Run(itoa(tt.amount), func(t *testing.T) {
			t.Parallel()
			if got := checkout.TaxOn(tt.amount); got != tt.want {
				t.Errorf("TaxOn(%d) = %d, want %d (%s)", tt.amount, got, tt.want, tt.why)
			}
		})
	}
}
