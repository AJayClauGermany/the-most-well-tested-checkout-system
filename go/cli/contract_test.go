package cli_test

// CONTRACT TESTS for what the program prints (rules §8): a receipt that follows
// spec/schema/receipt.schema.json, or an error that follows spec/schema/cli-error.schema.json.
// Checked here on the JSON TEXT the program prints – exactly what another program reads.

import (
	"encoding/json"
	"testing"

	"github.com/AJayClauGermany/the-most-well-tested-checkout-system/go/cli"
	"github.com/AJayClauGermany/the-most-well-tested-checkout-system/go/internal/spec"
)

func TestTheOutputForEverySpecCaseFollowsItsSchema(t *testing.T) {
	t.Parallel()
	catalog, err := spec.LoadCatalog()
	if err != nil {
		t.Fatal(err)
	}
	cases, err := spec.LoadCases()
	if err != nil {
		t.Fatal(err)
	}

	for _, c := range cases {
		t.Run(c.ID, func(t *testing.T) {
			t.Parallel()
			input, _ := json.Marshal(c.Basket)
			result, err := cli.Run(string(input), catalog)
			if err != nil {
				t.Fatalf("Run() returned an error: %v", err)
			}

			schema, wantExitCode := "receipt.schema.json", 0
			if c.ExpectedError != "" {
				schema, wantExitCode = "cli-error.schema.json", 1
			}
			if result.ExitCode != wantExitCode {
				t.Errorf("exit code = %d, want %d", result.ExitCode, wantExitCode)
			}
			var printed any
			if err := json.Unmarshal([]byte(result.Output), &printed); err != nil {
				t.Fatalf("the output is not JSON: %v", err)
			}
			problems, err := spec.SchemaErrors(schema, printed)
			if err != nil {
				t.Fatal(err)
			}
			if len(problems) != 0 {
				t.Errorf("the output does not follow %s: %v\n%s", schema, problems, result.Output)
			}
		})
	}
}

func TestTheErrorSchemaRejectsAnInventedErrorCode(t *testing.T) {
	t.Parallel()
	// Testing the test: the schema really does check the code.
	problems, err := spec.SchemaErrors("cli-error.schema.json", map[string]any{"error": "OUT_OF_STOCK"})

	if err != nil {
		t.Fatal(err)
	}
	if len(problems) == 0 {
		t.Error("an invented error code was accepted")
	}
}
