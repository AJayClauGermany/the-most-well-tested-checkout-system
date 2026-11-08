package checkout

// Checkout prices a basket using the products in catalog (see spec/checkout-rules.md §4).
//
// It returns either a Receipt and a nil error, or an empty Receipt and an error.
// For an invalid basket the error is a *CheckoutError carrying one of the six codes.
//
// TODO: not implemented yet. This repository is test-first: the shared spec
// cases in spec_cases_test.go are waiting (skipped) for this function. When you
// start implementing it, remove the t.Skip in that test and make them pass.
func Checkout(basket Basket, catalog Catalog) (Receipt, error) {
	return Receipt{}, ErrNotImplemented
}
