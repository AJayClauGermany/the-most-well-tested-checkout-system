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

The rules and the shared tests exist. **The checkout itself is not implemented yet.**
That is on purpose: the project starts test-first (TDD / BDD). Each language already has
its project set up, its tests running in CI, and a skipped suite of 60 shared test cases
waiting to go from red to green.

## What's in here

| Folder                          | What it is                                                                     |
|---------------------------------|--------------------------------------------------------------------------------|
| [`spec/`](spec/)                | The **language-agnostic** rules, product catalog, shared test cases and BDD scenarios |
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
   | TypeScript | Node.js 22+      | `npm install` then `npm test`                   |
   | Go         | Go 1.24+         | `go test ./...`                                 |
   | Java       | JDK 21 + Maven   | `mvn verify`                                    |
   | C#         | .NET 10 SDK      | `dotnet test`                                   |
   | Rust       | Rust (rustup)    | `cargo test`                                    |

3. Expect to see passing tests **and** skipped tests. The skipped ones are the shared
   spec cases, waiting for an implementation.

## How the shared spec works

The [`spec/`](spec/) folder describes the checkout **once**, without any code:

* [`checkout-rules.md`](spec/checkout-rules.md) – the rules in plain English
* [`catalog.json`](spec/catalog.json) – the six products every test uses
* [`test-cases/`](spec/test-cases/) – 60 data-driven cases: basket in → totals or error out
* [`features/`](spec/features/) – Given/When/Then scenarios for BDD tools (Cucumber & friends)

Every language reads these files directly, so all six implementations are checked
against exactly the same examples.

## Continuous integration

CI runs on **pull requests only** (or by hand from the *Actions* tab), and each
workflow only runs when its own language folder or the shared `spec/` changes.

## Contributing

New tests, new kinds of tests, new tools, better explanations – all welcome.
Start with [`CONTRIBUTING.md`](CONTRIBUTING.md).

## License

[Boost Software License 1.0](LICENSE)
