---
layout: page
title: Getting started
---

## 1. Get the code

```bash
git clone {{ site.repository_url }}.git
cd the-most-well-tested-checkout-system
```

## 2. Look around

```
spec/        the rules and shared test data – the same for every language
python/      ┐
typescript/  │
go/          │ one folder per language
java/        │
csharp/      │
rust/        ┘
docs/        this website
```

Start with the rules: [`spec/checkout-rules.md`]({{ site.repository_files }}/spec/checkout-rules.md).

## 3. Pick a language and run its tests

You only need the tools for the language you pick. Every language folder has a
`README.md` with full instructions; here is the short version:

| Language   | You need          | Run the tests (from the language folder)                         |
|------------|-------------------|------------------------------------------------------------------|
| Python     | Python 3.10+      | `pip install -e ".[dev]"` then `pytest`                          |
| TypeScript | Node.js 22+       | `npm install` then `npm test`                                    |
| Go         | Go 1.24+          | `go test ./...`                                                  |
| Java       | JDK 21 + Maven    | `mvn verify`                                                     |
| C#         | .NET 10 SDK       | `dotnet test`                                                    |
| Rust       | Rust (via rustup) | `cargo test`                                                     |

You should see some tests **pass** and one group of tests **skipped**. The skipped
tests are the shared spec cases, waiting for someone to implement the checkout.

## 4. Make a test fail on purpose

The quickest way to understand a test is to break it. Open one of the example unit
tests in your language, change an expected value, and run the tests again. Read the
failure message: a good test tells you *what* was expected and *what* actually happened.
Then put it back.

## 5. Next steps

* Read about the [kinds of tests](kinds-of-tests.html) this project wants to show.
* Pick something from your language's *"Testing ideas we'd love contributions for"* list.
