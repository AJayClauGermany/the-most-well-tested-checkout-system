# Checkout – C# (.NET)

The C# version of the checkout system described in [`../spec/checkout-rules.md`](../spec/checkout-rules.md).

> **Status: done.** `Checkout.Calculate`, saving orders (in memory, as JSON files and in
> SQLite) and the command line program are implemented, and every row of the C# column of
> the [testing matrix](../docs/testing-matrix.md) is ticked. It was built test-first: the
> commit that adds the tests comes before the commit that makes them pass – have a look at
> the git history.

## Tools we use (and why)

| Tool | What it does | Why this one |
|------|--------------|--------------|
| [.NET 10 SDK](https://dotnet.microsoft.com/download) | The C# compiler and the `dotnet` command | Builds, formats and tests everything; the version is pinned in [`global.json`](global.json) |
| [xUnit.net v3](https://xunit.net/) | The test framework: `[Fact]`, `[Theory]`, `Assert.Equal(...)` | The most widely used test framework for modern .NET |
| [Microsoft.Testing.Platform](https://learn.microsoft.com/dotnet/core/testing/microsoft-testing-platform-intro) | The engine `dotnet test` uses to run the tests, plus its coverage extension | The new, faster .NET test runner; switched on in `global.json` |
| [Reqnroll](https://reqnroll.net/) (`Reqnroll.xunit.v3`) | Runs the Given/When/Then scenarios in `spec/features/` | The open-source successor of SpecFlow, the standard Gherkin tool for .NET |
| [FsCheck](https://fscheck.github.io/FsCheck/) (`FsCheck.Xunit.v3`) | Property-based testing | The best-known property-based library for .NET; `[Property]` works like `[Fact]` |
| [NSubstitute](https://nsubstitute.github.io/) | Mocks: "was `Save` called? With what?" | Reads like English (`repository.DidNotReceive().Save(...)`). We skip Moq (its 2023 SponsorLink episode) |
| [`Microsoft.Extensions.TimeProvider.Testing`](https://learn.microsoft.com/dotnet/standard/datetime/timeprovider-overview) | `FakeTimeProvider`: a clock that says what the test tells it | Fakes .NET's own `TimeProvider`, so the production code needs no clock interface of its own |
| [JsonSchema.Net](https://docs.json-everything.net/schema/basics/) | JSON Schema validation (contract tests) | Supports draft 2020-12, which the spec schemas use |
| [Verify](https://github.com/VerifyTests/Verify) (`Verify.XunitV3`) | Snapshot / approval tests | The standard snapshot library for .NET. Pinned to 32.0.1, the last version before a "SponsorCheck" build gate that fails the build unless a sponsorship property is set |
| [Microsoft.Data.Sqlite](https://learn.microsoft.com/dotnet/standard/data/sqlite/) | SQLite for the SQLite order repository | Made by the .NET team; brings SQLite along, nothing to install |
| [ArchUnitNET](https://github.com/TNG/ArchUnitNET) | Architecture tests ("the checkout never touches files") | Checks the rules against the compiled assemblies |
| [BenchmarkDotNet](https://benchmarkdotnet.org/) | A micro-benchmark | The standard benchmarking tool for .NET |
| [ReportGenerator](https://reportgenerator.io/) | Merges the coverage files into one HTML report | Reads the Cobertura files that `dotnet test --coverage` writes |
| [Stryker.NET](https://stryker-mutator.io/docs/stryker-net/introduction/) | Mutation testing | The mutation testing tool for .NET; runs on Microsoft.Testing.Platform (preview) |
| `dotnet format` + .NET analyzers | Static checks: style, formatting and code analysis (rules in [`.editorconfig`](.editorconfig)) | Built into the SDK. `latest-recommended` analysis, checked at build time, warnings are errors |

Every NuGet package is pinned to an exact version in the `.csproj` files, and every project
has a `packages.lock.json` that records every package it pulls in (`dotnet restore
--locked-mode` in CI fails if one is out of date). ReportGenerator and Stryker.NET are local
.NET tools, pinned in [`.config/dotnet-tools.json`](.config/dotnet-tools.json).

## Getting started

You need the **.NET 10 SDK** (<https://dotnet.microsoft.com/download>). Check it with
`dotnet --version` (10.0.100 or newer). Then, from the repository root:

```bash
cd csharp
dotnet tool restore   # installs ReportGenerator and Stryker.NET from .config/dotnet-tools.json
dotnet test           # builds and runs every test project
```

You should see `total: 552`, `failed: 0`. Any editor works; Visual Studio, JetBrains Rider and
VS Code with the *C# Dev Kit* extension can also run and debug single tests with one click.

## Running the tests and checks

All commands are run from this `csharp/` folder. `dotnet test` finds `Checkout.slnx` (the
solution file that lists every project) and runs every test project in it.

| I want to… | Command |
|------------|---------|
| run all tests | `dotnet test` |
| see every test name | `dotnet test --output detailed` |
| run one kind of test (one test project) | `dotnet test --project tests/Checkout.UnitTests` |
| run one test class | `dotnet test --project tests/Checkout.UnitTests --filter-class "*.CouponsTests"` |
| run one test method | `dotnet test --project tests/Checkout.UnitTests --filter-method "*.Discount_Welcome5OneCentBelowTheMinimum_ThrowsCouponMinimumNotMet"` |
| run one shared spec case | `dotnet test --project tests/Checkout.UnitTests --filter "DisplayName~promo-coffee-three-bags-one-free"` |
| run the **BDD** scenarios from `spec/features/` | `dotnet test --project tests/Checkout.Specs` |
| run the **property-based** tests (P1–P14) | `dotnet test --project tests/Checkout.PropertyTests` |
| run the tests with **coverage** and JUnit/TRX **reports** | see [Coverage](#coverage) below |
| run **mutation testing** (about five minutes) | `dotnet stryker --output reports/mutation` → `reports/mutation/reports/mutation-report.html` |
| run the **shared e2e runner** (from the repository root) | `e2e/run.sh dotnet run --project csharp/src/Checkout.Cli` |
| … faster, after a `dotnet build` | `e2e/run.sh dotnet csharp/src/Checkout.Cli/bin/Debug/net10.0/Checkout.Cli.dll` |
| run the **benchmark** | `dotnet run -c Release --project benchmarks/Checkout.Benchmarks -- --job short --filter '*'` |
| check formatting (as CI does) | `dotnet format --verify-no-changes` |
| fix formatting | `dotnet format` |
| try the command line program | `echo '{"items":[{"sku":"APPLE","quantity":1}]}' \| dotnet run --project src/Checkout.Cli` |

The shared runner with `dotnet run` takes about three minutes, because `dotnet run` checks
the build before each of the 60 cases; the `dotnet <dll>` form takes about ten seconds and is
what CI uses. Both run the same program.

### How filtering works

xUnit gives `dotnet test` a few ways to pick tests. Use **one kind** of filter per command,
and name the test project that holds the tests with `--project`: without it, the filter runs in
every test project, and each project where it matches nothing ends with exit code 8 ("zero tests
ran"), so the whole run says `Failed!` although nothing failed (add `--ignore-exit-code 8` if you
really want to filter across all projects).

* **Simple filters** (xUnit's own): `--filter-class`, `--filter-method`, `--filter-namespace`,
  `--filter-trait`. They take full names, and `*` is a wildcard at the start or end:
  `--filter-class "*SpecLoaderTests"`, `--filter-method "*FindBySku*"`.
  Repeat an option to allow several values. There are `--filter-not-…` versions too.
* **VSTest-style filter** (the classic .NET syntax you will find in many tutorials):
  `--filter "<property><operator><value>"`, where the property is `FullyQualifiedName`,
  `DisplayName` or a trait name such as `Category`, and the operator is
  `=` (equals), `!=`, `~` (contains) or `!~`. Combine with `&` (and) and `|` (or):

  ```bash
  dotnet test --project tests/Checkout.UnitTests --filter "FullyQualifiedName~CouponsTests"
  dotnet test --project tests/Checkout.UnitTests --filter "DisplayName~validation-"   # every validation spec case
  dotnet test --project tests/Checkout.PropertyTests --filter "DisplayName~P10"       # one property
  ```

**One file?** .NET runs tests by *name*, not by file. By convention one file holds one
test class with the same name (`CouponsTests.cs` → `class CouponsTests`), so
"run one file" means "run one class".

### Coverage

```bash
dotnet test --coverage --coverage-output-format cobertura \
  --coverage-settings "$PWD/codecoverage.settings.xml" \
  --results-directory "$PWD/reports/test-results" \
  --report-xunit-junit --report-xunit-trx
dotnet reportgenerator "-reports:reports/test-results/*.cobertura.xml" -targetdir:reports/coverage \
  "-reporttypes:Html;TextSummary;MarkdownSummaryGithub;Cobertura" \
  "-classfilters:-System.Text.RegularExpressions.Generated*"
./check-coverage.sh      # fails below 100% line / 95% branch coverage
```

Then open `reports/coverage/index.html`. Three things worth knowing:

* Each test project writes its own coverage file; ReportGenerator merges them.
  [`codecoverage.settings.xml`](codecoverage.settings.xml) limits the measurement to our three
  `src/` assemblies and measures **in memory** – the default rewrites `Checkout.dll` on disk,
  and the architecture tests then (rightly) complain that it uses `System.IO`.
* Pass the settings file as an **absolute** path (`$PWD/...`): every test program starts in
  its own folder, and with a relative path `dotnet test` fails with an unhelpful
  `KeyNotFoundException`.
* The class filter leaves out the code that the `[GeneratedRegex]` source generator writes.
  Coverage is currently 100% of lines and 99% of branches (the branch left is the "look in
  the current folder" fallback when the program is not inside the repository).
  ReportGenerator's own `-minimumCoverageThresholds` is not in the version we pin, so the
  threshold lives in the small script [`check-coverage.sh`](check-coverage.sh).

`--report-xunit-junit` and `--report-xunit-trx` write one JUnit XML and one TRX file per test
project next to the coverage files; CI uploads them all.

## Folder layout

```
csharp/
├── Checkout.slnx                 ← the solution: lists every project below
├── global.json                   ← pins the .NET SDK and turns on Microsoft.Testing.Platform
├── Directory.Build.props         ← settings for every project (nullable, analyzers, warnings = errors, lock files)
├── .editorconfig                 ← code style, checked by `dotnet format` and the build
├── .config/dotnet-tools.json     ← pinned local tools: ReportGenerator, Stryker.NET
├── codecoverage.settings.xml     ← what code coverage measures
├── check-coverage.sh             ← fails the build below the coverage thresholds
├── stryker-config.json           ← Stryker.NET mutation testing
├── src/
│   ├── Checkout/                 ← the checkout itself: pure calculation, no files, no packages
│   │   ├── Checkout.cs           ← Checkout.Calculate(basket, catalog): the pipeline from rules §4; PlaceOrder
│   │   ├── Basket.cs, Catalog.cs, Product.cs, Promotion.cs, Receipt.cs
│   │   ├── ErrorCode.cs, CheckoutException.cs
│   │   ├── BasketValidator.cs    ← rows 1-4 of rules §4.1
│   │   ├── ReceiptLines.cs       ← merging basket items into receipt lines, pricing a line
│   │   ├── Promotions.cs         ← buy X get Y free, bulk discount
│   │   ├── Coupons.cs            ← SAVE10, WELCOME5
│   │   ├── Shipping.cs           ← standard / express
│   │   ├── Tax.cs                ← 8%
│   │   ├── Money.cs              ← RoundHalfUp, with whole numbers only
│   │   ├── Orders.cs             ← Order, IOrderRepository, InMemoryOrderRepository, order ids, timestamps
│   │   └── CheckoutJson.cs       ← the snake_case JSON settings everybody shares
│   ├── Checkout.Storage/         ← the order stores that touch the disk
│   │   ├── JsonFileOrderRepository.cs  ← one ORDER-<n>.json file per order
│   │   └── SqliteOrderRepository.cs    ← one table in an SQLite file (Microsoft.Data.Sqlite)
│   └── Checkout.Cli/             ← the command line program (rules §8)
│       ├── Program.cs            ← stdin → CheckoutCommand → stdout, exit code
│       ├── CheckoutCommand.cs    ← the CLI's logic: JSON text in, output + exit code out
│       └── SpecFolder.cs         ← finds ../spec/ from anywhere inside the repository
├── tests/                        ← one project per kind of test, see tests/README.md
│   ├── Directory.Build.props     ← what every test project shares (xUnit, coverage, TestSupport)
│   ├── xunit.runner.json         ← parallel settings
│   ├── Checkout.TestSupport/     ← helpers, not tests: spec loader, schema validator, fixtures, clocks, "$12.39" → 1239
│   ├── Checkout.UnitTests/       ← one small piece of code at a time, with made-up data; the 60 spec cases
│   ├── Checkout.IntegrationTests/  ← real files in a temp folder, a real SQLite file, the real spec files
│   ├── Checkout.Specs/           ← Reqnroll step definitions for spec/features/*.feature
│   ├── Checkout.PropertyTests/   ← P1 to P14 from spec/properties.md with FsCheck
│   ├── Checkout.ContractTests/   ← receipts and CLI errors checked against spec/schema/
│   ├── Checkout.SnapshotTests/   ← whole receipts compared with Snapshots/*.verified.txt
│   ├── Checkout.ArchitectureTests/ ← "the checkout does not depend on System.IO or the CLI"
│   └── Checkout.EndToEndTests/   ← the real program, run as a separate process
└── benchmarks/
    └── Checkout.Benchmarks/      ← a BenchmarkDotNet micro-benchmark of Checkout.Calculate
```

The shared data is read straight from the repository's [`spec/`](../spec/) folder (never
copied) – the feature files too, which `Checkout.Specs.csproj` includes with a relative path –
so every language is tested against exactly the same cases.

## The kinds of tests, and where to look

| Kind of test | Question it answers | Look at |
|---|---|---|
| **Unit** | Does this one method do the right thing? | `tests/Checkout.UnitTests/CatalogTests.cs` is the first read, then `MoneyTests.cs` |
| **Parameterised** | Does it work for all these examples? | `[Theory]` + `[InlineData]` in `PromotionsTests.cs`; `[MemberData]` in `CheckoutTests.cs` and `SpecCaseTests.cs` |
| **Shared spec cases** | Does it match the 60 cases every language must pass? | `tests/Checkout.UnitTests/SpecCaseTests.cs` |
| **BDD** | Does it behave as the Given/When/Then scenarios say? | `tests/Checkout.Specs/StepDefinitions.cs` – each sentence from `spec/features/README.md` appears once |
| **Test doubles** | How do I test without the real clock? Was `Save` really never called? | `tests/Checkout.UnitTests/OrdersTests.cs`: a fake (`InMemoryOrderRepository`), a stub clock (`FakeTimeProvider`), NSubstitute mocks |
| **Integration** | Do the repositories really write files and database rows? | `tests/Checkout.IntegrationTests/JsonFileOrderRepositoryTests.cs`, `SqliteOrderRepositoryTests.cs` |
| **Property-based** | Is it true for *every* basket, not just my examples? | `tests/Checkout.PropertyTests/CheckoutProperties.cs` – P1…P14, generators in `Generators.cs` |
| **Contract** | Is the output even the right shape? | `tests/Checkout.ContractTests/ReceiptSchemaTests.cs` |
| **Snapshot** | Did the output change since last time? | `tests/Checkout.SnapshotTests/ReceiptSnapshotTests.cs` and its `Snapshots/` |
| **End-to-end** | Does the real program work from the outside? | `tests/Checkout.EndToEndTests/CliTests.cs`, then the shared runner `e2e/run.sh` |
| **Architecture** | Does the code keep its shape (no files in the checkout)? | `tests/Checkout.ArchitectureTests/DependencyRulesTests.cs` |
| **Benchmark** | How fast is it? | `benchmarks/Checkout.Benchmarks/CheckoutBenchmarks.cs` |
| **Coverage** | Which lines did the tests never run? | [Coverage](#coverage) – 100% of lines, 99% of branches |
| **Mutation** | Would the tests notice if the code were wrong? | `dotnet stryker` – see below |
| **Static checks** | Can we catch mistakes without running anything? | `dotnet format --verify-no-changes`; `dotnet build` with `latest-recommended` analyzers |

### Tests that run in parallel, and tests that share

xUnit runs **test classes in parallel** with each other (each class is a "collection" unless
you say otherwise) and the tests **inside** one class one after another; `dotnet test` also runs
the test projects side by side. That is why no test depends on another: every test class gets a
fresh instance per test, integration tests each get their own temporary folder, and nothing
uses static mutable state. `tests/xunit.runner.json` spells those settings out.

When tests *should* share something expensive, put them in the same `[Collection]`:
`SpecCaseTests` is in the "standard catalog" collection, whose `StandardCatalogFixture` reads
the spec files once for all 60 cases (see `SpecCaseTests.cs` and
`Checkout.TestSupport/StandardCatalogFixture.cs`). Tests in one collection never run at the
same time.

### Mutation testing, and why the code looks the way it does

Stryker.NET makes ~120 small changes to `src/Checkout` ("mutants": `>=` becomes `>`, `+`
becomes `-`, a string becomes `""`) and runs the tests against each one. A mutant that survives
means no test noticed the change. The score is currently **100%** (94 mutants tested, all
killed; 10 more do not compile and 18 are skipped by Stryker because their code is already
covered by another mutant); CI fails below 95%. A run takes about five minutes.

Getting there shaped both the tests and the code:

* Boundaries are tested on both sides (`2499` / `2500`, `4999` / `5000`, `9` / `10`,
  `99` / `100`, `0` / `1`), because `>` vs `>=` is the classic surviving mutant.
* Error **messages** are tested, not only error codes – the first run showed that the messages
  in `ErrorCode.cs` could have been anything.
* Coupons were first a static dictionary. A mutant in a static field initializer is created
  only once per test process, before Stryker can switch it on, so with only the unit tests
  those mutants "survived". A `switch` expression reads at least as well and has no such
  problem.
* Stryker runs on the whole solution (`stryker-config.json` names `Checkout.slnx`), so all
  test projects take part. Its Microsoft.Testing.Platform runner is still marked *preview*;
  it needs each test project to reference `src/Checkout` directly, which
  `tests/Directory.Build.props` does.
* The thresholds in `stryker-config.json` must satisfy `high ≥ low ≥ break`. Otherwise
  Stryker.NET 5.0 exits with code 1 **without printing anything** (add `--verbosity trace` to
  see that it stopped).
* Only `src/Checkout` is mutated. The stores and the CLI are covered by integration and
  end-to-end tests, which run real files and processes and would make every mutant slow.

## How to add a new test

Say you want to check that `Money.RoundHalfUp` rounds 10% of $0.05 (0.5 cents) up to 1 cent.

1. **Pick the kind of test.** It is one method with made-up numbers: a unit test, so open
   `tests/Checkout.UnitTests/MoneyTests.cs` (or create a new class named `<Thing>Tests`).
2. **Write the test** with `[Fact]`, named `Method_Scenario_ExpectedResult`, following
   *arrange / act / assert*:

   ```csharp
   [Fact]
   public void RoundHalfUp_HalfACent_RoundsUpToOneCent()
   {
       // Arrange
       var amount = 5;

       // Act
       var rounded = Money.RoundHalfUp(amount, percent: 10);

       // Assert
       Assert.Equal(1, rounded);
   }
   ```

3. **Run it:** `dotnet test --project tests/Checkout.UnitTests --filter-method "*HalfACent*"`.
4. **Make sure it can fail.** Change `1` to `0` and check that the test goes red.
   A test that never fails does not test anything! (Mutation testing automates this idea.)
5. Same check, many inputs? Use `[Theory]` with `[InlineData(...)]` rows, or `[MemberData]`
   for rows built in code – see `PromotionsTests.cs` and `CheckoutTests.cs`.
6. **Run everything** before you open a pull request: `dotnet format`, `dotnet test`, and
   `dotnet stryker` if you changed `src/Checkout`.

Want to add a case for *every* language instead? Add it to
[`spec/test-cases/`](../spec/test-cases/) and it is picked up automatically by the spec cases,
contract and end-to-end tests and the shared runner (and update the count in
`tests/Checkout.IntegrationTests/SpecLoaderTests.cs`).

A snapshot test fails the first time it runs, because nothing is approved yet: Verify writes
`Snapshots/<Class>.<Test>.received.txt`. Read it; if it is right, rename it to
`.verified.txt` and commit it.

## Things to know about this implementation

* **Whole numbers only.** `Money.RoundHalfUp(amount, percent)` is `(amount × percent + 50) / 100`
  in integer maths (with a `long` in between), as the hint in rules §7 suggests. No `double` or
  `decimal` anywhere near money.
* **Parse, don't validate.** `Shipping.Parse` turns the basket's free-form `shipping` text into
  the enum `ShippingMethod` or throws; after that the compiler knows the value is fine.
* **A basket quantity is a `long`.** The basket schema allows any JSON integer, so
  `10000000000` must give `INVALID_QUANTITY` like `100` does, not "not a basket". After
  validation every quantity is from 1 to 99, so `ReceiptLine.Quantity` is an `int` again.
* **Records compare by value – except lists.** `Receipt` writes its own `Equals` so that two
  receipts with equal lines are equal, which the tests (and property P5) rely on.
* **The clock is .NET's `TimeProvider`.** Production code uses `TimeProvider.System`; tests pass
  a `FakeTimeProvider` from `Microsoft.Extensions.TimeProvider.Testing`. No clock interface of
  our own was needed.
* **Repositories are an interface** (`IOrderRepository` with `Save` and `Find`), with three
  implementations. The in-memory one doubles as the test fake. Ids are `ORDER-<n>`; the
  JSON-file store takes the next number after the highest existing file, the SQLite store uses
  the row id. `Find` returns `null` for "not found".
* **The domain has no packages and no files.** Files and SQLite live in `Checkout.Storage`,
  the terminal in `Checkout.Cli`; the architecture tests keep it that way.
* **The CLI** reads stdin, finds `spec/catalog.json` by walking up from the program's folder
  (then from the current folder), prints one line of JSON, and exits with 0 (receipt), 1
  (`{"error":"<CODE>"}`) or 2 (not a basket at all – a message on stderr, nothing on stdout).

## Ideas for the next contribution

Tier 1, 2 and 3 of the matrix are complete for C#, so what is left is polish:

* Compare with the TypeScript version: same spec, same kinds of tests, different idioms.
  That comparison is the point of the project.
* [AwesomeAssertions](https://awesomeassertions.org/) or [Shouldly](https://docs.shouldly.org/)
  for more readable failure messages – a nice before/after on one test class.
* [CliWrap](https://github.com/Tyrrzz/CliWrap) instead of `System.Diagnostics.Process` in the
  end-to-end tests, to show what a library adds.
* ReportGenerator's built-in coverage thresholds (`-minimumCoverageThresholds`) once the pinned
  version has them, replacing `check-coverage.sh`.
