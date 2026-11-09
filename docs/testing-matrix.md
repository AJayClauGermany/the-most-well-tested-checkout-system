---
layout: page
title: The testing matrix
---

This page is the **goal of the project** in one table: which kinds of testing we want to
show, with which tool, in each of the six languages. Everything else (the language READMEs,
the issue tracker, the [kinds of tests](kinds-of-tests.html) page) points here, so when the
goal changes, change it **here first**.

The guiding rule is still *keep it simple*: one small program, mainstream tools a beginner
can run on a laptop, and no kind of test that has nothing real to bite on in a checkout.

**Legend:** ✅ done · ◻️ wanted – pick it up! · – not idiomatic in this language, and that
is fine.
Where two tools are named, the first is the one we recommend; the second is a reasonable
alternative.

## Tier 1 – every language does these

The core. A language is "complete" when every row below is ✅ – and all six languages are
now complete. Each cell names the tool that language uses; the language README says why.

| Kind of test | Shared input | Python | TypeScript | Go | Java | C# | Rust |
|---|---|---|---|---|---|---|---|
| **Unit tests** | – | ✅ pytest | ✅ Vitest | ✅ `testing` | ✅ JUnit 6 + AssertJ | ✅ xUnit v3 | ✅ built-in harness |
| **Parameterised tests** (the language's native table style) | – | ✅ `@pytest.mark.parametrize` | ✅ `it.each` | ✅ table-driven `t.Run` subtests | ✅ `@ParameterizedTest` | ✅ `[Theory]` + `[InlineData]` / `[MemberData]` | ✅ `rstest` `#[case]` |
| **Shared spec cases** (data-driven) | `spec/test-cases/` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| **BDD** (Given / When / Then) | `spec/features/` | ✅ pytest-bdd | ✅ `@cucumber/cucumber` | ✅ godog | ✅ Cucumber-JVM 8 via `cucumber-junit-platform-engine` | ✅ Reqnroll | ✅ `cucumber` crate |
| **Test doubles** – a hand-written fake *and* a mock of the order repository and the clock | rules §6 | ✅ hand-written fake, `unittest.mock`, `monkeypatch` | ✅ hand-written fake, `vi.fn()`, fake timers | ✅ hand-written fake; `uber-go/mock` | ✅ hand-written fake, `Clock.fixed`, Mockito | ✅ hand-written fake, `FakeTimeProvider`, NSubstitute | ✅ hand-written fake; `mockall` |
| **Integration** – JSON-file and SQLite order repositories against a temp folder | `spec/schema/order.schema.json` | ✅ `tmp_path` + `sqlite3` | ✅ `node:fs` + `os.tmpdir()` + `node:sqlite` | ✅ `t.TempDir()` + `modernc.org/sqlite` | ✅ JDBC + H2 (Failsafe `*IT`) | ✅ temp folder + `Microsoft.Data.Sqlite` | ✅ `tempfile` + `rusqlite` |
| **Property-based** | `spec/properties.md` | ✅ Hypothesis | ✅ fast-check | ✅ rapid | ✅ jqwik 1.9.x ¹ | ✅ FsCheck (`FsCheck.Xunit.v3`) | ✅ proptest |
| **Contract (schema) tests** of the receipt | `spec/schema/receipt.schema.json` | ✅ `jsonschema` | ✅ ajv | ✅ `santhosh-tekuri/jsonschema` | ✅ networknt json-schema-validator | ✅ JsonSchema.Net | ✅ `jsonschema` crate |
| **Snapshot / approval** of whole receipts | `spec/examples/` | ✅ syrupy | ✅ `toMatchSnapshot` / inline | ✅ golden files with `-update` | ✅ ApprovalTests.Java | ✅ Verify (`Verify.XunitV3`) ⁵ | ✅ insta |
| **End-to-end** through the CLI, from inside the language | rules §8 | ✅ `subprocess` | ✅ `node:child_process` | ✅ `os/exec` | ✅ `ProcessBuilder` | ✅ `System.Diagnostics.Process` | ✅ `assert_cmd` |
| **Code coverage** | – | ✅ pytest-cov | ✅ `@vitest/coverage-v8` | ✅ `go test -cover` | ✅ JaCoCo | ✅ `Microsoft.Testing.Extensions.CodeCoverage` + ReportGenerator | ✅ cargo-llvm-cov |
| **Mutation testing** | – | ✅ mutmut ² | ✅ StrykerJS (Vitest runner) | ✅ gremlins ³ | ✅ PIT (pitest) | ✅ Stryker.NET (MTP runner, preview) | ✅ cargo-mutants ⁶ |
| **Static checks** | – | ✅ Ruff · ✅ mypy (strict) | ✅ `tsc` · ✅ Biome ⁴ | ✅ `gofmt`, `go vet` · ✅ staticcheck · ✅ govulncheck ⁷ | ✅ `-Xlint:all -Werror` · ✅ Error Prone · ✅ Checkstyle | ✅ analyzers (`latest-recommended`), warnings as errors | ✅ rustfmt + Clippy |

¹ jqwik is in maintenance mode and its 1.10 release added an unusual "anti-AI usage"
licence clause; we pin 1.9.x (plain Apache 2.0) and note that Java currently has no
other mainstream property-based library.
² mutmut uses `fork()`, so on Windows run it under WSL.
³ gremlins is the only maintained Go mutation tool and it moves slowly; expect rough edges.
⁴ The TypeScript folder pins TypeScript 7, the native (Go) compiler, which has no JavaScript
compiler API. typescript-eslint needs that API, so the one-binary alternative Biome does
linting and formatting instead.
⁵ Pinned to Verify 32.x: from 33.0 Verify fails the build unless a yearly sponsorship (or
exemption) property is set.
⁶ A full cargo-mutants run takes about fifteen minutes, so pull requests mutate only the lines they
change (`--in-diff`); starting the Rust workflow by hand runs the full set.
⁷ govulncheck runs in CI as an informational step: a new advisory against the pinned Go
toolchain should not turn an unrelated pull request red. It is the one small exception to
"security scanning" being out of scope (see below): Go ships it, and it costs one line.

## Tier 2 – where the language has an idiomatic story

These are not required everywhere. A "–" means the language has no natural tool for it,
and forcing one in would make the example *less* typical of that language.

| Kind of test | Python | TypeScript | Go | Java | C# | Rust |
|---|---|---|---|---|---|---|
| **Documentation tests** (examples in docs that also run) | ✅ `doctest` | – | ✅ `Example` functions with `// Output:` | – | – | ✅ doctests (`cargo test --doc`) |
| **Fuzzing** (random bytes/JSON must never crash – property P10) | – (Hypothesis covers it) | – | ✅ built-in `go test -fuzz` | ✅ Jazzer `@FuzzTest` | – | ✅ cargo-fuzz target ⁸ |
| **Benchmarks** (one micro-benchmark of `checkout`) | – | ✅ `vitest bench` | ✅ `testing.B` with `b.Loop` | – | ✅ BenchmarkDotNet | ✅ criterion |
| **Architecture tests** ("the domain never imports the CLI / IO") | – | – | – | ✅ ArchUnit | ✅ ArchUnitNET | – |
| **Type-level tests** | – | ✅ `expectTypeOf` + `vitest --typecheck` | – | – | – | – |
| **Concurrency** (the repository used from several goroutines) | – | – | ✅ `go test -race` | – | – | – |

⁸ cargo-fuzz needs a nightly compiler, so CI only checks that the fuzz target compiles;
`cargo +nightly fuzz run basket_checkout` fuzzes it locally. Go and Java fuzz from inside
their normal test runners, and their saved inputs replay on every run.

## Tier 3 – shared plumbing

Not kinds of tests, but things every language should show once because every real
project has them.

| What | Python | TypeScript | Go | Java | C# | Rust |
|---|---|---|---|---|---|---|
| **Test report for CI** (JUnit XML or similar) | ✅ `--junitxml` | ✅ `--reporter=junit` | ✅ gotestsum `--junitfile` | ✅ Surefire / Failsafe (default) | ✅ `--report-xunit-junit` + TRX | ✅ `cargo nextest --profile ci` |
| **Tests run in parallel / independently** | ✅ pytest-xdist, pytest-randomly | ✅ Vitest (default) | ✅ `t.Parallel()`, `-shuffle=on` | ✅ JUnit parallel config | ✅ xUnit collections | ✅ cargo-nextest |

And two things that are shared by all six and live outside the language folders:

| What | Where | Status |
|---|---|---|
| Spec validation (schemas, unique ids, example outputs) | [`.github/workflows/spec.yml`]({{ site.repository_files }}/.github/workflows/spec.yml) | ✅ |
| One e2e runner for every language's CLI | [`e2e/run.sh`]({{ site.repository_files }}/e2e/run.sh) | ✅ runner · ✅ all six CLIs |

## Deliberately out of scope

The checkout is deterministic, has no network, no screen and no users logging in. These
kinds of testing have nothing real to test here, so they would only add noise:

* UI / browser testing (Playwright, Selenium)
* Load and performance testing beyond a single micro-benchmark
* Contract testing between services (Pact) – there is only one program
* Security scanning and penetration testing (Go's built-in, informational `govulncheck` is the
  one exception – footnote ⁷)
* Flaky-test retries – a test here that is flaky is simply a bug
* A real database server behind Docker (Testcontainers) – SQLite shows the same idea
  and runs on any laptop

**One possible stretch goal, later:** wrap the checkout in a small HTTP API, which would
make *API testing* (request/response contracts, OpenAPI validation) meaningful. It is
explicitly not part of the current goal; the matrix above comes first.

## Order stores used by the integration tests

The rules ([§6]({{ site.repository_files }}/spec/checkout-rules.md#6-saving-an-order-optional))
say what saving an order *does*, not where orders are kept. Which stores to build is a
testing decision, so it is recorded here:

| Store | Why it is here | Kind of test it shows |
|---|---|---|
| **In memory** | The simplest thing that works, and a hand-written *fake* for tests of the checkout itself | Unit tests, test doubles |
| **JSON file** (orders written as [`order.schema.json`]({{ site.repository_files }}/spec/schema/order.schema.json)) | Real files in a real temporary folder | Integration tests, contract tests |
| **SQLite** (the language's usual embedded library) | A real database with no server to install | Integration tests |

A database server behind Docker (Testcontainers) is deliberately not on the list: SQLite
shows the same idea and runs on any laptop.

## Order of work

Every language got here the same way, and its git history shows it: first one commit with
the **whole test suite** (red – the code was only stubs), then the commit that makes it green,
and, in several languages, later commits where coverage or mutation testing found gaps in
the tests.

1. Implement the checkout, test-first, until the shared spec cases are ✅ (they unblock
   almost every other row).
2. Add saving orders with a replaceable time source (rules §6) – this unlocks test
   doubles and integration tests.
3. Add the CLI (rules §8) – this unlocks e2e tests and the shared runner.
4. Fill in the rest of Tier 1, then Tier 2 and 3 where they fit.

The matrix is complete, so new work is a new row (a kind of test we do not show yet), a better
tool for an existing cell, or a new shared case or property that every language then has to
pass. Update this table in the same pull request that adds a test. Each language README's
*"Testing ideas"* section mirrors its column and adds the practical hints.
