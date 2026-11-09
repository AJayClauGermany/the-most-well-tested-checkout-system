package checkout

// RequireProduct is like Catalog.Find, but a missing product is an error:
// a *CheckoutError with the code UNKNOWN_PRODUCT (rules §4.1 row 2).
//
// TODO: not implemented yet.
func RequireProduct(catalog Catalog, sku string) (Product, error) {
	return Product{}, ErrNotImplemented
}
