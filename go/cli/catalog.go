package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/AJayClauGermany/the-most-well-tested-checkout-system/go/checkout"
	"github.com/AJayClauGermany/the-most-well-tested-checkout-system/go/specdir"
)

// StandardCatalogPath returns where the standard catalog lives: spec/catalog.json at the
// root of the repository (found by walking up from the working directory or the program).
func StandardCatalogPath() (string, error) {
	dir, err := specdir.Locate()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "catalog.json"), nil
}

// LoadStandardCatalog reads spec/catalog.json.
func LoadStandardCatalog() (checkout.Catalog, error) {
	path, err := StandardCatalogPath()
	if err != nil {
		return checkout.Catalog{}, err
	}
	return LoadCatalog(path)
}

// LoadCatalog reads a catalog from a JSON file shaped like spec/catalog.json.
func LoadCatalog(path string) (checkout.Catalog, error) {
	var catalog checkout.Catalog
	data, err := os.ReadFile(path)
	if err != nil {
		return catalog, err
	}
	if err := json.Unmarshal(data, &catalog); err != nil {
		return catalog, fmt.Errorf("%s: %w", path, err)
	}
	return catalog, nil
}
