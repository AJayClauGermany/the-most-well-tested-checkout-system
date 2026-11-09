//go:build e2e

// The line above is a BUILD TAG (build constraint). Go only compiles this file
// when you ask for the "e2e" tag:
//
//	go test -tags e2e ./e2e/...
//
// A plain `go test ./...` ignores this file completely. That is how Go keeps
// slow end-to-end (and integration) tests out of the everyday test run.

// Package e2e holds END-TO-END tests: they build the real command line program
// (rules §8), start it as a separate process, feed it a basket on standard input
// and check what it prints and its exit code – testing the whole program from the
// outside, exactly as a user (or the shared runner in ../../e2e/) would.
package e2e

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/AJayClauGermany/the-most-well-tested-checkout-system/go/internal/spec"
)

// binary is the path of the CLI, built once for all tests by TestMain.
var binary string

// TestMain runs before the tests of this package: it builds cmd/checkout into a temporary
// folder (go build, like a user would), runs the tests, and cleans up.
func TestMain(m *testing.M) {
	dir, err := os.MkdirTemp("", "checkout-e2e-")
	if err != nil {
		panic(err)
	}
	binary = filepath.Join(dir, "checkout")
	build := exec.Command("go", "build", "-o", binary, "../cmd/checkout")
	build.Stderr = os.Stderr
	if err := build.Run(); err != nil {
		panic("building the CLI: " + err.Error())
	}

	code := m.Run()
	os.RemoveAll(dir)
	os.Exit(code)
}

// output is what one run of the program produced.
type output struct {
	stdout, stderr string
	exitCode       int
}

// runCLI starts the program in folder dir with input on standard input and waits for it.
func runCLI(t *testing.T, dir, input string) output {
	t.Helper()
	cmd := exec.Command(binary)
	cmd.Dir = dir
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr

	err := cmd.Run()

	var exitErr *exec.ExitError
	if err != nil && !errors.As(err, &exitErr) {
		t.Fatalf("could not run the program: %v", err)
	}
	return output{stdout: stdout.String(), stderr: stderr.String(), exitCode: cmd.ProcessState.ExitCode()}
}

func TestEverySharedSpecCaseThroughTheRealProgram(t *testing.T) {
	t.Parallel()
	cases, err := spec.LoadCases()
	if err != nil {
		t.Fatal(err)
	}

	for _, c := range cases {
		t.Run(c.ID, func(t *testing.T) {
			t.Parallel()
			input, _ := json.Marshal(c.Basket)

			result := runCLI(t, ".", string(input)+"\n")

			if c.ExpectedError != "" {
				if result.exitCode != 1 {
					t.Errorf("exit code = %d, want 1 (stderr: %s)", result.exitCode, result.stderr)
				}
				want := `{"error":"` + string(c.ExpectedError) + `"}` + "\n"
				if result.stdout != want {
					t.Errorf("stdout = %q, want %q", result.stdout, want)
				}
				return
			}
			if result.exitCode != 0 {
				t.Fatalf("exit code = %d, want 0 (stdout: %s, stderr: %s)", result.exitCode, result.stdout, result.stderr)
			}
			var receipt any
			if err := json.Unmarshal([]byte(result.stdout), &receipt); err != nil {
				t.Fatalf("stdout is not JSON: %v\n%s", err, result.stdout)
			}
			var totals spec.Totals
			_ = json.Unmarshal([]byte(result.stdout), &totals)
			if diff := cmp.Diff(*c.Expected, totals); diff != "" {
				t.Errorf("%s\ntotals mismatch (-want +got):\n%s", c.Description, diff)
			}
			problems, err := spec.SchemaErrors("receipt.schema.json", receipt)
			if err != nil {
				t.Fatal(err)
			}
			if len(problems) != 0 {
				t.Errorf("the receipt does not follow receipt.schema.json: %v", problems)
			}
		})
	}
}

func TestPrintsExactlyOneLineOfJSON(t *testing.T) {
	t.Parallel()
	result := runCLI(t, ".", `{"items":[{"sku":"APPLE","quantity":1}]}`)

	if !strings.HasSuffix(result.stdout, "\n") || strings.Count(result.stdout, "\n") != 1 {
		t.Errorf("stdout = %q, want one line ending in a newline", result.stdout)
	}
}

func TestExitsWith2AndExplainsOnStderrWhenTheInputIsNotJSON(t *testing.T) {
	t.Parallel()
	result := runCLI(t, ".", "this is not json")

	if result.exitCode != 2 {
		t.Errorf("exit code = %d, want 2", result.exitCode)
	}
	if result.stdout != "" {
		t.Errorf("stdout = %q, want nothing", result.stdout)
	}
	if !strings.HasPrefix(result.stderr, "checkout: ") {
		t.Errorf("stderr = %q, want a message starting with 'checkout: '", result.stderr)
	}
}

func TestExitsWith2WhenTheJSONIsNotABasket(t *testing.T) {
	t.Parallel()
	result := runCLI(t, ".", `{"hello": "world"}`)

	if result.exitCode != 2 || !strings.Contains(result.stderr, "not a basket") {
		t.Errorf("exit code = %d, stderr = %q; want 2 and 'not a basket'", result.exitCode, result.stderr)
	}
}

func TestExitsWith2OnEmptyInput(t *testing.T) {
	t.Parallel()
	if result := runCLI(t, ".", ""); result.exitCode != 2 {
		t.Errorf("exit code = %d, want 2", result.exitCode)
	}
}

func TestFindsTheCatalogFromAnyFolderInsideTheRepository(t *testing.T) {
	t.Parallel()
	specDir, err := spec.Dir()
	if err != nil {
		t.Fatal(err)
	}
	// From the repository root, and from deep inside the spec folder.
	for _, dir := range []string{filepath.Dir(specDir), filepath.Join(specDir, "test-cases")} {
		result := runCLI(t, dir, `{"items":[{"sku":"APPLE","quantity":1}]}`)
		if result.exitCode != 0 {
			t.Errorf("from %s: exit code = %d, stderr = %q", dir, result.exitCode, result.stderr)
		}
	}
}

func TestExitsWith2WhenTheCatalogCannotBeFound(t *testing.T) {
	t.Parallel()
	// The binary lives in a temporary folder and this test runs it from another one: no
	// spec/ above either place.
	result := runCLI(t, t.TempDir(), `{"items":[{"sku":"APPLE","quantity":1}]}`)

	if result.exitCode != 2 || !strings.Contains(result.stderr, "spec/catalog.json") {
		t.Errorf("exit code = %d, stderr = %q; want 2 and a message about spec/catalog.json", result.exitCode, result.stderr)
	}
}
