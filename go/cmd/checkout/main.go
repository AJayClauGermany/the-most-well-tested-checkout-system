// The command line program (rules §8): basket JSON in on stdin, receipt JSON out on stdout.
//
//	echo '{"items":[{"sku":"APPLE","quantity":1}]}' | go run ./cmd/checkout
//
// All the logic is in package cli; this file only connects it to the terminal, which is
// why it is left out of the coverage report and tested end-to-end (e2e/) instead.
package main

import (
	"fmt"
	"io"
	"os"

	"github.com/AJayClauGermany/the-most-well-tested-checkout-system/go/cli"
)

func main() {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		fail(err)
	}
	catalog, err := cli.LoadStandardCatalog()
	if err != nil {
		fail(err)
	}
	result, err := cli.Run(string(input), catalog)
	if err != nil {
		fail(err) // not a basket at all (rules §8 only covers valid and invalid baskets)
	}
	fmt.Println(result.Output)
	os.Exit(result.ExitCode)
}

// fail explains the problem on stderr and exits with code 2: nothing on stdout, so a
// script reading our output never mistakes the message for a receipt.
func fail(err error) {
	fmt.Fprintf(os.Stderr, "checkout: %v\n", err)
	os.Exit(2)
}
