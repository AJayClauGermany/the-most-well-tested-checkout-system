package spec

// Tests for the schema helper (test-support code). The spec's own example outputs must
// pass, and something clearly wrong must fail – otherwise every contract test that relies
// on SchemaErrors would prove nothing.

import (
	"path/filepath"
	"testing"

	"github.com/AJayClauGermany/the-most-well-tested-checkout-system/go/checkout"
)

func TestTheSpecExamplesFollowTheirSchemas(t *testing.T) {
	t.Parallel()
	dir, err := Dir()
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		example string
		schema  string
	}{
		{example: "receipt.json", schema: "receipt.schema.json"},
		{example: "order.json", schema: "order.schema.json"}, // its "$ref" to the receipt schema must resolve
		{example: "cli-error.json", schema: "cli-error.schema.json"},
	}

	for _, tt := range tests {
		t.Run(tt.example, func(t *testing.T) {
			t.Parallel()
			var value any
			if err := ReadJSON(filepath.Join(dir, "examples", tt.example), &value); err != nil {
				t.Fatal(err)
			}

			problems, err := SchemaErrors(tt.schema, value)

			if err != nil {
				t.Fatalf("SchemaErrors() returned an error: %v", err)
			}
			if len(problems) != 0 {
				t.Errorf("%s does not follow %s: %v", tt.example, tt.schema, problems)
			}
		})
	}
}

func TestSchemaErrorsReportsEachProblemWithItsLocation(t *testing.T) {
	t.Parallel()
	order := map[string]any{"id": "ORDER-01", "created_at": "today", "receipt": map[string]any{}}

	problems, err := SchemaErrors("order.schema.json", order)

	if err != nil {
		t.Fatalf("SchemaErrors() returned an error: %v", err)
	}
	if len(problems) < 3 {
		t.Errorf("got %d problems, want at least 3 (id, created_at, receipt): %v", len(problems), problems)
	}
	for i := 1; i < len(problems); i++ {
		if problems[i-1] > problems[i] {
			t.Errorf("problems are not sorted: %v", problems)
		}
	}
}

func TestSchemaErrorsFailsForAnUnknownSchema(t *testing.T) {
	t.Parallel()
	if _, err := SchemaErrors("no-such.schema.json", map[string]any{}); err == nil {
		t.Error("SchemaErrors() with an unknown schema returned no error")
	}
}

func TestSchemaErrorsFailsForAValueThatIsNotJSON(t *testing.T) {
	t.Parallel()
	if _, err := SchemaErrors("receipt.schema.json", make(chan int)); err == nil {
		t.Error("SchemaErrors() with a channel returned no error")
	}
}

func TestTotalsOfPicksTheSixTotals(t *testing.T) {
	t.Parallel()
	dir, err := Dir()
	if err != nil {
		t.Fatal(err)
	}
	var receipt checkout.Receipt
	if err := ReadJSON(filepath.Join(dir, "examples", "receipt.json"), &receipt); err != nil {
		t.Fatal(err)
	}

	got := TotalsOf(receipt)

	want := Totals{Subtotal: 2797, PromotionDiscount: 899, CouponDiscount: 190, Shipping: 499, Tax: 137, Total: 2344}
	if got != want {
		t.Errorf("TotalsOf() = %+v, want %+v", got, want)
	}
}
