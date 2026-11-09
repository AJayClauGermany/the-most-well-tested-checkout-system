package checkout

// Checkout prices a basket using the products in catalog (see spec/checkout-rules.md §4).
//
// It returns either a Receipt and a nil error, or an empty Receipt and an error.
// For an invalid basket the error is a *CheckoutError carrying one of the six codes.
//
// TODO: not implemented yet. This repository is test-first: every test is written and
// red, waiting for the implementation.
func Checkout(basket Basket, catalog Catalog) (Receipt, error) {
	return Receipt{}, ErrNotImplemented
}
