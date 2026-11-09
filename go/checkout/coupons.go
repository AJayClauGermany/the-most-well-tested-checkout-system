package checkout

// Coupons (rules §4.4).

// NormalizeCouponCode returns the coupon code as the shop knows it: trimmed and upper case.
//
// TODO: not implemented yet.
func NormalizeCouponCode(code string) string {
	return ""
}

// CouponDiscount returns the coupon discount in cents for afterPromotions cents of
// shopping, or a *CheckoutError (UNKNOWN_COUPON, COUPON_MINIMUM_NOT_MET).
//
// TODO: not implemented yet.
func CouponDiscount(code string, afterPromotions int) (int, error) {
	return 0, ErrNotImplemented
}
