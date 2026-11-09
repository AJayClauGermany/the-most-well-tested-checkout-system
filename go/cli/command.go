// Package cli holds the logic of the command line program (rules §8), kept free of
// os.Stdin and os.Stdout so it can be unit tested in-process.
package cli

import (
	"errors"

	"github.com/AJayClauGermany/the-most-well-tested-checkout-system/go/checkout"
)

var errNotImplemented = errors.New("cli: not implemented yet")

// Result is what the program prints (one line, without the newline) and its exit code.
type Result struct {
	Output   string
	ExitCode int
}

// Run runs one checkout from the JSON text of a basket. TODO: not implemented yet.
func Run(input string, catalog checkout.Catalog) (Result, error) {
	return Result{}, errNotImplemented
}
