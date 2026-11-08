# Checkout – C# (.NET)

The C# version of the checkout system described in [`../spec/checkout-rules.md`](../spec/checkout-rules.md).

> **Status: tests first, code later.** The data types and the test setup are ready,
> but `Checkout.Calculate` is not implemented yet – it throws `NotImplementedException`.
> The 60 shared spec cases are written and **skipped**, waiting for someone to
> implement it (that is test-driven development!).

## Tools we use, and why

| Tool | What it does | Why |
|------|--------------|-----|
| [.NET 10 SDK](https://dotnet.microsoft.com/download) | The C# compiler and the `dotnet` command | Builds, formats and tests everything; the version is pinned in [`global.json`](global.json) |
| [xUnit.net v3](https://xunit.net/) | The test framework: `[Fact]`, `[Theory]`, `Assert.Equal(...)` | The most widely used test framework for modern .NET |
| [Microsoft.Testing.Platform](https://learn.microsoft.com/dotnet/core/testing/microsoft-testing-platform-intro) | The engine `dotnet test` uses to run the tests | The new, faster .NET test runner; switched on in `global.json` |
| `System.Text.Json` | Reads the shared spec JSON files | Built into .NET – nothing extra to install |
| `dotnet format` | Static check: code style and formatting (rules in [`.editorconfig`](.editorconfig)) | Built into the SDK; CI fails when code is not formatted |

Only **one** NuGet package is used: `xunit.v3.mtp-v2` (xUnit v3 for Microsoft.Testing.Platform).
Its exact version is in the `.csproj`, and `packages.lock.json` records every package it
pulls in, so everybody gets the same versions. Extra libraries (BDD, assertions,
property-based testing, …) are welcome as [contributions](#testing-ideas-wed-love-contributions-for).

## Prerequisites and install

1. Install the **.NET 10 SDK** from <https://dotnet.microsoft.com/download> and check it works:

   ```bash
   dotnet --version    # 10.0.100 or newer
   ```

2. From the repository root:

   ```bash
   cd csharp
   dotnet restore      # downloads the NuGet packages
   dotnet test         # builds and runs all the tests
   ```

Any editor works. Visual Studio, JetBrains Rider and VS Code with the
*C# Dev Kit* extension can also run and debug single tests with one click.

## Running tests

All commands are run from this `csharp/` folder. `dotnet test` finds `Checkout.slnx`
(the solution file that lists every project) and runs every test project in it.

| I want to… | Command |
|------------|---------|
| run all tests | `dotnet test` |
| see every test name (passed / skipped) | `dotnet test --output detailed` |
| run only the unit tests (one test project) | `dotnet test --project tests/Checkout.UnitTests` |
| run one file (= one test class) | `dotnet test --filter-class "*.CatalogTests"` |
| run one test method | `dotnet test --filter-method "*.FindBySku_KnownSku_ReturnsTheProduct"` |
| run one spec case | `dotnet test --filter "DisplayName~promo-coffee-three-bags-one-free"` |
| check formatting (as CI does) | `dotnet format --verify-no-changes` |
| fix formatting | `dotnet format` |

### How filtering works

xUnit gives `dotnet test` a few ways to pick tests. Use **one kind** of filter per command.

* **Simple filters** (xUnit's own): `--filter-class`, `--filter-method`, `--filter-namespace`,
  `--filter-trait`. They take full names, and `*` is a wildcard at the start or end:
  `--filter-class "*SpecLoaderTests"`, `--filter-method "*FindBySku*"`.
  Repeat an option to allow several values. There are `--filter-not-…` versions too.
* **VSTest-style filter** (the classic .NET syntax you will find in many tutorials):
  `--filter "<property><operator><value>"`, where the property is `FullyQualifiedName`,
  `DisplayName` or a trait name such as `Category`, and the operator is
  `=` (equals), `!=`, `~` (contains) or `!~`. Combine with `&` (and) and `|` (or):

  ```bash
  dotnet test --filter "FullyQualifiedName~CatalogTests"
  dotnet test --filter "FullyQualifiedName~UnitTests&FullyQualifiedName!~SpecCaseTests"
  dotnet test --filter "DisplayName~validation-"         # every validation spec case
  ```

**One file?** .NET runs tests by *name*, not by file. By convention one file holds one
test class with the same name (`CatalogTests.cs` → `class CatalogTests`), so
"run one file" means "run one class".

## Folder layout

```
csharp/
├── Checkout.slnx                 ← the solution: lists every project below
├── global.json                   ← pins the .NET SDK and turns on Microsoft.Testing.Platform
├── Directory.Build.props         ← settings shared by every project (nullable, warnings = errors, …)
├── .editorconfig                 ← code style, checked by `dotnet format`
├── src/
│   └── Checkout/                 ← the checkout itself (a class library)
│       ├── Checkout.cs           ← Checkout.Calculate(basket, catalog) – not implemented yet
│       ├── Basket.cs, Catalog.cs, Product.cs, Promotion.cs, Receipt.cs
│       └── ErrorCode.cs, CheckoutException.cs
└── tests/
    ├── README.md                 ← where integration, BDD and end-to-end test projects go
    └── Checkout.UnitTests/       ← unit tests (xUnit)
        ├── CatalogTests.cs       ← START HERE: an example unit test, explained line by line
        ├── ErrorCodeTests.cs     ← an example [Theory] (one test, many inputs)
        ├── SpecLoaderTests.cs    ← checks that the shared spec files load correctly
        ├── SpecCaseTests.cs      ← one test per case in spec/test-cases/*.json (skipped for now)
        └── SpecSupport/          ← test helpers: reads ../spec/ (the files are never copied)
```

### Where do other kinds of tests go?

The C# convention is **one test project per kind of test**, all in `tests/` and all
listed in `Checkout.slnx`:

* `tests/Checkout.UnitTests/` – unit tests (exists)
* `tests/Checkout.IntegrationTests/` – tests with real files or a real database
* `tests/Checkout.Specs/` – BDD: the Gherkin files in [`../spec/features/`](../spec/features/), run with [Reqnroll](https://reqnroll.net/)
* `tests/Checkout.EndToEndTests/` – start the command line program (rules §8) and check its output

This is different from, say, Go (tests next to the code) or Python (one `tests/`
folder with sub-folders). Separate projects keep each kind's packages apart and make
it easy to run one kind: `dotnet test --project tests/Checkout.IntegrationTests`.
See [`tests/README.md`](tests/README.md) for how to add a project.

## The code in 30 seconds

```csharp
var catalog = new Catalog([
    new Product("APPLE", "Apple", UnitPrice: 50),
    new Product("COFFEE", "Bag of Coffee Beans", UnitPrice: 899, new BuyXGetYFree(Buy: 2, Free: 1)),
]);
var basket = new Basket([new BasketItem("COFFEE", 3)], Coupon: "SAVE10");

try
{
    Receipt receipt = Checkout.Calculate(basket, catalog);   // throws NotImplementedException for now
    Console.WriteLine(receipt.Total);                        // money is always whole cents (int)
}
catch (CheckoutException e)
{
    Console.WriteLine(e.Code.ToSpecCode());                  // e.g. "UNKNOWN_COUPON"
}
```

* All domain types are **records**: small immutable classes that compare by value,
  so `Assert.Equal(new BuyXGetYFree(2, 1), product.Promotion)` just works.
* `Promotion` is an abstract record with two kinds, `BuyXGetYFree` and `BulkDiscount`.
* The namespace is `CheckoutSystem` (not `Checkout`), so that it does not clash with the
  `Checkout` class.

## How to add a new test

1. Pick (or create) a test class in `tests/Checkout.UnitTests/`. One class per thing
   you test, named `<Thing>Tests`, e.g. `RoundingTests.cs`.
2. Write a public method with `[Fact]`, named `Method_Scenario_ExpectedResult`:

   ```csharp
   namespace CheckoutSystem.UnitTests;

   public class RoundingTests
   {
       [Fact]
       public void RoundHalfUp_ExactHalf_RoundsUp()
       {
           // Arrange
           var amount = 225;   // 225 × 10% = 22.5 cents

           // Act
           var discount = Money.PercentOf(amount, percent: 10);

           // Assert
           Assert.Equal(23, discount);
       }
   }
   ```

3. Run just that class while you work: `dotnet test --filter-class "*.RoundingTests"`.
   Watch it fail first (red), write the code (green), then tidy up (refactor).
4. Same check, many inputs? Use `[Theory]` with `[InlineData(...)]` rows –
   see `ErrorCodeTests.cs`.
5. Before you push: `dotnet format` and `dotnet test`.

**Adding a shared spec case?** Add it to `../spec/test-cases/*.json` (see
[the README there](../spec/test-cases/README.md)). `SpecCaseTests` picks it up
automatically; only the case count in `SpecLoaderTests` needs updating.

**Starting the implementation?** Remove `Skip = "..."` from the `[MemberData]` line in
`SpecCaseTests.cs`. All 60 cases will fail – now make them pass, one rule at a time.

## Testing ideas we'd love contributions for

This list mirrors the C# column of the [testing matrix](../docs/testing-matrix.md) – the
project's single list of goals – and adds the practical hints. Tick the matrix when you add
one. We are on **xUnit v3** with **Microsoft.Testing.Platform v2**; pin companion packages
(coverage, Verify) to versions that say they support it.

| Tier | Idea | Tools to try |
|------|------|--------------|
| 1 | Implement `Checkout.Calculate` test-first, with unit tests for each step | xUnit `[Fact]` |
| 1 | Parameterised tests – the native table style | `[Theory]` + `[InlineData]` / `[MemberData]` |
| 1 | BDD: run the Gherkin scenarios in `../spec/features/` | [Reqnroll](https://reqnroll.net/) (`Reqnroll.xunit.v3`) in a new `tests/Checkout.Specs` project |
| 1 | Test doubles for the order repository and the clock (rules §6) | A hand-written `InMemoryOrderRepository` and a fixed `TimeProvider` first; then [NSubstitute](https://nsubstitute.github.io/) to check a failed checkout never calls `Save` |
| 1 | Integration tests: a JSON-file and an SQLite repository against a temp folder; stored orders match [`order.schema.json`](../spec/schema/order.schema.json) | `System.IO` + `Path.GetTempPath()`, [Microsoft.Data.Sqlite](https://learn.microsoft.com/dotnet/standard/data/sqlite/) |
| 1 | Property-based tests – one per property in [`spec/properties.md`](../spec/properties.md) | [FsCheck](https://fscheck.github.io/FsCheck/) (`FsCheck.Xunit.v3`), or [CsCheck](https://github.com/AnthonyLloyd/CsCheck) |
| 1 | Contract tests: every receipt satisfies [`receipt.schema.json`](../spec/schema/receipt.schema.json) | [JsonSchema.Net](https://docs.json-everything.net/schema/basics/) |
| 1 | Snapshot / approval tests of whole receipts | [Verify](https://github.com/VerifyTests/Verify) (`Verify.XunitV3`) |
| 1 | Nicer assertions | [AwesomeAssertions](https://awesomeassertions.org/) (the open-source fork of FluentAssertions, which went commercial in v8) or [Shouldly](https://docs.shouldly.org/) |
| 1 | End-to-end tests of the CLI (rules §8), then the shared runner `e2e/run.sh dotnet run --project csharp/src/Checkout.Cli` | A `Checkout.Cli` console project, tested with `System.Diagnostics.Process` or [CliWrap](https://github.com/Tyrrzz/CliWrap) |
| 1 | Code coverage | [`Microsoft.Testing.Extensions.CodeCoverage`](https://learn.microsoft.com/dotnet/core/testing/microsoft-testing-platform-extensions-code-coverage) (`dotnet test --coverage`) + [ReportGenerator](https://reportgenerator.io/) |
| 1 | Mutation testing | [Stryker.NET](https://stryker-mutator.io/docs/stryker-net/introduction/) with its Microsoft.Testing.Platform runner (still marked preview) |
| 1 | Stricter static analysis | `<AnalysisLevel>latest-recommended</AnalysisLevel>` and `<EnforceCodeStyleInBuild>true</EnforceCodeStyleInBuild>` in `Directory.Build.props` |
| 2 | Architecture tests ("the domain does not depend on `System.IO`") | [ArchUnitNET](https://github.com/TNG/ArchUnitNET) |
| 2 | One micro-benchmark of `Calculate` | [BenchmarkDotNet](https://benchmarkdotnet.org/) |
| 3 | Test report for CI; filtering; parallel collections | `dotnet test --report-junit` (or `--report-xunit-trx`), `[Trait]` + `--filter`, `[Collection]` |

We deliberately skip Moq (see its 2023 SponsorLink episode) in favour of NSubstitute, and
SharpFuzz (needs external fuzzer binaries) in favour of FsCheck for random input.
