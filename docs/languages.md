---
layout: page
title: Languages
---

All six languages solve the same problem and load the same shared spec files, but
each one organises its tests in its own way. Comparing them side by side is half the fun.

## At a glance

| | Python | TypeScript | Go | Java | C# | Rust |
|---|---|---|---|---|---|---|
| **Folder** | [`python/`]({{ site.repository_files }}/python/) | [`typescript/`]({{ site.repository_files }}/typescript/) | [`go/`]({{ site.repository_files }}/go/) | [`java/`]({{ site.repository_files }}/java/) | [`csharp/`]({{ site.repository_files }}/csharp/) | [`rust/`]({{ site.repository_files }}/rust/) |
| **Build tool** | pip + `pyproject.toml` | npm | Go modules | Maven | .NET SDK | Cargo |
| **Test framework** | pytest | Vitest | `testing` (built in) | JUnit 6 | xUnit v3 | built-in test harness (+ cargo-nextest) |
| **Static checks** | Ruff (lint + format), mypy (strict) | `tsc` type check, Biome (lint + format) | `gofmt`, `go vet`, staticcheck | `javac -Xlint:all -Werror`, Error Prone, Checkstyle | `dotnet format`, analyzers, warnings as errors | `rustfmt`, Clippy |
| **Run all tests** | `pytest` | `npm test` | `go test -tags integration ./...` | `mvn verify` | `dotnet test` | `cargo test` |
| **Entry point** | `checkout(basket, catalog)` | `checkout(basket, catalog)` | `checkout.Checkout(basket, catalog)` | `Checkout.checkout(basket, catalog)` | `Checkout.Calculate(basket, catalog)` | `checkout(&basket, &catalog)` |
| **Errors** | raises `CheckoutError` | throws `CheckoutError` | returns `*CheckoutError` | throws `CheckoutException` | throws `CheckoutException` | returns `Err(CheckoutError)` |
| **Checkout + save** | `place_order(...)` | `placeOrder(...)` | `checkout.PlaceOrder(...)` | `Orders.placeOrder(...)` | `Checkout.PlaceOrder(...)` | `place_order(...)` |
| **Replaceable clock** (rules §6) | a function returning a `datetime` | a function returning a `Date` | a `Clock` interface | `java.time.Clock` | .NET's own `TimeProvider` | a `Clock` trait |
| **Command line program** | `python -m checkout` | `node src/cli.ts` | `go run ./cmd/checkout` | `java -jar target/checkout.jar` | `dotnet run --project src/Checkout.Cli` | `cargo run --bin checkout` |

## Where do the different kinds of tests live?

Each language follows its own community's habits, so the same kind of test ends up in very
different places:

| | Python | TypeScript | Go | Java | C# | Rust |
|---|---|---|---|---|---|---|
| **Unit** | `tests/unit/` | `tests/unit/` | `*_test.go` next to the code | `src/test/java/…/*Test.java` (Surefire) | `tests/Checkout.UnitTests/` | `#[cfg(test)] mod tests` inside each source file |
| **Shared spec cases** | `tests/acceptance/` | `tests/integration/` | `checkout/spec_cases_test.go` | `SharedSpecCasesTest.java` | `tests/Checkout.UnitTests/SpecCaseTests.cs` | `tests/spec_cases.rs`, `tests/spec_cases_named.rs` |
| **BDD** | `tests/bdd/` | `tests/bdd/` | `bdd/` | `bdd/` package | `tests/Checkout.Specs/` | `tests/features.rs` |
| **Integration** | `tests/integration/` | `tests/integration/` | `*_integration_test.go` + build tag `integration` | `*IT.java` (Failsafe) | `tests/Checkout.IntegrationTests/` | `tests/*_repository.rs` |
| **Property-based** | `tests/property/` | `tests/property/` | `checkout/properties_test.go` | `property/` package | `tests/Checkout.PropertyTests/` | `tests/properties.rs` |
| **Contract** | `tests/contract/` | `tests/contract/` | `checkout/contract_test.go`, `cli/contract_test.go` | `contract/` package | `tests/Checkout.ContractTests/` | `tests/contract.rs` |
| **Snapshot** | `tests/snapshot/` | `tests/snapshot/` | `checkout/golden_test.go` + `testdata/` | `snapshot/` package | `tests/Checkout.SnapshotTests/` | `tests/snapshots.rs` |
| **End-to-end** | `tests/e2e/` | `tests/e2e/` | `e2e/` (build tag `e2e`) | `e2e/` package, `*IT.java` | `tests/Checkout.EndToEndTests/` | `tests/cli.rs` |

In short: Python and TypeScript use one folder per kind of test; Go and Rust keep unit tests
right next to the code; Java decides by class name (`*Test` or `*IT`); C# gives each kind its
own test project.

## From skipped to green

Every language was built test-first: the commit that adds the tests comes before the commit
that makes them pass. Before that, all 60 shared spec cases were written but **skipped**, so
the build stayed green while the work waited. Each language skips a test in its own way – and
each one has since removed its skip:

| Language   | How the shared spec cases were skipped |
|------------|----------------------------------------|
| Python     | `@pytest.mark.skip(reason=...)` |
| TypeScript | `it.skip.each(cases)(...)` |
| Go         | `t.Skip("...")` inside each subtest |
| Java       | `@Disabled("...")` |
| C#         | `[MemberData(..., Skip = "...")]` |
| Rust       | `#[ignore = "..."]` (run anyway with `cargo test -- --ignored`) |

Run `git log --oneline -- <language folder>` to see the order: the tests, then the
implementation – and, in several languages, later commits where coverage or mutation testing
found gaps in the tests (Python, Go and Rust have good examples).

## Which kinds of tests exist in which language?

That is tracked in one place, [the testing matrix](testing-matrix.html): every kind of
test we show, and the tool chosen for it in each language. Each language's `README.md`
explains its own column in more detail.
