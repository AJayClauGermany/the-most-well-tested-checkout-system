# Checkout system – Rust

The Rust version of the checkout system. The rules are in
[`../spec/checkout-rules.md`](../spec/checkout-rules.md). Read them first.

> **Status: tests first, code later.** The types and the shared tests are ready, but
> `checkout()` only contains `todo!()`. The shared spec-case test is *ignored* until
> someone implements it. That is test-driven development (TDD): the tests are written
> first and wait for the code that makes them pass.

## Tools we use (and why)

| Tool | What it does | Why |
|------|--------------|-----|
| [Rust](https://www.rust-lang.org/) + **Cargo** | Compiler, build tool and package manager in one | Cargo is how every Rust project is built and tested |
| `cargo test` (built in) | Finds and runs every `#[test]` function | No extra test framework needed: testing is part of the language |
| [`serde`](https://serde.rs/) + [`serde_json`](https://docs.rs/serde_json) | Turns JSON into Rust structs | The shared spec files are JSON |
| `rustfmt` (`cargo fmt`) | Formats the code in the one standard Rust style | Nobody argues about spaces and braces |
| `clippy` (`cargo clippy`) | A linter: spots common mistakes and un-idiomatic code | Catches bugs before the tests even run |

The exact dependency versions are pinned in `Cargo.lock`, which is committed so that
everyone (and CI) builds with exactly the same versions.

## Prerequisites

Install Rust with [rustup](https://rustup.rs/) (Rust 1.85 or newer, because we use the
2024 edition), then add the formatter and linter:

```bash
rustup component add rustfmt clippy
```

Cargo downloads the dependencies by itself the first time you build or test.
All commands below are run from this `rust/` folder.

## Running the tests

| What | Command |
|------|---------|
| All tests | `cargo test` |
| Only unit tests (the ones inside `src/`) | `cargo test --lib` |
| Only one integration test file | `cargo test --test spec_files` |
| Only tests whose name contains some text | `cargo test find_is_case_sensitive` |
| Only the *ignored* tests | `cargo test -- --ignored` |
| Everything, *including* ignored tests | `cargo test -- --include-ignored` |
| Show `println!` output of passing tests | `cargo test -- --nocapture` |

Notice the `--` in some commands: options **before** it go to Cargo (which tests to
build), options **after** it go to the test runner (which tests to run and how).

The static checks that CI also runs:

```bash
cargo fmt --check                           # is the code formatted? (`cargo fmt` fixes it)
cargo clippy --all-targets -- -D warnings   # any lint warnings? (-D turns them into errors)
```

Today `cargo test` shows 8 unit tests, 5 spec-loader tests and 1 doc test passing,
plus 1 ignored test (the shared spec cases).

## Folder layout

Rust has its own conventions for where tests go, and they are different from most
other languages:

```
rust/
├── Cargo.toml              ← project file: name, edition, dependencies
├── Cargo.lock              ← exact dependency versions (committed)
├── src/
│   ├── lib.rs              ← the public API, including the checkout() stub
│   ├── catalog.rs          ← Product, Promotion, Catalog      + their unit tests
│   ├── basket.rs           ← BasketItem, Basket               + their unit tests
│   ├── receipt.rs          ← Receipt, ReceiptLine
│   └── error.rs            ← ErrorCode, CheckoutError         + their unit tests
└── tests/                  ← integration tests: one program per .rs file
    ├── common/mod.rs       ← shared helpers: loads ../spec/ (not a test program itself)
    ├── spec_files.rs       ← checks that the spec files load correctly
    └── spec_cases.rs       ← runs every case in ../spec/test-cases/ (ignored for now)
```

* **Unit tests live next to the code**, at the bottom of each `src/*.rs` file inside a
  `#[cfg(test)] mod tests { ... }` block. `#[cfg(test)]` means the block is only compiled
  for `cargo test`, so it never ends up in the real program. Because the tests sit inside
  the module, they can also test private functions.
* **Integration tests live in `tests/`.** Cargo compiles each file there as a separate
  program that can only use the crate's *public* API – just like a real user of the crate.
  Code shared between test files goes in a sub-folder (`tests/common/mod.rs`), because
  Cargo would otherwise treat it as one more test program.
* **Doc tests** are the code examples in `///` comments. `cargo test` compiles (and
  normally runs) them, so the documentation can never go out of date. See the example
  on `checkout()` in `src/lib.rs`; it is marked `no_run` until checkout is implemented.
* **BDD tests** (for `../spec/features/*.feature`) would go in `tests/features.rs`, run by
  the `cucumber` crate (see the ideas list below).
* **End-to-end tests** would test a small command line program (rules §8) that would live
  in `src/bin/checkout.rs`; its tests would go in `tests/cli.rs`.

The tests read the spec straight from the repository's `../spec/` folder (it is never
copied). `tests/common/mod.rs` finds it through `env!("CARGO_MANIFEST_DIR")`, the folder
containing `Cargo.toml`, so the tests work no matter which folder you start them from.

### Why is the shared spec-case test one test with a loop?

Rust's built-in test runner has no *parameterized tests* (one test generated per row of
data). So `tests/spec_cases.rs` has a single test that loops over all 60 cases, collects
every failure together with its case id, and fails once at the end with the full list.
That way a single run shows **all** broken cases, not just the first one.

It is marked `#[ignore = "..."]` because `checkout()` is not implemented yet. Ignored
tests are still compiled (so they cannot rot), but not run. Try
`cargo test --test spec_cases -- --ignored` to see it fail with "not implemented yet".
When you start implementing, delete the `#[ignore]` line.

## How to add a new test

Say you implement a small `round_half_up` helper in a new file `src/money.rs`
(and add `mod money;` to `src/lib.rs`). Add a unit test at the bottom of that file:

```rust
pub fn round_half_up(amount: i64, percent: i64) -> i64 {
    (amount * percent + 50) / 100
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn exact_halves_round_up() {
        // Arrange: 225 × 10% = 22.5 cents
        let amount = 225;

        // Act
        let rounded = round_half_up(amount, 10);

        // Assert
        assert_eq!(rounded, 23);
    }
}
```

1. Write the test **first** and run `cargo test exact_halves` – watch it fail (red).
2. Write just enough code to make it pass (green).
3. Tidy up the code while the test stays green (refactor), then run
   `cargo fmt` and `cargo clippy --all-targets -- -D warnings`.

Give tests names that read like a sentence about the behaviour
(`find_is_case_sensitive`, not `test3`): the name is what you see when it fails.

For a test that uses only the public API (or the spec files), create a new file in
`tests/`, for example `tests/receipts.rs`, and write `mod common;` at the top if you
need the spec loader.

## Testing ideas we'd love contributions for

* **BDD:** run `../spec/features/*.feature` with the [`cucumber`](https://crates.io/crates/cucumber) crate.
* **One test per spec case:** [`rstest`](https://crates.io/crates/rstest) (parameterized tests),
  [`libtest-mimic`](https://crates.io/crates/libtest-mimic) or
  [`datatest-stable`](https://crates.io/crates/datatest-stable), so every case id shows up
  as its own test.
* **Property-based testing:** [`proptest`](https://crates.io/crates/proptest) or
  [`quickcheck`](https://crates.io/crates/quickcheck) – e.g. "the total is never negative",
  "merging lines never changes the subtotal", "shuffling the basket gives the same totals".
* **Mocking / fakes:** an in-memory order repository (rules §6) behind a trait, with
  [`mockall`](https://crates.io/crates/mockall) to check how it is called.
* **Snapshot testing:** [`insta`](https://crates.io/crates/insta) for whole receipts.
* **Better assertion messages:** [`pretty_assertions`](https://crates.io/crates/pretty_assertions)
  shows a coloured diff when `assert_eq!` fails.
* **Coverage:** [`cargo-llvm-cov`](https://github.com/taiki-e/cargo-llvm-cov) or
  [`cargo-tarpaulin`](https://github.com/xd009642/tarpaulin).
* **Mutation testing:** [`cargo-mutants`](https://mutants.rs/) – do our tests notice when
  the code is changed on purpose?
* **End-to-end:** the command line program from rules §8 in `src/bin/checkout.rs`, tested
  with [`assert_cmd`](https://crates.io/crates/assert_cmd) in `tests/cli.rs`.
* **Fuzzing:** [`cargo-fuzz`](https://github.com/rust-fuzz/cargo-fuzz) on the basket JSON parser.
* **Benchmarks:** [`criterion`](https://crates.io/crates/criterion) or
  [`divan`](https://crates.io/crates/divan).
* **A faster test runner:** [`cargo-nextest`](https://nexte.st/).
