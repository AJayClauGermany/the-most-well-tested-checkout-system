package cli_test

// UNIT TESTS for the command line program's logic (rules §8), run in-process. The real
// program (cmd/checkout) is started as a separate process by the tests in e2e/.

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/AJayClauGermany/the-most-well-tested-checkout-system/go/checkout"
	"github.com/AJayClauGermany/the-most-well-tested-checkout-system/go/cli"
	"github.com/AJayClauGermany/the-most-well-tested-checkout-system/go/internal/fixtures"
	"github.com/AJayClauGermany/the-most-well-tested-checkout-system/go/specdir"
)

// mustRun runs the command and stops the test if the input was not even a basket.
func mustRun(t *testing.T, input string) cli.Result {
	t.Helper()
	result, err := cli.Run(input, fixtures.TinyCatalog)
	if err != nil {
		t.Fatalf("Run(%q) returned an error: %v", input, err)
	}
	return result
}

func TestStandardCatalogPathPointsAtSpecCatalogJSON(t *testing.T) {
	t.Parallel()
	path, err := cli.StandardCatalogPath()

	if err != nil {
		t.Fatalf("StandardCatalogPath() returned an error: %v", err)
	}
	if !strings.HasSuffix(path, filepath.Join("spec", "catalog.json")) {
		t.Errorf("path = %s, want it to end in spec/catalog.json", path)
	}
	// The repository root is the folder that also holds the rules.
	if _, err := os.Stat(filepath.Join(filepath.Dir(path), "checkout-rules.md")); err != nil {
		t.Errorf("no checkout-rules.md next to %s", path)
	}
}

func TestLoadStandardCatalogLoadsTheSixStandardProducts(t *testing.T) {
	t.Parallel()
	catalog, err := cli.LoadStandardCatalog()

	if err != nil {
		t.Fatalf("LoadStandardCatalog() returned an error: %v", err)
	}
	var skus []string
	for _, product := range catalog.Products {
		skus = append(skus, product.SKU)
	}
	want := []string{"APPLE", "BREAD", "COFFEE", "NOTEBOOK", "PEN", "HEADPHONES"}
	if diff := cmp.Diff(want, skus); diff != "" {
		t.Errorf("SKUs mismatch (-want +got):\n%s", diff)
	}
}

func TestLoadCatalogReportsAMissingFile(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "catalog.json")

	if _, err := cli.LoadCatalog(path); err == nil {
		t.Error("LoadCatalog() of a missing file returned no error")
	}
}

func TestLoadCatalogReportsAFileThatIsNotACatalog(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "catalog.json")
	if err := os.WriteFile(path, []byte(`{"products": "none"}`), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := cli.LoadCatalog(path)

	if err == nil || !strings.Contains(err.Error(), path) {
		t.Errorf("LoadCatalog() error = %v, want it to name the file", err)
	}
}

// Not parallel: t.Setenv and t.Chdir change the whole test process.
func TestLoadStandardCatalogFailsClearlyOutsideTheRepository(t *testing.T) {
	t.Setenv(specdir.EnvVar, "") // no shortcut to the spec folder either
	t.Chdir(t.TempDir())         // no spec/ above here, nor above the test binary

	if _, err := cli.LoadStandardCatalog(); err == nil {
		t.Error("LoadStandardCatalog() found a catalog outside the repository")
	}
	if _, err := cli.StandardCatalogPath(); err == nil {
		t.Error("StandardCatalogPath() found a catalog outside the repository")
	}
}

func TestRunPrintsTheReceiptAndExitsWith0ForAValidBasket(t *testing.T) {
	t.Parallel()
	result := mustRun(t, `{"items":[{"sku":"APPLE","quantity":1}]}`)

	if result.ExitCode != 0 {
		t.Errorf("exit code = %d, want 0", result.ExitCode)
	}
	var receipt checkout.Receipt
	if err := json.Unmarshal([]byte(result.Output), &receipt); err != nil {
		t.Fatalf("the output is not JSON: %v\n%s", err, result.Output)
	}
	if diff := cmp.Diff(fixtures.OneAppleReceipt(), receipt); diff != "" {
		t.Errorf("receipt mismatch (-want +got):\n%s", diff)
	}
}

func TestRunWritesSnakeCaseFieldNamesExactlyAsTheSchemaSays(t *testing.T) {
	t.Parallel()
	result := mustRun(t, `{"items":[{"sku":"APPLE","quantity":1}]}`)

	want := `{"lines":[{"sku":"APPLE","name":"Apple","quantity":1,"unit_price":50,"line_total":50,"line_discount":0}],` +
		`"subtotal":50,"promotion_discount":0,"coupon_discount":0,"shipping":499,"tax":4,"total":553}`
	if result.Output != want {
		t.Errorf("output =\n%s\nwant\n%s", result.Output, want)
	}
}

func TestRunPrintsTheErrorCodeAndExitsWith1ForAnInvalidBasket(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input string
		want  string
	}{
		{input: `{"items":[]}`, want: `{"error":"EMPTY_BASKET"}`},
		{input: `{"items":[{"sku":"APPLE","quantity":1}],"coupon":"NOPE"}`, want: `{"error":"UNKNOWN_COUPON"}`},
		{input: `{"items":[{"sku":"APPLE","quantity":1}],"shipping":"overnight"}`, want: `{"error":"UNKNOWN_SHIPPING_METHOD"}`},
		{input: `{"items":[{"sku":"APPLE","quantity":0}]}`, want: `{"error":"INVALID_QUANTITY"}`},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			t.Parallel()
			result := mustRun(t, tt.input)

			want := cli.Result{Output: tt.want, ExitCode: 1}
			if result != want {
				t.Errorf("Run(%s) = %+v, want %+v", tt.input, result, want)
			}
		})
	}
}

func TestRunPrintsCompactJSONOnASingleLine(t *testing.T) {
	t.Parallel()
	result := mustRun(t, `{"items":[{"sku":"COFFEE","quantity":3},{"sku":"APPLE","quantity":2}]}`)

	if strings.Contains(result.Output, "\n") {
		t.Errorf("output has a line break:\n%s", result.Output)
	}
}

func TestRunAcceptsSurroundingWhitespaceAsAShellPipelineAddsANewline(t *testing.T) {
	t.Parallel()
	result := mustRun(t, "  {\"items\":[{\"sku\":\"APPLE\",\"quantity\":1}]}\n")

	if result.ExitCode != 0 {
		t.Errorf("exit code = %d, want 0", result.ExitCode)
	}
}

func TestRunTreatsNullAndMissingCouponAndShippingAsNotGiven(t *testing.T) {
	t.Parallel()
	for _, input := range []string{
		`{"items":[{"sku":"APPLE","quantity":1}],"coupon":null}`,
		`{"items":[{"sku":"APPLE","quantity":1}],"coupon":"   "}`,
		`{"items":[{"sku":"APPLE","quantity":1}]}`,
	} {
		result := mustRun(t, input)
		if result.ExitCode != 0 {
			t.Errorf("Run(%s) exit code = %d, want 0 (%s)", input, result.ExitCode, result.Output)
		}
	}
}

func TestRunFailsWithAMessageWhenTheInputIsNotABasketAtAll(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input string
		want  string
	}{
		{input: "not json", want: "the input is not JSON"},
		{input: "", want: "the input is not JSON"},
		{input: `{"items":`, want: "the input is not JSON"},
		{input: "null", want: "the input is not a basket: expected a JSON object"},
		{input: "42", want: "the input is not a basket: expected a JSON object"},
		{input: `"text"`, want: "the input is not a basket: expected a JSON object"},
		{input: "[]", want: "the input is not a basket: expected a JSON object"},
		{input: "{}", want: "the input is not a basket: 'items' must be a list"},
		{input: `{"items":null}`, want: "the input is not a basket: 'items' must be a list"},
		{input: `{"items":"APPLE"}`, want: "the input is not a basket: json: cannot unmarshal"},
		{input: `{"items":[1,2]}`, want: "the input is not a basket: json: cannot unmarshal"},
		{input: `{"items":[null]}`, want: "the input is not a basket: every item needs a string 'sku'"},
		{input: `{"items":[{"quantity":1}]}`, want: "the input is not a basket: every item needs a string 'sku'"},
		{input: `{"items":[{"sku":5,"quantity":1}]}`, want: "the input is not a basket: json: cannot unmarshal"},
		{input: `{"items":[{"sku":"APPLE"}]}`, want: "the input is not a basket: every item needs a whole-number 'quantity'"},
		{input: `{"items":[{"sku":"APPLE","quantity":1.5}]}`, want: "the input is not a basket: json: cannot unmarshal"},
		{input: `{"items":[],"coupon":5}`, want: "the input is not a basket: json: cannot unmarshal"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()
			_, err := cli.Run(tt.input, fixtures.TinyCatalog)

			if err == nil || !strings.HasPrefix(err.Error(), tt.want) {
				t.Errorf("Run(%q) error = %v, want it to start with %q", tt.input, err, tt.want)
			}
		})
	}
}
