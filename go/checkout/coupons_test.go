package checkout_test

// UNIT TESTS for coupons (rules §4.4). The coupon minimum (2500) is tested on both sides.

import (
	"fmt"
	"testing"

	"github.com/AJayClauGermany/the-most-well-tested-checkout-system/go/checkout"
)

func TestNormalizeCouponCode(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input string
		want  string
	}{
		{input: "SAVE10", want: "SAVE10"},
		{input: "save10", want: "SAVE10"},
		{input: "  save10 ", want: "SAVE10"},
		{input: "\tWelcome5\n", want: "WELCOME5"},
		{input: "nope", want: "NOPE"},
		// Missing, empty or blank all mean "no coupon".
		{input: "", want: ""},
		{input: "   ", want: ""},
		{input: "\t\n", want: ""},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%q", tt.input), func(t *testing.T) {
			t.Parallel()
			if got := checkout.NormalizeCouponCode(tt.input); got != tt.want {
				t.Errorf("NormalizeCouponCode(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

// discountOf is CouponDiscount for a coupon that must be accepted.
func discountOf(t *testing.T, code string, afterPromotions int) int {
	t.Helper()
	discount, err := checkout.CouponDiscount(code, afterPromotions)
	if err != nil {
		t.Fatalf("CouponDiscount(%q, %d) returned an error: %v", code, afterPromotions, err)
	}
	return discount
}

func TestCouponDiscountIsZeroWithoutACoupon(t *testing.T) {
	t.Parallel()
	for _, code := range []string{"", "  ", "\t"} {
		if got := discountOf(t, code, 10000); got != 0 {
			t.Errorf("CouponDiscount(%q, 10000) = %d, want 0", code, got)
		}
	}
}

func TestSAVE10TakesTenPercentOffRoundedHalfUp(t *testing.T) {
	t.Parallel()
	tests := []struct {
		afterPromotions int
		want            int
	}{
		{afterPromotions: 500, want: 50},
		{afterPromotions: 225, want: 23},   // 22.5 → 23
		{afterPromotions: 224, want: 22},   // 22.4 → 22
		{afterPromotions: 1798, want: 180}, // 179.8 → 180
		{afterPromotions: 1, want: 0},      // 0.1 → 0: no minimum spend
		{afterPromotions: 5, want: 1},      // 0.5 → 1
	}

	for _, tt := range tests {
		t.Run(itoa(tt.afterPromotions), func(t *testing.T) {
			t.Parallel()
			if got := discountOf(t, "SAVE10", tt.afterPromotions); got != tt.want {
				t.Errorf("SAVE10 on %d = %d, want %d", tt.afterPromotions, got, tt.want)
			}
		})
	}
}

func TestWELCOME5TakesExactlyFiveDollarsOff(t *testing.T) {
	t.Parallel()
	if got := discountOf(t, "WELCOME5", 2500); got != 500 {
		t.Errorf("WELCOME5 on 2500 = %d, want 500", got)
	}
	if got := discountOf(t, "WELCOME5", 99999); got != 500 {
		t.Errorf("WELCOME5 on 99999 = %d, want 500", got)
	}
}

func TestWELCOME5WorksAtExactlyTheMinimumAndFailsOneCentBelow(t *testing.T) {
	t.Parallel()
	if got := discountOf(t, "WELCOME5", 2500); got != 500 {
		t.Errorf("WELCOME5 at exactly the minimum = %d, want 500", got)
	}

	_, err := checkout.CouponDiscount("WELCOME5", 2499)

	checkoutErr := assertErrorCode(t, err, checkout.CodeCouponMinimumNotMet)
	assertMessage(t, checkoutErr, "WELCOME5 needs a spend of at least 2500 cents, got 2499")
}

func TestCouponsIgnoreCaseAndSurroundingSpaces(t *testing.T) {
	t.Parallel()
	if got := discountOf(t, " save10 ", 500); got != 50 {
		t.Errorf("' save10 ' on 500 = %d, want 50", got)
	}
	if got := discountOf(t, "welcome5", 2500); got != 500 {
		t.Errorf("'welcome5' on 2500 = %d, want 500", got)
	}
}

func TestAnUnknownCouponIsRejected(t *testing.T) {
	t.Parallel()
	_, err := checkout.CouponDiscount("freestuff", 10000)

	checkoutErr := assertErrorCode(t, err, checkout.CodeUnknownCoupon)
	// The message shows the normalized code: that is the code the shop looked for.
	assertMessage(t, checkoutErr, "Unknown coupon: FREESTUFF")
}

func TestACodeThatMerelyStartsOrEndsWithARealCodeIsUnknown(t *testing.T) {
	t.Parallel()
	for _, code := range []string{"SAVE100", "WELCOME50", "SAVE1", "XSAVE10"} {
		_, err := checkout.CouponDiscount(code, 10000)

		assertErrorCode(t, err, checkout.CodeUnknownCoupon)
	}
}
