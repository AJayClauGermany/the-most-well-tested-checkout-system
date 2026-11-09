package bdd

// BDD (behaviour-driven development) TESTS: the Given/When/Then scenarios in
// spec/features/*.feature, run by godog inside a normal `go test`.
//
// The feature files are written in Gherkin, a near-English language that people who do
// not program can read (and write). They are shared by all six languages and read straight
// from the spec folder – never copied – so every language is held to the same scenarios.
// Each scenario becomes a subtest, so one can be run on its own:
//
//	go test ./bdd -run 'TestFeatures/Half_a_cent_is_rounded_up' -v

import (
	"path/filepath"
	"testing"

	"github.com/cucumber/godog"

	"github.com/AJayClauGermany/the-most-well-tested-checkout-system/go/internal/spec"
)

func TestFeatures(t *testing.T) {
	dir, err := spec.Dir()
	if err != nil {
		t.Fatalf("finding the spec folder: %v", err)
	}

	suite := godog.TestSuite{
		Name:                "checkout",
		ScenarioInitializer: InitializeScenario, // registers the steps, see steps_test.go
		Options: &godog.Options{
			Format:    "progress", // one dot per step; use "pretty" to see every sentence
			Paths:     []string{filepath.Join(dir, "features")},
			TestingT:  t,    // run every scenario as a subtest of this test
			Strict:    true, // a step without a definition fails instead of being skipped
			Randomize: -1,   // run the scenarios in a random order: they must not depend on each other
		},
	}

	if status := suite.Run(); status != 0 {
		t.Fatalf("godog finished with status %d", status)
	}
}
