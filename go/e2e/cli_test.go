//go:build e2e

// The line above is a BUILD TAG (build constraint). Go only compiles this file
// when you ask for the "e2e" tag:
//
//	go test -tags e2e ./e2e/...
//
// A plain `go test ./...` ignores this file completely. That is how Go keeps
// slow end-to-end (and integration) tests out of the everyday test run.

// Package e2e will hold END-TO-END tests: they build the real command line
// program (rules §8), feed it a basket on standard input and check the JSON it
// prints and its exit code – testing the whole program from the outside.
package e2e

import "testing"

func TestCLIPlaceholder(t *testing.T) {
	t.Skip("the checkout command line program (rules §8) does not exist yet – " +
		"add it under go/cmd/checkout, then replace this placeholder with real e2e tests")
}
