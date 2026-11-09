# Checkout – Go

The Go version of the checkout system described in [`../spec/checkout-rules.md`](../spec/checkout-rules.md).

> **Status: done.** `checkout.Checkout`, saving orders (in memory, as JSON files and in
> SQLite) and the command line program are implemented, and every row of the Go column of
> the [testing matrix](../docs/testing-matrix.md) is covered. It was built strictly
> test-first: the commit that adds the tests comes before the commit that makes them
> pass – have a look at the git history.

## Tools we use (and why)

Go ships with most of what a tester needs. Everything outside the standard library is
pinned in `go.mod` (exact versions) and `go.sum` (checksums), and is downloaded by the
`go` command itself – there is nothing to install by hand.

| Tool | What it does | Why this one |
|------|--------------|--------------|
| [Go](https://go.dev/) 1.24+ | The language and its `go` command | `go test`, `go vet`, `gofmt`, coverage, fuzzing, benchmarks and the race detector are all built in |
| [`testing`](https://pkg.go.dev/testing) (standard library) | Tests, table-driven subtests, examples, fuzz tests, benchmarks | Every Go project uses it. Our tests are plain `if got != want { t.Errorf(...) }` |
| [go-cmp](https://pkg.go.dev/github.com/google/go-cmp/cmp) | `cmp.Diff(want, got)` for whole receipts | Prints exactly which field differs, instead of two long dumps. Keeps the stdlib style (it is not an assertion library) |
| [godog](https://github.com/cucumber/godog) | Runs the Given/When/Then scenarios in `spec/features/` | The official Cucumber for Go. Runs inside `go test`, so no extra binary |
| [uber-go/mock](https://github.com/uber-go/mock) (`mockgen`, `gomock`) | Generates mocks from our interfaces | The maintained successor of `golang/mock`, the usual Go mocking tool |
| [rapid](https://pkg.go.dev/pgregory.net/rapid) | Property-based testing | Generates and *shrinks* random inputs; `testing/quick` is frozen and cannot shrink |
| [santhosh-tekuri/jsonschema](https://github.com/santhosh-tekuri/jsonschema) v6 | JSON Schema validation (contract tests) | Supports draft 2020-12, which the spec schemas use, and resolves their relative `$ref`s |
| [modernc.org/sqlite](https://pkg.go.dev/modernc.org/sqlite) | SQLite through Go's `database/sql` | Pure Go: no C compiler needed, works the same on every laptop and in CI |
| [staticcheck](https://staticcheck.dev/) | Static checks beyond `go vet` | The most widely used Go linter that needs no configuration |
| [gotestsum](https://github.com/gotestyourself/gotestsum) | Readable test output and a JUnit XML report | CI dashboards read JUnit XML; `go test -json` alone is hard to read |
| [gremlins](https://github.com/go-gremlins/gremlins) | Mutation testing | The only maintained mutation tool for Go. Pinned to v0.5.1, the last one that runs on Go 1.24 |
| [govulncheck](https://go.dev/blog/govulncheck) | Known vulnerabilities in dependencies and the standard library | The Go team's own tool; only reports vulnerable code that is actually called |

`mockgen`, `staticcheck`, `gotestsum` and `govulncheck` are listed in the `tool` block of
`go.mod` (a Go 1.24 feature), so `go tool staticcheck ./...` runs exactly the pinned version.
gremlins is run with `go run …@v0.5.1` instead, so its many dependencies stay out of our
`go.mod`.

## Getting started

You need **Go 1.24 or newer**. Check with `go version`.

```bash
cd go
go test ./...        # downloads the libraries on the first run, then runs the tests
```

You should see `ok` for 5 packages – about 630 tests and subtests, counted by
`go tool gotestsum`. Go caches test results: add `-count=1` to force a real re-run.

## Running the tests and checks

All commands are run from this `go/` folder. `./...` means "this folder and every folder below it".

| I want to… | Command |
|------------|---------|
| run the everyday tests (unit, spec cases, BDD, properties, contracts, golden files, examples) | `go test ./...` |
| see every test name | `go test -v ./...` (or `go tool gotestsum --format testname`) |
| run one package (folder) | `go test ./checkout` |
| run tests whose name matches a pattern | `go test ./checkout -run 'TestBulkDiscount' -v` |
| run one subtest (e.g. one spec case or one table row) | `go test ./checkout -run 'TestSpecCases/promo-coffee-three-bags-one-free' -v` |
| also run the **integration** tests (JSON files, SQLite) | `go test -tags integration ./...` |
| run the **end-to-end** tests (builds and runs the real CLI) | `go test -tags e2e ./e2e/...` |
| run the **shared e2e runner** (from the repository root) | `e2e/run.sh go run ./go/cmd/checkout` |
| run the **BDD** scenarios only | `go test ./bdd -v` |
| run one **property** with more cases | `go test ./checkout -run TestP5 -rapid.checks=2000` |
| check for **data races** | `go test -race -tags integration ./...` |
| prove tests are **independent** (random order) | `go test -shuffle=on ./...` |
| **coverage** with an HTML report | `go test -tags integration -coverpkg=./checkout/...,./cli/...,./jsonstore/...,./sqlitestore/...,./specdir/... -coverprofile=coverage.out ./...` then `go tool cover -html=coverage.out` |
| **fuzz** for 30 seconds | `go test ./checkout -run '^$' -fuzz=FuzzCheckout -fuzztime=30s` |
| run the **benchmark** | `go test ./checkout -run '^$' -bench . -benchmem` |
| run **mutation testing** (about a minute) | see [Mutation testing](#mutation-testing) below |
| update the **golden files** after an intended change | `go test ./checkout ./cli -update` |
| re-generate the **mocks** after changing an interface | `go generate ./...` |
| a **JUnit XML** report, like CI | `go tool gotestsum --junitfile reports/junit.xml -- -tags integration ./...` |
| format, vet and lint | `gofmt -l .` (should print nothing), `go vet ./...`, `go tool staticcheck ./...` |
| check for known vulnerabilities | `go tool govulncheck ./...` |
| try the command line program | `echo '{"items":[{"sku":"APPLE","quantity":1}]}' \| go run ./cmd/checkout` |

**One file?** Go runs tests by *package* (folder) and *name*, not by file. To run the tests
of one file, use `-run` with a pattern that matches their names, e.g.
`go test ./checkout -run 'TestCatalog' -v` for `catalog_test.go`. `-run` takes a regular
expression; a `/` separates a test from its subtests.

## Folder layout

```
go/
├── go.mod, go.sum              ← module name, pinned library and tool versions, checksums
├── checkout/                   ← the checkout package (the "production" code)
│   ├── types.go                ← Product, Promotion, Catalog, Basket, Receipt, …
│   ├── errors.go               ← ErrorCode, ErrorCodes, CheckoutError
│   ├── checkout.go             ← Checkout(basket, catalog): the pipeline from rules §4
│   ├── validation.go           ← rows 1-4 of rules §4.1
│   ├── catalog.go              ← RequireProduct (Catalog.Find is in types.go)
│   ├── lines.go                ← MergeItems, PriceLine
│   ├── promotions.go           ← buy X get Y free, bulk discount
│   ├── coupons.go              ← SAVE10, WELCOME5
│   ├── shipping.go             ← standard / express
│   ├── tax.go                  ← 8%
│   ├── money.go                ← RoundHalfUp, with whole numbers only
│   ├── orders.go               ← Order, Clock, OrderRepository, InMemoryOrderRepository, PlaceOrder
│   ├── *_test.go               ← unit, table-driven, property, contract, golden, example, fuzz, bench tests
│   └── testdata/               ← golden files (the go tool ignores folders named testdata)
├── jsonstore/                  ← order repository: one ORDER-<n>.json file per order
├── sqlitestore/                ← order repository: one table in an SQLite file
├── specdir/                    ← finds ../spec/ by walking up the folder tree
├── cli/                        ← the CLI's logic: JSON text in, output + exit code out
├── cmd/checkout/main.go        ← the CLI itself: stdin → cli.Run → stdout
├── bdd/                        ← godog step definitions for ../spec/features/
├── e2e/                        ← end-to-end tests of the built CLI (build tag e2e)
└── internal/                   ← test support, only importable from inside this module
    ├── spec/                   ← loads ../spec/ files, SchemaErrors, DollarsToCents
    ├── fixtures/               ← a tiny made-up catalog, a sample receipt, a fixed clock
    └── mocks/                  ← generated by mockgen from checkout/orders.go – do not edit
```

The shared data is read straight from the repository's [`spec/`](../spec/) folder (never
copied). A two-line [`go.work`](../go.work) at the repository root tells the `go` command
that this module belongs to the repository, so `go run ./go/cmd/checkout` works from the
root too (that is the command the shared e2e runner uses).

### How Go separates different kinds of tests

Other languages often have a `tests/` folder per kind of test. Go does it differently:

* **Unit tests live right next to the code.** A file named `something_test.go` is only
  compiled by `go test`, never by `go build`.
  - `package checkout` in a test file = *white-box* test (can see unexported names, e.g.
    `errors_test.go` tests the unexported `newError`).
  - `package checkout_test` = *black-box* test (only the public API, like a real user).
    Most tests here are black-box.
* **The kind of test is in the file name and the comment at its top**:
  `properties_test.go`, `contract_test.go`, `golden_test.go`, `fuzz_test.go`, …
* **Build tags hide whole files.** A file that starts with `//go:build integration` (or
  `e2e`) is only compiled when you pass `-tags integration`. That keeps tests that touch
  the disk, a database or a separate process out of the everyday `go test ./...`.

## The kinds of tests, and where to look

| Kind of test | Question it answers | Look at |
|---|---|---|
| **Unit** | Does this one function do the right thing? | `checkout/catalog_test.go` is a good first read, then `checkout/money_test.go` |
| **Parameterised (table-driven)** | Does it work for all these examples? | `[]struct{…}` tables + `t.Run` in `checkout/promotions_test.go`, `validation_test.go`, `tax_test.go` |
| **Shared spec cases** | Does it match the 60 cases every language must pass? | `checkout/spec_cases_test.go` |
| **BDD** | Does it behave as the Given/When/Then scenarios say? | `bdd/steps_test.go` – each sentence from `spec/features/README.md` once; `bdd/features_test.go` runs godog |
| **Test doubles** | How do I test without the real clock? Was `Save` really never called? | `checkout/orders_test.go`: a fake (`InMemoryOrderRepository`), a stub (`FixedClock`), gomock mocks from `internal/mocks` |
| **Integration** | Do the repositories really write files and database rows? | `jsonstore/jsonstore_integration_test.go`, `sqlitestore/sqlitestore_integration_test.go` |
| **Property-based** | Is it true for *every* basket, not just my examples? | `checkout/properties_test.go` – `TestP1_…` to `TestP14_…` |
| **Contract** | Is the output even the right shape? | `checkout/contract_test.go` (receipts), `cli/contract_test.go` (CLI output and errors) |
| **Snapshot / golden files** | Did the output change since last time? | `checkout/golden_test.go` + `checkout/testdata/`, `cli/golden_test.go` |
| **End-to-end** | Does the real program work from the outside? | `e2e/cli_test.go`, then `e2e/run.sh go run ./go/cmd/checkout` from the root |
| **Examples (documentation tests)** | Do the examples in the docs still work? | `checkout/example_test.go` – shown by `go doc` and checked by `go test` |
| **Fuzzing** | Can *any* input make it crash? (property P10) | `checkout/fuzz_test.go` |
| **Benchmark** | How fast is it? | `checkout/bench_test.go` |
| **Concurrency** | Is the repository safe to use from many goroutines? | the `…ManyGoroutines…` tests, run with `-race` |
| **Coverage** | Which lines did the tests never run? | 100% of `checkout`, `cli`, `jsonstore`, `sqlitestore` and `specdir`; CI fails below 100% |
| **Mutation** | Would the tests notice if the code were wrong? | see below |
| **Static checks** | Can we catch mistakes without running anything? | `gofmt`, `go vet`, `go tool staticcheck` |

### Golden files: why hand-rolled

The matrix offers two options: [go-snaps](https://github.com/gkampitakis/go-snaps) or golden
files with an `-update` flag. We picked **hand-rolled golden files** because they teach the
most: the whole mechanism is about twenty lines in `checkout/golden_test.go` (`flag.Bool`,
`os.ReadFile`, `os.WriteFile`, `cmp.Diff`), the approved output is an ordinary JSON file you
can read in `testdata/`, and the same pattern is used by the Go standard library itself.
After an intended change, run `go test ./checkout ./cli -update` and **read the diff** before
committing it.

### Mutation testing

gremlins makes small changes to the code ("mutants": `>=` becomes `>`, `+` becomes `-`,
a condition is negated, `&&` becomes `||`, `break` becomes `continue`, …) and runs the
tests against each one. A mutant that survives means no test noticed the change.

```bash
CHECKOUT_SPEC_DIR="$(cd .. && pwd)/spec" \
go run github.com/go-gremlins/gremlins/cmd/gremlins@v0.5.1 unleash \
  --tags integration \
  --exclude-files 'cmd/' --exclude-files 'internal/mocks/' \
  --invert-logical --invert-loopctrl --invert-assignments --remove-self-assignments \
  --threshold-efficacy 95 --threshold-mcover 95 --output reports/mutation.json
```

Why `CHECKOUT_SPEC_DIR`? gremlins copies `go/` into a temporary folder and runs the tests
there – and from there, no `spec/` folder can be found by looking up the folder tree. The
variable names it directly (see package `specdir`).

This is worth a story, because it fooled us once. Without the variable, **every** test run in
the copy failed – and gremlins counts any failing test run as a "killed" mutant. So the first
report said "120 killed, 0 lived, 100%" while the tests had not checked a single mutant. A
review caught it: a line that no test looks at was reported as killed too. Now CI first runs
the whole suite in a copy outside the repository, and only then gremlins. The lesson: a tool
that can only ever say "all good" is not testing anything – check that it can fail.

The real result is **118 mutants, 118 killed, 0 lived: 100% test efficacy and 100%
mutator coverage**, in about a minute, so CI runs it on every pull request and fails below
95%. Getting there shaped both the tests and the code:

* Boundaries are tested on both sides (`2499` / `2500`, `4999` / `5000`, `9` / `10`,
  `99` / `100`, `0` / `1`), because `>` versus `>=` is the classic surviving mutant.
* Error **messages** are tested, not only error codes.
* The first real run found three survivors. Turning `continue` into `break` in the JSON-file
  store went unnoticed, because no test put a file *before* the orders in the folder (now
  `AAA-notes.txt` and an `ARCHIVE` folder do). `n > highest` versus `n >= highest` can never be
  told apart, so the line became `highest = max(highest, n)`. And nothing checked that a
  program started through a symlink looks for `spec/` next to the real program, so that is
  now its own small function, `programFolder`, tested with a real symlink.
* Branches that cannot happen were removed instead of left untested: `json.Marshal` of a
  receipt cannot fail, so its error is ignored with a comment; the SQLite store uses
  `INSERT … RETURNING` instead of a second call with its own error.
* `cmd/checkout/main.go` is excluded: it runs in a separate process (see `e2e/`), which an
  in-process tool cannot observe. Its logic is in package `cli`, which *is* mutated.
* gremlins moves slowly (see the matrix); v0.6.0 needs Go 1.25, so we pin v0.5.1. It also
  mutates the test-support package `internal/spec`, which is fine – those mutants are killed too.

## How to add a new test

Say you want to check that `RoundHalfUp` rounds 10% of $0.05 (0.5 cents) up to 1 cent.

1. **Pick the kind of test.** It is one function with made-up numbers: a unit test, so open
   `checkout/money_test.go` (or create a new file ending in `_test.go` next to the code).
2. **Write the test**, following the *arrange / act / assert* pattern:

   ```go
   func TestRoundHalfUpRoundsHalfACentUpToOneCent(t *testing.T) {
       // Arrange
       amount := 5

       // Act
       rounded := checkout.RoundHalfUp(amount, 10)

       // Assert
       if rounded != 1 {
           t.Errorf("RoundHalfUp(5, 10) = %d, want 1", rounded)
       }
   }
   ```

   Many similar cases? Add a row to a **table** instead (see `TestRoundHalfUpMoreExamples`).
3. **Run it:** `go test ./checkout -run TestRoundHalfUpRoundsHalfACentUpToOneCent -v`.
4. **Make sure it can fail.** Change `1` to `0` and check that the test goes red.
   A test that never fails does not test anything! (Mutation testing automates this idea.)
5. **Run everything** before you open a pull request: `gofmt -l .`, `go vet ./...`,
   `go tool staticcheck ./...` and `go test -race -tags integration ./...`.

Want to add a case for *every* language instead? Add it to
[`spec/test-cases/`](../spec/test-cases/) and it is picked up automatically by the spec
cases, contract, e2e and shared-runner tests (and update the count in
`internal/spec/spec_test.go`).

## Things to know about this implementation

* **Whole numbers only.** `RoundHalfUp(amount, percent)` is `(amount*percent + 50) / 100`,
  as the hint in rules §7 suggests. Go's integer division does the rounding down.
* **Errors are values.** `Checkout` returns `(Receipt, error)`; the error is always a
  `*CheckoutError` with a `Code` (for programs) and a `Message` (for people). Callers use
  `errors.As`. A missing coupon, JSON `null` and `""` all arrive as `""` in Go – which is
  exactly what the rules want, since they all mean "no coupon".
* **Parse, don't validate.** `ParseShippingMethod` turns the basket's free-form string into
  the type `ShippingMethod` or an error; after that the value is known to be fine.
* **The clock is an interface** with one method, `Now()`. `ClockFunc` turns a plain function
  into a clock; `SystemClock()` is the real one and `FixedClock(t)` the test stub.
* **Repositories are an interface** (`Save` and `Find`) with three implementations. `Find`
  returns `nil, nil` for "not found" (a normal answer) and an error only for real trouble.
  The in-memory one has a mutex, so it is safe from several goroutines (checked with `-race`).
* **`PlaceOrder` = `Checkout` + `Save`.** It cannot save anything when checkout fails, which
  is what the gomock test and property P14 check.
* **The CLI** exits `0` with a receipt, `1` with `{"error":"<CODE>"}`, and `2` with a message
  on stderr (and nothing on stdout) when the input is not a basket at all, or when it cannot
  find `spec/catalog.json`. It looks for `spec/` above the working directory first, then above
  the program file, so it works from anywhere inside the repository.

## Ideas for the next contribution

Every row of the Go column of the matrix is covered, so what is left is polish:

* Compare with another language: same spec, same kinds of tests, different idioms. That
  comparison is the point of the project.
* Try [go-snaps](https://github.com/gkampitakis/go-snaps) next to the hand-rolled golden
  files and write down the trade-off.
* [golangci-lint](https://golangci-lint.run/) bundles staticcheck with many more linters;
  it is a bigger tool to configure, which is why we started with staticcheck alone.
* A state-machine (model-based) test with rapid's `T.Repeat`: random sequences of `Save`
  and `Find` against the JSON-file store, checked against the in-memory fake as a model.
