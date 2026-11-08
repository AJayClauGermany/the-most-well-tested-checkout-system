// Package spec loads the shared, language-independent test data from the
// repository's top-level spec/ folder, so Go tests can use it.
//
// It is test-support code. It lives under internal/ because Go only lets
// packages inside this module import an "internal" package – it is not part
// of the checkout library that other people would use.
package spec

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/AJayClauGermany/the-most-well-tested-checkout-system/go/checkout"
)

// Totals are the six amounts a spec case expects (all in cents).
type Totals struct {
	Subtotal          int `json:"subtotal"`
	PromotionDiscount int `json:"promotion_discount"`
	CouponDiscount    int `json:"coupon_discount"`
	Shipping          int `json:"shipping"`
	Tax               int `json:"tax"`
	Total             int `json:"total"`
}

// Case is one shared test case from spec/test-cases/*.json.
// Exactly one of Expected and ExpectedError is set.
type Case struct {
	File          string             `json:"-"` // which file it came from, e.g. "01-subtotal.json"
	ID            string             `json:"id"`
	Description   string             `json:"description"`
	Basket        checkout.Basket    `json:"basket"`
	Expected      *Totals            `json:"expected"`       // nil when an error is expected
	ExpectedError checkout.ErrorCode `json:"expected_error"` // "" when totals are expected
}

// Dir finds the repository's spec/ folder.
//
// `go test` runs each package's tests inside that package's own folder
// (for example go/internal/spec), so we cannot rely on a fixed relative path.
// Instead we start at the current folder and walk up, one parent at a time,
// until we find a folder that contains spec/catalog.json.
func Dir() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		candidate := filepath.Join(dir, "spec")
		if _, err := os.Stat(filepath.Join(candidate, "catalog.json")); err == nil {
			return candidate, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir { // reached the root of the file system
			return "", errors.New("could not find the spec/ folder in any parent directory")
		}
		dir = parent
	}
}

// LoadCatalog reads spec/catalog.json.
func LoadCatalog() (checkout.Catalog, error) {
	var catalog checkout.Catalog
	dir, err := Dir()
	if err != nil {
		return catalog, err
	}
	err = readJSON(filepath.Join(dir, "catalog.json"), &catalog)
	return catalog, err
}

// LoadCases reads every spec/test-cases/*.json file, in file name order.
func LoadCases() ([]Case, error) {
	dir, err := Dir()
	if err != nil {
		return nil, err
	}
	paths, err := filepath.Glob(filepath.Join(dir, "test-cases", "*.json"))
	if err != nil {
		return nil, err
	}
	sort.Strings(paths)

	var all []Case
	for _, path := range paths {
		var file struct {
			Cases []Case `json:"cases"`
		}
		if err := readJSON(path, &file); err != nil {
			return nil, err
		}
		for _, c := range file.Cases {
			c.File = filepath.Base(path)
			all = append(all, c)
		}
	}
	return all, nil
}

// readJSON decodes the JSON file at path into target.
func readJSON(path string, target any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("%s: %w", path, err)
	}
	return nil
}
