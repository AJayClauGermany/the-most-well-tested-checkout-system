package checkout

// Tax (rules §4.6).

// TaxPercent is the tax rate: 8% of the discounted subtotal.
const TaxPercent = 8

// TaxOn returns 8% of discountedSubtotal, rounded half up. Shipping is never passed
// in here – it is not taxed.
func TaxOn(discountedSubtotal int) int {
	return RoundHalfUp(discountedSubtotal, TaxPercent)
}
