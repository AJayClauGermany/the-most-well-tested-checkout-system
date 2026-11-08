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
| **Test framework** | pytest | Vitest | `testing` (built in) | JUnit | xUnit v3 | built-in test harness |
| **Static checks** | Ruff (lint + format) | `tsc` type check | `gofmt`, `go vet` | `javac -Xlint:all -Werror` | `dotnet format`, warnings as errors | `rustfmt`, Clippy |
| **Run all tests** | `pytest` | `npm test` | `go test ./...` | `mvn verify` | `dotnet test` | `cargo test` |
| **Entry point** | `checkout(basket, catalog)` | `checkout(basket, catalog)` | `checkout.Checkout(basket, catalog)` | `Checkout.checkout(basket, catalog)` | `Checkout.Calculate(basket, catalog)` | `checkout(&basket, &catalog)` |
| **Errors** | raises `CheckoutError` | throws `CheckoutError` | returns `*CheckoutError` | throws `CheckoutException` | throws `CheckoutException` | returns `Err(CheckoutError)` |

## Where do the different kinds of tests live?

| | Python | TypeScript | Go | Java | C# | Rust |
|---|---|---|---|---|---|---|
| **Unit** | `tests/unit/` | `tests/unit/` | `*_test.go` next to the code | `src/test/java/…/*Test.java` (Surefire) | `tests/Checkout.UnitTests/` | `#[cfg(test)] mod tests` inside each source file |
| **Integration** | `tests/integration/` | `tests/integration/` | `*_test.go` + build tag | `*IT.java` (Failsafe) | a separate test project | `tests/*.rs` |
| **Shared spec cases** | `tests/acceptance/` | `tests/integration/` | `checkout/spec_cases_test.go` | `SharedSpecCasesTest.java` | `SpecCaseTests.cs` | `tests/spec_cases.rs` |
| **End-to-end** | `tests/e2e/` | `tests/e2e/` | `e2e/` (build tag `e2e`) | `e2e/` package, `*IT.java` | a separate test project | `tests/` |

## How "not implemented yet" tests are skipped

Every language runs all 60 shared spec cases against the checkout, but skips them until
the checkout exists. Each one does it differently:

| Language   | How a test is skipped                                  |
|------------|--------------------------------------------------------|
| Python     | `@pytest.mark.skip(reason=...)`                         |
| TypeScript | `it.skip.each(cases)(...)`                              |
| Go         | `t.Skip("...")` inside each subtest                     |
| Java       | `@Disabled("...")`                                      |
| C#         | `[MemberData(..., Skip = "...")]`                       |
| Rust       | `#[ignore = "..."]` (run anyway with `cargo test -- --ignored`) |

## Testing showcase tracker

What each language demonstrates today. ✅ = done, ⏸️ = written but skipped until the
checkout is implemented, blank = waiting for a contributor. Update this table in your
pull request when you add something!

| Kind of test                         | Python | TypeScript | Go | Java | C# | Rust |
|--------------------------------------|:------:|:----------:|:--:|:----:|:--:|:----:|
| Unit tests                           | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| Loading the shared spec files        | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| Shared spec cases (data-driven)      | ⏸️ | ⏸️ | ⏸️ | ⏸️ | ⏸️ | ⏸️ |
| BDD with `spec/features`             |    |    |    |    |    |    |
| Integration (e.g. saving orders to a file) |    |    |    | ✅¹ |    |    |
| Test doubles (fakes / mocks)         |    |    |    |    |    |    |
| Property-based                       |    |    |    |    |    |    |
| Snapshot / approval                  |    |    |    |    |    |    |
| End-to-end via the CLI               |    |    |    |    |    |    |
| Code coverage                        |    |    |    |    |    |    |
| Mutation testing                     |    |    |    |    |    |    |

¹ A single example `*IT` test that proves the Failsafe wiring works.

Each language's `README.md` lists concrete tools we would love to see added.
