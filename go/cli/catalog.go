package cli

import "github.com/AJayClauGermany/the-most-well-tested-checkout-system/go/checkout"

// StandardCatalogPath returns where spec/catalog.json is. TODO: not implemented yet.
func StandardCatalogPath() (string, error) {
	return "", errNotImplemented
}

// LoadStandardCatalog reads spec/catalog.json. TODO: not implemented yet.
func LoadStandardCatalog() (checkout.Catalog, error) {
	return checkout.Catalog{}, errNotImplemented
}

// LoadCatalog reads a catalog from a JSON file. TODO: not implemented yet.
func LoadCatalog(path string) (checkout.Catalog, error) {
	return checkout.Catalog{}, errNotImplemented
}
