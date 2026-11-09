package checkout

// Basket validation (rules §4.1, rows 1 to 4). Rows 5 to 7 (shipping and coupon) are
// checked by shipping.go and coupons.go, because the coupon check needs the price after
// promotions.

// The quantity limits from rules §4.1 row 3.
const (
	MinQuantity = 1
	MaxQuantity = 99
)

// IsValidQuantity reports whether quantity is from 1 to 99 (rules §4.1 row 3).
// (It is an int, so it is always a whole number – Go's type system checked that already.)
func IsValidQuantity(quantity int) bool {
	return quantity >= MinQuantity && quantity <= MaxQuantity
}

// ValidateItems returns a *CheckoutError for rows 1-4 of rules §4.1, in that order,
// or nil when the items are fine.
func ValidateItems(basket Basket, catalog Catalog) error {
	// Row 1: no items at all.
	if len(basket.Items) == 0 {
		return newError(CodeEmptyBasket, "The basket has no items")
	}

	// Row 2: every SKU must exist. All SKUs are checked before any quantity, so an unknown
	// product on line 3 still wins over a bad quantity on line 1.
	for _, item := range basket.Items {
		if _, err := RequireProduct(catalog, item.SKU); err != nil {
			return err
		}
	}

	// Row 3: every quantity must be from 1 to 99.
	for _, item := range basket.Items {
		if !IsValidQuantity(item.Quantity) {
			return newError(CodeInvalidQuantity, "Invalid quantity for %s: %d", item.SKU, item.Quantity)
		}
	}

	// Row 4: the lines of one SKU, added together, may not go over 99 either.
	for _, line := range MergeItems(basket.Items) {
		if line.Quantity > MaxQuantity {
			return newError(CodeInvalidQuantity, "Too many %s in total: %d", line.SKU, line.Quantity)
		}
	}

	return nil
}
