# Test projects

In .NET, each **kind** of test usually gets its **own test project** (its own
`.csproj`), named `<Project>.<Kind>Tests`. Each one is added to `Checkout.slnx`,
so `dotnet test` in the `csharp/` folder runs them all.

| Folder                          | Kind of test | Status |
|---------------------------------|--------------|--------|
| `Checkout.UnitTests/`           | Unit tests: small, fast, no files/network. Also the shared data-driven spec cases. | exists |
| `Checkout.IntegrationTests/`    | Integration tests: real files or a real database, e.g. a JSON-file order repository (rules §6). | idea – add me! |
| `Checkout.Specs/` (Reqnroll)    | BDD: runs the Gherkin scenarios in `../../spec/features/`. | idea – add me! |
| `Checkout.EndToEndTests/`       | End-to-end: start the command line program (rules §8) and talk to it through stdin/stdout. | idea – add me! |

Why separate projects? Each kind needs different packages (Reqnroll, Testcontainers, …),
and you can run one kind on its own:

```bash
dotnet test --project tests/Checkout.UnitTests
```

## Adding a new test project

```bash
cd csharp
dotnet new xunit3 -n Checkout.IntegrationTests -o tests/Checkout.IntegrationTests
dotnet sln add tests/Checkout.IntegrationTests
dotnet add tests/Checkout.IntegrationTests reference src/Checkout
```

Then delete the lines from the new `.csproj` that `../Directory.Build.props` already
sets (`TargetFramework`, `Nullable`, `ImplicitUsings`), and delete the `global.json`
that the template may create next to the test project or at the repository root
(ours lives in `csharp/global.json`).

If the new project needs the shared spec too, move `Checkout.UnitTests/SpecSupport/`
into a small class library (e.g. `tests/Checkout.TestSupport/`) that both projects reference.
