# The Most Well-Tested Checkout System

> **Given a basket of products, calculate the final price using quantities, discounts,
> shipping, and tax, then optionally save the resulting order.**

A deliberately small checkout system, written in **six languages**, used to show and
compare **as many ways of testing as possible** – from unit tests all the way up to
end-to-end tests.

It is written for people who are **getting started with programming** (late high school,
early college, self-taught). The code favours being easy to read over being fast or
production-ready.

📖 **Documentation website:** <https://ajayclaugermany.github.io/the-most-well-tested-checkout-system/>
(built from [`docs/`](docs/) by GitHub Pages)

## Status: tests first

The rules and the shared tests exist, and the project is built test-first (TDD / BDD).
**TypeScript is the first language with a finished checkout**: every row of the testing
matrix is ticked for it, so it is the point of comparison for the other five. In Python,
Go, Java, C# and Rust the project is set up and the tests run in CI, but the 60 shared
test cases are still skipped, waiting to go from red to green.

## What's in here

| Folder                          | What it is                                                                     |
|---------------------------------|--------------------------------------------------------------------------------|
| [`spec/`](spec/)                | The **language-agnostic** rules, catalog, shared test cases, BDD scenarios, properties and schemas |
| [`e2e/`](e2e/)                  | One end-to-end runner that drives every language's command line program          |
| [`python/`](python/)            | Python – pytest                                                                |
| [`typescript/`](typescript/)    | TypeScript on Node.js – Vitest                                                  |
| [`go/`](go/)                    | Go – the built-in `testing` package                                            |
| [`java/`](java/)                | Java – JUnit with Maven                                                         |
| [`csharp/`](csharp/)            | C# on .NET – xUnit                                                              |
| [`rust/`](rust/)                | Rust – the built-in test harness with Cargo                                     |
| [`docs/`](docs/)                | The documentation website (Jekyll, published by GitHub Pages)                 |
| [`.github/workflows/`](.github/workflows/) | CI: one workflow per language, plus one for the spec              |

## Quick start

1. Read the rules: [`spec/checkout-rules.md`](spec/checkout-rules.md).
2. Pick a language and follow its `README.md`. In short (run inside the language folder):

   | Language   | Needs            | Run the tests                                   |
   |------------|------------------|-------------------------------------------------|
   | Python     | Python 3.10+     | `pip install -e ".[dev]"` then `pytest`         |
   | TypeScript | Node.js 22.18+   | `npm install` then `npm test`                   |
   | Go         | Go 1.24+         | `go test ./...`                                 |
   | Java       | JDK 21 + Maven   | `mvn verify`                                    |
   | C#         | .NET 10 SDK      | `dotnet test`                                   |
   | Rust       | Rust (rustup)    | `cargo test`                                    |

3. In TypeScript every test passes. In the other languages expect passing tests **and**
   skipped tests: the skipped ones are the shared spec cases, waiting for an implementation.

## How the shared spec works

The [`spec/`](spec/) folder describes the checkout **once**, without any code:

* [`checkout-rules.md`](spec/checkout-rules.md) – the rules in plain English
* [`catalog.json`](spec/catalog.json) – the six products every test uses
* [`test-cases/`](spec/test-cases/) – 60 data-driven cases: basket in → totals or error out
* [`features/`](spec/features/) – Given/When/Then scenarios for BDD tools (Cucumber & friends)
* [`properties.md`](spec/properties.md) – rules that hold for *every* basket, for property-based tests
* [`schema/`](spec/schema/) – JSON Schemas for baskets, receipts and saved orders, for contract tests

Every language reads these files directly, so all six implementations are checked
against exactly the same examples.

## What we want to show

The goal is written down once, as a table of *kind of test × language → tool*:
**[the testing matrix](docs/testing-matrix.md)**. Tier 1 is what every language does
(unit, data-driven, BDD, test doubles, integration, property-based, contract, snapshot,
end-to-end, coverage, mutation, static checks); Tier 2 is what only some languages do
idiomatically (doctests, fuzzing, benchmarks, …). It also lists what we leave out on purpose.

## Continuous integration

CI runs on **pull requests only** (or by hand from the *Actions* tab), and each
workflow only runs when its own language folder or the shared `spec/` changes.

## Contributing

New tests, new kinds of tests, new tools, better explanations – all welcome.
Start with [`CONTRIBUTING.md`](CONTRIBUTING.md).

## License

[Boost Software License 1.0](LICENSE)
