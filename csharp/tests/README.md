# Test projects

In .NET, each **kind** of test usually gets its **own test project** (its own
`.csproj`), named `Checkout.<Kind>Tests`. Each one is added to `../Checkout.slnx`,
so `dotnet test` in the `csharp/` folder runs them all.

| Folder                          | Kind of test | Tests |
|---------------------------------|--------------|------:|
| `Checkout.UnitTests/`           | Unit tests: small, fast, no files. Also the 60 shared data-driven spec cases and the test doubles (fake, stub clock, NSubstitute mocks). | 309 |
| `Checkout.IntegrationTests/`    | Integration tests: the JSON-file and SQLite order repositories against a temporary folder, and the spec loader against the real spec files. | 41 |
| `Checkout.Specs/` (Reqnroll)    | BDD: runs the Gherkin scenarios in `../../spec/features/` (read from there, never copied). | 44 |
| `Checkout.PropertyTests/`       | Property-based tests P1–P14 from `../../spec/properties.md`, with FsCheck. | 14 |
| `Checkout.ContractTests/`       | Receipts, saved orders and CLI errors checked against `../../spec/schema/` with JsonSchema.Net. | 67 |
| `Checkout.SnapshotTests/`       | Whole receipts compared with the approved `Snapshots/*.verified.txt` (Verify). | 5 |
| `Checkout.ArchitectureTests/`   | "The checkout does not depend on `System.IO`, the console, the CLI or the stores" (ArchUnitNET). | 5 |
| `Checkout.EndToEndTests/`       | End-to-end: starts the real command line program (rules §8) as a separate process for every spec case, plus input that is not a basket. | 67 |
| `Checkout.TestSupport/`         | Not tests: helpers every test project uses (spec loader, schema validator, fixtures, clocks, dollars → cents). | – |

Why separate projects? Each kind needs different packages (Reqnroll, FsCheck, Verify, …),
and you can run one kind on its own:

```bash
dotnet test --project tests/Checkout.IntegrationTests
```

[`Directory.Build.props`](Directory.Build.props) in this folder gives every test project the
same basics, so each `.csproj` only lists what makes its kind of test special: xUnit v3 on
Microsoft.Testing.Platform, the code coverage extension, a reference to `src/Checkout` and to
`Checkout.TestSupport`, and the shared [`xunit.runner.json`](xunit.runner.json).

## How the end-to-end tests find the program

`Checkout.EndToEndTests` does not reference `src/Checkout.Cli` itself, but
`Checkout.TestSupport` does, so building the tests always builds the program first. The test
project writes the program's path (`src/Checkout.Cli/bin/<Configuration>/net10.0/Checkout.Cli.dll`)
into its own assembly as `[AssemblyMetadata]`, and the tests start it with `dotnet <that path>`.
`CliTests.TheProgram_HasBeenBuilt` fails with a clear message if it is missing.

## Adding a new test project

```bash
cd csharp
dotnet new xunit3 -n Checkout.SomethingTests -o tests/Checkout.SomethingTests
dotnet sln add tests/Checkout.SomethingTests
```

Then delete from the new `.csproj` everything that `../Directory.Build.props` and
`tests/Directory.Build.props` already set (`TargetFramework`, `Nullable`, `ImplicitUsings`,
`OutputType`, the xUnit package and the project reference), delete the `global.json` that the
template may create (ours lives in `csharp/global.json`), and run `dotnet restore` so that the
new `packages.lock.json` is written. Commit it.
