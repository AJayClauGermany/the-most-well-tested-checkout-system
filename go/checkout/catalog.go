package checkout

// RequireProduct is like Catalog.Find, but a missing product is an error:
// a *CheckoutError with the code UNKNOWN_PRODUCT (rules §4.1 row 2).
// SKUs are case sensitive, so "apple" does not find "APPLE".
func RequireProduct(catalog Catalog, sku string) (Product, error) {
	product, found := catalog.Find(sku)
	if !found {
		return Product{}, newError(CodeUnknownProduct, "Unknown product: %s", sku)
	}
	return product, nil
}
