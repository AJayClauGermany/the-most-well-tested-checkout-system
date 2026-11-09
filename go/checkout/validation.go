package checkout

// Basket validation (rules §4.1, rows 1 to 4).

// The quantity limits from rules §4.1 row 3.
const (
	MinQuantity = 1
	MaxQuantity = 99
)

// IsValidQuantity reports whether quantity is from 1 to 99 (rules §4.1 row 3).
//
// TODO: not implemented yet.
func IsValidQuantity(quantity int) bool {
	return false
}

// ValidateItems returns a *CheckoutError for rows 1-4 of rules §4.1, in that order,
// or nil when the items are fine.
//
// TODO: not implemented yet.
func ValidateItems(basket Basket, catalog Catalog) error {
	return ErrNotImplemented
}
