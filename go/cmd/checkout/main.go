// The command line program (rules §8): basket JSON in on stdin, receipt JSON out on stdout.
//
// TODO: not implemented yet – the end-to-end tests in e2e/ are waiting for it.
package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Fprintln(os.Stderr, "checkout: not implemented yet")
	os.Exit(2)
}
