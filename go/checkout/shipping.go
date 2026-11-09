package checkout

// Shipping (rules §4.5).

// ShippingMethod is a shipping method the shop knows: ShippingStandard or ShippingExpress.
// Basket.Shipping is a plain string because a basket can hold a wrong value such as
// "overnight"; ParseShippingMethod turns it into this type or reports the error.
type ShippingMethod string

// The shipping prices from rules §4.5, in cents.
const (
	StandardShippingCost = 499
	ExpressShippingCost  = 999
	// FreeShippingFrom is the discounted subtotal from which standard shipping is free.
	FreeShippingFrom = 5000
)

// ParseShippingMethod turns the basket's shipping field into a known method.
// "" (the field was missing) means standard. Anything else that is not exactly
// "standard" or "express" is a *CheckoutError with UNKNOWN_SHIPPING_METHOD
// (rules §4.1 row 5).
func ParseShippingMethod(shipping string) (ShippingMethod, error) {
	switch shipping {
	case "":
		return ShippingStandard, nil
	case ShippingStandard, ShippingExpress:
		return ShippingMethod(shipping), nil
	}
	return "", newError(CodeUnknownShippingMethod, "Unknown shipping method: %s", shipping)
}

// ShippingCost returns the shipping cost in cents for this method and this discounted
// subtotal (the price after promotions and coupon).
func ShippingCost(method ShippingMethod, discountedSubtotal int) int {
	if method == ShippingExpress {
		return ExpressShippingCost
	}
	if discountedSubtotal >= FreeShippingFrom {
		return 0
	}
	return StandardShippingCost
}
