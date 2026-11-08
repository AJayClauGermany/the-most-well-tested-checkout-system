# Checkout – Go

The Go version of the checkout system described in [`../spec/checkout-rules.md`](../spec/checkout-rules.md).

> **Status: tests first, code later.** The data types and the test setup are ready,
> but `checkout.Checkout` is not implemented yet – it returns `ErrNotImplemented`.
> The 60 shared spec cases are written and **skipped**, waiting for someone to
> implement it (that is test-driven development!).

## Tools we use, and why

| Tool | What it does | Why |
|------|--------------|-----|
| [Go](https://go.dev/) 1.24+ | The language and its `go` command | Everything below is built in – nothing to install |
| `go test` + the [`testing`](https://pkg.go.dev/testing) package | Finds and runs tests | Go's standard test tool; every Go project uses it |
| `gofmt` | Formats code | There is one official Go style, so nobody argues about it |
| `go vet` | Static check: finds suspicious code (e.g. wrong `Printf` arguments) without running it | Ships with Go and catches real bugs |

This baseline uses **only the standard library**, so there is no `go.sum` file and
nothing to download. Extra libraries (assertions, BDD, …) are welcome as
[contributions](#testing-ideas-wed-love-contributions-for).

## Prerequisites and install

1. Install Go 1.24 or newer from <https://go.dev/dl/> and check it works:

   ```bash
   go version
   ```

2. There is nothing else to install. From the repository root:

   ```bash
   cd go
   go test ./...
   ```

## Running tests

All commands are run from this `go/` folder. `./...` means "this folder and every folder below it".

| I want to… | Command |
|------------|---------|
| run all tests | `go test ./...` |
| see every test name (passed / skipped) | `go test -v ./...` |
| run the tests of one package (folder) | `go test ./checkout` |
| run one test function | `go test ./checkout -run TestCatalogFindReturnsTheProductWithThatSKU -v` |
| run one subtest (e.g. one spec case) | `go test ./checkout -run 'TestSpecCases/promo-coffee-three-bags-one-free' -v` |
| run only the quick unit tests | `go test -short ./...` |
| also run the end-to-end tests | `go test -tags e2e ./...` |
| check formatting and static check | `gofmt -l .` (should print nothing) and `go vet ./...` |
| fix formatting | `gofmt -w .` |

**One file?** Go runs tests by *package* (folder) and *name*, not by file.
To run the tests from one file, use `-run` with a pattern that matches their
names, e.g. `go test ./checkout -run 'TestCatalog' -v` for `catalog_test.go`.
`-run` takes a regular expression; a `/` separates a test from its subtests.

Go caches test results. If you see `(cached)` and want a real re-run, add `-count=1`.

## Folder layout

```
go/
├── go.mod                      ← module name and Go version (like package.json / pom.xml)
├── checkout/                   ← the checkout package (the "production" code)
│   ├── types.go                ← Product, Promotion, Catalog, Basket, Receipt, …
│   ├── errors.go               ← ErrorCode, CheckoutError, ErrNotImplemented
│   ├── checkout.go             ← Checkout(basket, catalog) – not implemented yet
│   ├── catalog_test.go         ← example unit tests (start reading here!)
│   ├── errors_test.go          ← unit tests for the error codes
│   └── spec_cases_test.go      ← one subtest per shared spec case (skipped for now)
├── internal/spec/              ← test support: loads ../spec/*.json into Go types
│   ├── spec.go
│   └── spec_test.go            ← tests for the loader itself
└── e2e/
    └── cli_test.go             ← placeholder for end-to-end tests (build tag "e2e")
```

The shared test data is **not copied** here. The loader in `internal/spec` walks up
from the folder `go test` runs in until it finds the repository's `spec/` folder.

## How Go separates different kinds of tests

Other languages often have a `tests/` folder. Go does it differently:

* **Unit tests live right next to the code.** A file named `something_test.go`
  in the same folder is only compiled by `go test`, never by `go build`.
  - `package checkout` in a test file = *white-box* test (can see unexported names).
  - `package checkout_test` = *black-box* test (only the public API, like a real user).
* **`-run` picks tests by name.** Good names (`TestCatalogFind…`, `TestSpecCases/<case id>`)
  make it easy to run just the ones you want.
* **`-short` skips slow tests.** A slow test starts with
  `if testing.Short() { t.Skip("slow test") }`, and `go test -short ./...` skips it.
  (No test does this yet, so today `-short` runs everything.)
* **Build tags hide whole files.** A file that starts with `//go:build integration`
  (or `e2e`) is only compiled when you pass `-tags integration`. This is the usual
  way to keep **integration tests** (real files, databases, …) and **end-to-end tests**
  out of the everyday run:
  - integration tests: next to the code they test, e.g. a future
    `checkout/order_file_repository_integration_test.go` with `//go:build integration`,
    run with `go test -tags integration ./...`;
  - end-to-end tests: in [`e2e/`](e2e/) with `//go:build e2e`, run with
    `go test -tags e2e ./e2e/...`. See the placeholder in `e2e/cli_test.go`.
* **BDD tests** (the Gherkin files in [`../spec/features`](../spec/features)) would use
  [godog](https://github.com/cucumber/godog), e.g. in a `bdd/` folder with a
  `features_test.go` that reads `../spec/features` directly.

## How to add a new test

Say you want to check that `Catalog.Find` finds the first product in the list.

1. Open (or create) a file ending in `_test.go` next to the code – here
   `checkout/catalog_test.go`.
2. Write a function whose name starts with `Test` and takes `t *testing.T`:

   ```go
   func TestCatalogFindFindsTheFirstProduct(t *testing.T) {
       // Arrange
       catalog := Catalog{Products: []Product{{SKU: "APPLE"}, {SKU: "PEN"}}}

       // Act
       _, found := catalog.Find("APPLE")

       // Assert
       if !found {
           t.Errorf("Find(%q) did not find the product", "APPLE")
       }
   }
   ```

3. Run it: `go test ./checkout -run TestCatalogFindFindsTheFirstProduct -v`.
4. Make it fail on purpose once (change `"APPLE"` to `"MANGO"`) to see that it really tests something.
5. Before you commit: `gofmt -w .`, `go vet ./...` and `go test ./...`.

Many similar cases? Use a **table-driven test** – see `TestCatalogFindDoesNotFindUnknownSKUs`.

**Starting the implementation?** Remove the `t.Skip(...)` line in
`checkout/spec_cases_test.go`, watch all 60 cases fail (red), then make them pass
one by one (green). Small helper functions (rounding, merging lines, …) deserve their
own unit tests too.

## Testing ideas we'd love contributions for

* **BDD:** run `../spec/features/*.feature` with [godog](https://github.com/cucumber/godog).
* **Nicer assertions and diffs:** [testify](https://github.com/stretchr/testify) (`assert`/`require`)
  or [go-cmp](https://github.com/google/go-cmp) (`cmp.Diff` for whole receipts).
* **Property-based testing:** [rapid](https://github.com/flyingmutant/rapid) or the built-in
  [`testing/quick`](https://pkg.go.dev/testing/quick), e.g. "total is never negative",
  "splitting a line into two lines does not change the price".
* **Fuzzing:** Go's built-in fuzzing (`func FuzzXxx(f *testing.F)`, run with `go test -fuzz`),
  e.g. feeding random JSON baskets – checkout must never panic.
* **Example tests:** `func ExampleCheckout()` with an `// Output:` comment – documentation
  that is also a test.
* **Mocks and fakes:** a hand-written in-memory `OrderRepository` fake (rules §6), or generated
  mocks with [gomock](https://github.com/uber-go/mock) / [mockery](https://github.com/vektra/mockery).
* **Snapshot / golden-file tests:** compare receipt JSON with files in `testdata/`
  (hand-rolled with an `-update` flag, or [cupaloy](https://github.com/bradleyjkemp/cupaloy)).
* **Integration tests:** a JSON-file or SQLite ([modernc.org/sqlite](https://pkg.go.dev/modernc.org/sqlite))
  order repository tested with `t.TempDir()` behind `//go:build integration`;
  [testcontainers-go](https://golang.testcontainers.org/) for a real database.
* **End-to-end tests:** build the CLI from rules §8 in `cmd/checkout`, then drive it with
  `os/exec` in `e2e/`, or with [testscript](https://pkg.go.dev/github.com/rogpeppe/go-internal/testscript).
* **Coverage:** `go test -cover ./...`, or `go test -coverprofile=coverage.out ./...` then
  `go tool cover -html=coverage.out`.
* **Mutation testing:** [gremlins](https://github.com/go-gremlins/gremlins) or
  [go-mutesting](https://github.com/avito-tech/go-mutesting) – do the tests notice when the code is changed?
* **Benchmarks:** `func BenchmarkCheckout(b *testing.B)`, run with `go test -bench .`.
* **Race detector:** `go test -race ./...` once something runs concurrently (e.g. a shared repository).
* **Stricter linting:** [staticcheck](https://staticcheck.dev/) or [golangci-lint](https://golangci-lint.run/).
