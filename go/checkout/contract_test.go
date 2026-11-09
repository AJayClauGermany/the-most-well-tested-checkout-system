package checkout_test

// CONTRACT TESTS: is the output even the right shape?
//
// spec/schema/receipt.schema.json is the contract every language's receipt must honour.
// Here a JSON Schema library checks real receipts against it. The shared e2e runner only
// compares totals; this is where the shape of the whole receipt (lines included) is held
// to the contract. (The CLI's error output has its own contract test in package cli.)

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/AJayClauGermany/the-most-well-tested-checkout-system/go/checkout"
	"github.com/AJayClauGermany/the-most-well-tested-checkout-system/go/internal/fixtures"
	"github.com/AJayClauGermany/the-most-well-tested-checkout-system/go/internal/spec"
)

// assertFollowsSchema fails the test unless value is valid against the named schema.
func assertFollowsSchema(t *testing.T, schema string, value any) {
	t.Helper()
	problems, err := spec.SchemaErrors(schema, value)
	if err != nil {
		t.Fatalf("loading %s: %v", schema, err)
	}
	if len(problems) != 0 {
		t.Errorf("does not follow %s:\n  %v", schema, problems)
	}
}

func TestTheReceiptForEveryValidSpecCaseFollowsTheSchema(t *testing.T) {
	t.Parallel()
	catalog := loadCatalog(t)

	for _, c := range loadCases(t) {
		if c.Expected == nil {
			continue
		}
		t.Run(c.ID, func(t *testing.T) {
			t.Parallel()
			receipt, err := checkout.Checkout(c.Basket, catalog)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			assertFollowsSchema(t, "receipt.schema.json", receipt)
		})
	}
}

func TestAReceiptFromAMadeUpCatalogFollowsTheSchemaToo(t *testing.T) {
	t.Parallel()
	receipt := mustCheckout(t, items(item("TV", 1), item("COFFEE", 6)))

	assertFollowsSchema(t, "receipt.schema.json", receipt)
}

func TestTheSchemaCheckItselfWorksABrokenReceiptIsRejected(t *testing.T) {
	t.Parallel()
	// Testing the test: if the schema accepted anything, the tests above would prove nothing.
	// A map lets us build a receipt the Go type system would never allow.
	broken := map[string]any{
		"lines":              fixtures.OneAppleReceipt().Lines,
		"subtotal":           50,
		"promotion_discount": 0,
		"coupon_discount":    0,
		"shipping":           123, // not 0, 499 or 999
		"tax":                -1,  // negative
		"total":              553,
		"extra":              true, // an unknown field
	}

	problems, err := spec.SchemaErrors("receipt.schema.json", broken)

	if err != nil {
		t.Fatalf("loading the schema: %v", err)
	}
	want := []string{
		"/ additional properties 'extra' not allowed",
		"/shipping value must be one of 0, 499, 999",
		"/tax minimum: got -1, want 0",
	}
	if diff := cmp.Diff(want, problems); diff != "" {
		t.Errorf("problems mismatch (-want +got):\n%s", diff)
	}
}

func TestEverySharedTestCasesBasketIsAValidBasket(t *testing.T) {
	t.Parallel()
	// The spec is consistent with itself: every basket in spec/test-cases/ follows
	// spec/schema/basket.schema.json. (CI checks this for the spec too; here it also proves
	// that our Basket type round-trips to the JSON the schema expects.)
	for _, c := range loadCases(t) {
		t.Run(c.ID, func(t *testing.T) {
			t.Parallel()
			assertFollowsSchema(t, "basket.schema.json", c.Basket)
		})
	}
}
