package spec

// Tests for the spec loader. If these fail, the shared test data could not be
// read correctly, and every test that uses it would be meaningless – so it is
// worth checking the loader itself.

import (
	"testing"

	"github.com/AJayClauGermany/the-most-well-tested-checkout-system/go/checkout"
)

func TestLoadCatalogHasSixProducts(t *testing.T) {
	catalog, err := LoadCatalog()
	if err != nil {
		t.Fatalf("LoadCatalog() returned an error: %v", err)
	}

	if got := len(catalog.Products); got != 6 {
		t.Errorf("catalog has %d products, want 6", got)
	}
}

func TestLoadCatalogReadsCoffeePromotion(t *testing.T) {
	catalog, err := LoadCatalog()
	if err != nil {
		t.Fatalf("LoadCatalog() returned an error: %v", err)
	}

	coffee, found := catalog.Find("COFFEE")
	if !found {
		t.Fatal("COFFEE is not in the catalog")
	}
	if coffee.Promotion == nil {
		t.Fatal("COFFEE has no promotion, want buy 2 get 1 free")
	}

	want := checkout.Promotion{Type: checkout.BuyXGetYFree, Buy: 2, Free: 1}
	if *coffee.Promotion != want {
		t.Errorf("COFFEE promotion = %+v, want %+v", *coffee.Promotion, want)
	}
}

func TestLoadCasesFindsAllSixtyCases(t *testing.T) {
	cases, err := LoadCases()
	if err != nil {
		t.Fatalf("LoadCases() returned an error: %v", err)
	}

	if got := len(cases); got != 60 {
		t.Errorf("loaded %d cases, want 60", got)
	}
}

func TestLoadCasesIDsAreUnique(t *testing.T) {
	cases, err := LoadCases()
	if err != nil {
		t.Fatalf("LoadCases() returned an error: %v", err)
	}

	seenIn := map[string]string{} // case id -> file it was first seen in
	for _, c := range cases {
		if firstFile, seen := seenIn[c.ID]; seen {
			t.Errorf("case id %q is in both %s and %s", c.ID, firstFile, c.File)
		}
		seenIn[c.ID] = c.File
	}
}

func TestLoadCasesHaveExactlyOneExpectation(t *testing.T) {
	cases, err := LoadCases()
	if err != nil {
		t.Fatalf("LoadCases() returned an error: %v", err)
	}

	for _, c := range cases {
		hasTotals := c.Expected != nil
		hasError := c.ExpectedError != ""
		if hasTotals == hasError {
			t.Errorf("%s (%s): want exactly one of expected / expected_error, got expected=%v expected_error=%q",
				c.ID, c.File, hasTotals, c.ExpectedError)
		}
	}
}
