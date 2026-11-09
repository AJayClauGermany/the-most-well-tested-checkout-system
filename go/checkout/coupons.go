package checkout

// Coupons (rules §4.4).

import "strings"

// coupon is what the shop knows about one coupon code.
type coupon struct {
	// discount gives the discount in cents for the price after promotions.
	discount func(afterPromotions int) int
	// minimumSpend is the least the customer must spend (after promotions) to use it.
	minimumSpend int
}

// coupons lists the coupons that exist, by their normalized code.
var coupons = map[string]coupon{
	"SAVE10":   {discount: func(afterPromotions int) int { return RoundHalfUp(afterPromotions, 10) }},
	"WELCOME5": {discount: func(int) int { return 500 }, minimumSpend: 2500},
}

// NormalizeCouponCode returns the coupon code as the shop knows it: trimmed and upper case.
// Missing, empty or blank all mean "no coupon" and give "".
func NormalizeCouponCode(code string) string {
	return strings.ToUpper(strings.TrimSpace(code))
}

// CouponDiscount returns the coupon discount in cents for afterPromotions cents of
// shopping. The error, if any, is a *CheckoutError with UNKNOWN_COUPON or
// COUPON_MINIMUM_NOT_MET (rules §4.1 rows 6 and 7).
func CouponDiscount(code string, afterPromotions int) (int, error) {
	normalized := NormalizeCouponCode(code)
	if normalized == "" {
		return 0, nil
	}

	known, exists := coupons[normalized]
	if !exists {
		return 0, newError(CodeUnknownCoupon, "Unknown coupon: %s", normalized)
	}
	if afterPromotions < known.minimumSpend {
		return 0, newError(CodeCouponMinimumNotMet,
			"%s needs a spend of at least %d cents, got %d", normalized, known.minimumSpend, afterPromotions)
	}
	return known.discount(afterPromotions), nil
}
