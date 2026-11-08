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

This list mirrors the Go column of the [testing matrix](../docs/testing-matrix.md) – the
project's single list of goals – and adds the practical hints. Tick the matrix when you add
one. Tip: pin tools like `godog` and `gremlins` with Go 1.24's `go tool` directive in `go.mod`
so nobody has to install them by hand.

**Tier 1 – every language does these**

* **Implement `Checkout`** test-first, with unit tests for each step.
* **Table-driven subtests** (`[]struct{...}` + `t.Run`) – *the* Go idiom, shown on its own
  next to the shared spec cases.
* **BDD:** run `../spec/features/*.feature` with [godog](https://github.com/cucumber/godog).
* **Nicer assertions and diffs:** [testify](https://github.com/stretchr/testify) (`assert`/`require`)
  or [go-cmp](https://github.com/google/go-cmp) (`cmp.Diff` for whole receipts).
* **Test doubles** for the `OrderRepository` and `Clock` interfaces (rules §6): a hand-written
  in-memory fake first, then [uber-go/mock](https://github.com/uber-go/mock) to check a failed
  checkout never calls `Save`, and a fixed clock for `created_at`.
* **Integration tests** for a JSON-file and an SQLite ([modernc.org/sqlite](https://pkg.go.dev/modernc.org/sqlite),
  pure Go) repository using `t.TempDir()`, behind `//go:build integration`. Check stored
  orders against [`spec/schema/order.schema.json`](../spec/schema/order.schema.json).
* **Property-based testing** with [rapid](https://github.com/flyingmutant/rapid): one test per
  property in [`spec/properties.md`](../spec/properties.md), named after its id.
  (Not `testing/quick`: it is frozen and cannot shrink failures.)
* **Contract tests** with [santhosh-tekuri/jsonschema](https://github.com/santhosh-tekuri/jsonschema):
  every receipt must satisfy [`spec/schema/receipt.schema.json`](../spec/schema/receipt.schema.json).
* **Golden-file / snapshot tests:** receipt JSON compared with files in `testdata/`, hand-rolled
  with an `-update` flag (teaches the most) or with [go-snaps](https://github.com/gkampitakis/go-snaps).
* **End-to-end tests:** build the CLI from rules §8 in `cmd/checkout`, drive it with `os/exec`
  in `e2e/`, then run the shared runner: `e2e/run.sh go run ./go/cmd/checkout`.
* **Coverage:** `go test -cover ./...`, or `-coverprofile=coverage.out` then `go tool cover -html=coverage.out`.
* **Mutation testing:** [gremlins](https://github.com/go-gremlins/gremlins) – the only maintained
  Go option; it moves slowly, so expect rough edges.
* **Stricter static checks:** [staticcheck](https://staticcheck.dev/) or
  [golangci-lint](https://golangci-lint.run/), plus `govulncheck`.

**Tier 2 – idiomatic extras**

* **Example tests:** `func ExampleCheckout()` with an `// Output:` comment – documentation
  that is also a test.
* **Fuzzing:** `func FuzzCheckout(f *testing.F)` run with `go test -fuzz` – property P10:
  random JSON must never panic.
* **Benchmarks:** `func BenchmarkCheckout(b *testing.B)` using `b.Loop()` (Go 1.24).
* **Race detector:** `go test -race ./...` once the repository is used from several goroutines.

**Tier 3 – plumbing**

* `go test -json ./... | gotestsum --junitfile report.xml` for a CI report;
  `t.Parallel()` and `go test -shuffle=on` to prove tests are independent.
