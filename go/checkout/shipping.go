package checkout

// Shipping (rules §4.5).

// ShippingMethod is a shipping method the shop knows: ShippingStandard or ShippingExpress.
type ShippingMethod string

// The shipping prices from rules §4.5, in cents.
const (
	StandardShippingCost = 499
	ExpressShippingCost  = 999
	// FreeShippingFrom is the discounted subtotal from which standard shipping is free.
	FreeShippingFrom = 5000
)

// ParseShippingMethod turns the basket's shipping field into a known method.
//
// TODO: not implemented yet.
func ParseShippingMethod(shipping string) (ShippingMethod, error) {
	return "", ErrNotImplemented
}

// ShippingCost returns the shipping cost in cents for this method and this discounted
// subtotal.
//
// TODO: not implemented yet.
func ShippingCost(method ShippingMethod, discountedSubtotal int) int {
	return 0
}
