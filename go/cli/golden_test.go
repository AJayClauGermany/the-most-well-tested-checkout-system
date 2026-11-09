package cli_test

// SNAPSHOT (golden file) TEST of the exact text the program prints for the example basket
// of rules §5. Update it after an intended change with:  go test ./cli -update

import (
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/AJayClauGermany/the-most-well-tested-checkout-system/go/cli"
	"github.com/AJayClauGermany/the-most-well-tested-checkout-system/go/internal/fixtures"
	"github.com/AJayClauGermany/the-most-well-tested-checkout-system/go/internal/spec"
)

var update = flag.Bool("update", false, "rewrite the golden files in testdata/ with the current output")

func TestTheExampleOutputMatchesItsGoldenFile(t *testing.T) {
	t.Parallel()
	catalog, err := spec.LoadCatalog()
	if err != nil {
		t.Fatal(err)
	}
	input, _ := json.Marshal(fixtures.ExampleBasket)

	result, err := cli.Run(string(input), catalog)
	if err != nil {
		t.Fatal(err)
	}

	got := result.Output + "\n" // the program prints one line, ending in a newline
	path := filepath.Join("testdata", "example-output.golden.txt")
	if *update {
		if err := os.WriteFile(path, []byte(got), 0o644); err != nil {
			t.Fatal(err)
		}
		return
	}
	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading the golden file: %v (run `go test ./cli -update` to create it)", err)
	}
	if diff := cmp.Diff(string(want), got); diff != "" {
		t.Errorf("output does not match the golden file (-golden +got):\n%s", diff)
	}
}
