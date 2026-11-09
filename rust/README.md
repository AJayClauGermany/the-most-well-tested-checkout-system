# Checkout system – Rust

The Rust version of the checkout described in [`spec/checkout-rules.md`](../spec/checkout-rules.md).

> **Status: done.** `checkout()`, saving orders (in memory, as JSON files and in SQLite) and
> the command line program are implemented, and every row of the Rust column of the
> [testing matrix](../docs/testing-matrix.md) is covered. It was built strictly test-first:
> the commit that adds the tests (all red, every function a `todo!()`) comes before the
> commit that makes them pass – have a look at the git history.

## Tools we use (and why)

| Tool | What it does | Why this one |
|------|--------------|--------------|
| [Rust](https://www.rust-lang.org/) + **Cargo** | Compiler, build tool and package manager in one | Cargo is how every Rust project is built and tested |
| `cargo test` (built in) | Finds and runs every `#[test]` function, and the examples in `///` comments | Testing is part of the language: no framework needed for unit, integration and doc tests |
| [`serde`](https://serde.rs/) + [`serde_json`](https://docs.rs/serde_json) | Turns JSON into Rust structs and back | The shared spec files, the CLI and the saved orders are JSON |
| [`rusqlite`](https://crates.io/crates/rusqlite) (`bundled`) | SQLite for the SQLite order repository | The usual SQLite driver; `bundled` compiles SQLite in, so there is nothing to install |
| [`rstest`](https://crates.io/crates/rstest) | Parameterised tests (`#[case]`) | Rust's built-in harness has no table tests; rstest is the common way to write them |
| [`libtest-mimic`](https://crates.io/crates/libtest-mimic) | Builds a list of tests at run time | Gives each of the 60 shared spec cases (and each feature file) its own test name, read from the spec folder |
| [`cucumber`](https://crates.io/crates/cucumber) | Runs the Given/When/Then scenarios in `spec/features/` | The Rust implementation of Cucumber; `async`, so we drive it with the small executor from [`futures`](https://crates.io/crates/futures) |
| [`mockall`](https://crates.io/crates/mockall) | Generates mock objects from a trait (`#[automock]`) | The most used mocking library; checks *how* a dependency was called |
| [`tempfile`](https://crates.io/crates/tempfile) | Temporary folders that delete themselves | Integration tests with real files and databases, without leftovers |
| [`proptest`](https://crates.io/crates/proptest) | Property-based testing | Hypothesis-style generation and shrinking, the standard choice for Rust |
| [`jsonschema`](https://crates.io/crates/jsonschema) | JSON Schema validation (contract tests) | Supports draft 2020-12, which the spec schemas use |
| [`insta`](https://insta.rs/) | Snapshot testing | The snapshot library for Rust, with a review tool (`cargo insta review`) |
| [`pretty_assertions`](https://crates.io/crates/pretty_assertions) | A drop-in `assert_eq!` that shows a coloured diff | Long receipts are much easier to compare when the difference is highlighted |
| [`assert_cmd`](https://crates.io/crates/assert_cmd) + [`predicates`](https://crates.io/crates/predicates) | Run the real program and check its output | The usual way to test a Rust command line program from the outside |
| [`criterion`](https://crates.io/crates/criterion) | Micro-benchmarks | Statistically sound timings that work on stable Rust |
| `rustfmt` + `clippy` | Formatter and linter | The standard Rust tools; CI treats every warning as an error |
| [`cargo-llvm-cov`](https://github.com/taiki-e/cargo-llvm-cov) | Code coverage | Source-based coverage that works on stable Rust |
| [`cargo-mutants`](https://mutants.rs/) | Mutation testing | The mutation testing tool for Rust; needs no changes to the code |
| [`cargo-nextest`](https://nexte.st/) | A faster test runner | Runs every test in its own process, in parallel, and writes a JUnit XML report for CI |
| [`cargo-fuzz`](https://rust-fuzz.github.io/book/cargo-fuzz.html) | Fuzzing with libFuzzer | The standard Rust fuzzer (needs a nightly compiler, so it is optional) |

The exact versions of every library are pinned in `Cargo.lock`, which is committed, so
everyone (and CI) builds with exactly the same code. The `cargo install` tools are pinned in
the commands below and in the CI workflow.

## Getting started

Install Rust with [rustup](https://rustup.rs/). You need **Rust 1.89 or newer** (the 2024
edition needs 1.85, the SQLite driver needs 1.89); check with `rustc --version`. Then add
the formatter and linter:

```bash
rustup component add rustfmt clippy
cd rust
cargo test        # builds everything (the first time takes a while: SQLite is compiled too)
```

You should see every test program report `test result: ok` – about 430 tests in all,
including 60 shared spec cases, 44 BDD scenarios (in 7 feature files), 14 properties and
10 doc tests. All commands below are run from this `rust/` folder.

The extra tools are optional; install them once if you want to use them:

```bash
cargo install cargo-nextest --version 0.9.151 --locked
cargo install cargo-llvm-cov --version 0.9.1 --locked && rustup component add llvm-tools-preview
cargo install cargo-mutants --version 27.1.0 --locked
cargo install cargo-insta --version 1.49.0 --locked     # only to review snapshots
```

## Running the tests and checks

| I want to… | Command |
|------------|---------|
| run all tests (unit, integration, BDD, doc tests) | `cargo test` |
| run only the unit tests (the ones inside `src/`) | `cargo test --lib` |
| run one test file from `tests/` | `cargo test --test sqlite_repository` |
| run tests whose name contains some text | `cargo test half` |
| run the **BDD** scenarios from `spec/features/` | `cargo test --test features` |
| run the scenarios of one feature file | `cargo test --test features shipping` |
| run one shared spec case by its id | `cargo test --test spec_cases_named promo-coffee` |
| list the shared spec cases as tests | `cargo test --test spec_cases_named -- --list` |
| run only the **doc tests** | `cargo test --doc` |
| see `println!` output of passing tests | `cargo test -- --nocapture` |
| run the tests **in parallel, one process each**, with a JUnit report | `cargo nextest run --profile ci` → `target/nextest/ci/junit.xml` |
| measure **coverage** (fails below 100% of the lines in `src/`) | `cargo llvm-cov --ignore-filename-regex 'src/bin/' --fail-under-lines 100 --html` → `target/llvm-cov/html/index.html` |
| run **mutation testing** (about 15 minutes, see below) | `cargo mutants --in-place` → `mutants.out/` |
| … only on the lines changed since `main` | `git diff --relative main > ../my.diff && cargo mutants --in-place --in-diff ../my.diff` |
| review changed **snapshots** | `cargo insta review` |
| run the **benchmark** | `cargo bench` |
| run the **shared e2e runner** (from the repository root) | `e2e/run.sh cargo run --quiet --manifest-path rust/Cargo.toml --bin checkout` |
| … or, faster, against a binary built once | `cargo build`, then from the root: `e2e/run.sh rust/target/debug/checkout` |
| check the formatting (`cargo fmt` fixes it) | `cargo fmt --check` |
| lint (warnings are errors) | `cargo clippy --locked --all-targets -- -D warnings` |
| check the documentation for broken links | `RUSTDOCFLAGS="-D warnings" cargo doc --no-deps` |
| try the command line program | `echo '{"items":[{"sku":"APPLE","quantity":1}]}' \| cargo run --quiet --bin checkout` |

Notice the `--` in some commands: options **before** it go to Cargo (which tests to build),
options **after** it go to the test runner (which tests to run and how).

## Folder layout

Rust has its own conventions for where tests go, and they differ from most other languages:

```
rust/
├── Cargo.toml              ← project file: dependencies, test targets, benchmark
├── Cargo.lock              ← exact dependency versions (committed)
├── .config/nextest.toml    ← cargo-nextest: the `ci` profile writes JUnit XML
├── .cargo/mutants.toml     ← cargo-mutants: what is not mutated, and why
├── src/
│   ├── lib.rs              ← the public API and checkout(): the pipeline from rules §4
│   ├── basket.rs           ← Basket, BasketItem                              (rules §3)
│   ├── catalog.rs          ← Product, Promotion, Catalog, require()          (rules §2)
│   ├── receipt.rs          ← Receipt, ReceiptLine                            (rules §5)
│   ├── error.rs            ← ErrorCode, CheckoutError                        (rules §9)
│   ├── validation.rs       ← rows 1-4 of rules §4.1
│   ├── lines.rs            ← merging basket items into receipt lines, pricing a line
│   ├── promotions.rs       ← buy X get Y free, bulk discount
│   ├── coupons.rs          ← SAVE10, WELCOME5
│   ├── shipping.rs         ← standard / express
│   ├── tax.rs              ← 8%
│   ├── money.rs            ← round_half_up, with whole numbers only
│   ├── orders.rs           ← Order, Clock, OrderRepository, InMemoryOrderRepository, place_order
│   ├── json_file_repository.rs ← one ORDER-<n>.json file per order
│   ├── sqlite_repository.rs    ← one table in an SQLite file
│   ├── spec_folder.rs      ← finds ../spec/ by walking up the folder tree
│   ├── command.rs          ← the CLI's logic: JSON text in, output + exit code out
│   ├── test_fixtures.rs    ← tiny made-up catalog for the unit tests (compiled for tests only)
│   └── bin/checkout.rs     ← the CLI itself: stdin → command.rs → stdout
├── tests/                  ← integration tests: one program per .rs file
│   ├── common/             ← shared helpers (not a test program itself)
│   │   ├── mod.rs          ← load_catalog(), load_spec_cases(), spec_dir(), fixed_clock()
│   │   ├── fixtures.rs     ← tiny catalog, sample receipt, a clock the test can move
│   │   ├── schemas.rs      ← schema_errors(schema, value) with the jsonschema crate
│   │   └── dollars.rs      ← "$12.39" → 1239, without floating point
│   ├── spec_cases.rs       ← all 60 shared cases in one test (collects every failure)
│   ├── spec_cases_named.rs ← the same 60 cases, one named test each (libtest-mimic)
│   ├── spec_files.rs       ← checks that the spec files load
│   ├── features.rs         ← Cucumber step definitions for spec/features/
│   ├── properties.rs       ← P1 to P14 from spec/properties.md with proptest
│   ├── contract.rs         ← receipts, orders and CLI errors checked against spec/schema/
│   ├── snapshots.rs        ← whole receipts compared with saved snapshots …
│   ├── snapshots/          ← … which live here (*.snap, committed)
│   ├── json_file_repository.rs ← real files in a temporary folder
│   ├── sqlite_repository.rs    ← a real SQLite file in a temporary folder
│   ├── cli.rs              ← the real program run as a separate process
│   └── dollars.rs          ← tests for the "$12.39" helper (test code deserves tests)
├── benches/checkout.rs     ← a Criterion micro-benchmark
└── fuzz/                   ← a separate crate for cargo-fuzz (nightly only), target for P10
```

* **Unit tests live next to the code**, at the bottom of each `src/*.rs` file inside a
  `#[cfg(test)] mod tests { ... }` block. `#[cfg(test)]` means the block is only compiled
  for `cargo test`, so it never ends up in the real program. Because the tests sit inside
  the module, they can also test private functions.
* **Integration tests live in `tests/`.** Cargo compiles each file there as a separate
  program that can only use the crate's *public* API – just like a real user of the crate.
  Code shared between test files goes in a sub-folder (`tests/common/`), because Cargo
  would otherwise treat it as one more test program.
* **Doc tests** are the code examples in `///` comments, for example on `checkout()`,
  `round_half_up()` and `run_checkout_command()`. `cargo test` runs them, so the
  documentation can never go out of date.
* Two test files bring their own `main` (`harness = false` in `Cargo.toml`): the BDD runner
  `features.rs` and the per-case runner `spec_cases_named.rs`. Both build their list of
  tests at run time with `libtest-mimic` – one test per feature file, one per spec case – so
  they accept the usual options (`--list`, name filters) and work with cargo-nextest.
  The benchmark has its own `main` too.

The tests read the spec straight from the repository's `../spec/` folder (it is never
copied). `tests/common/mod.rs` finds it through `env!("CARGO_MANIFEST_DIR")`, the folder
containing `Cargo.toml`, so the tests work no matter which folder you start them from.

## The kinds of tests, and where to look

| Kind of test | Question it answers | Look at |
|---|---|---|
| **Unit** | Does this one function do the right thing? | the `mod tests` at the bottom of `src/money.rs` is a good first read |
| **Parameterised** | Does it work for all these examples? | `#[rstest]` with `#[case]` rows in `src/promotions.rs`, `src/validation.rs`, `src/lib.rs` |
| **Shared spec cases** | Does it match the 60 cases every language must pass? | `tests/spec_cases.rs` (one loop) and `tests/spec_cases_named.rs` (one test per case) |
| **BDD** | Does it behave as the Given/When/Then scenarios say? | `tests/features.rs` – each sentence from `spec/features/README.md` appears once; an undefined sentence fails the run |
| **Test doubles** | How do I test without the real clock? Was `save()` really never called? | `src/orders.rs`: a fake (`InMemoryOrderRepository`), a stub (`FixedClock`), mocks (`MockOrderRepository`, `MockClock` from `#[automock]`) |
| **Integration** | Do the repositories really write files and database rows? | `tests/json_file_repository.rs`, `tests/sqlite_repository.rs` |
| **Property-based** | Is it true for *every* basket, not just my examples? | `tests/properties.rs` – P1…P14 |
| **Contract** | Is the output even the right shape? | `tests/contract.rs` |
| **Snapshot** | Did the output change since last time? | `tests/snapshots.rs` and `tests/snapshots/*.snap` |
| **End-to-end** | Does the real program work from the outside? | `tests/cli.rs`, then the shared runner `e2e/run.sh` |
| **Doc tests** | Do the examples in the documentation still work? | the `///` examples in `src/lib.rs`, `src/money.rs`, `src/command.rs`, … |
| **Benchmark** | How fast is it? | `benches/checkout.rs` |
| **Fuzzing** | Does any input at all crash it? | `fuzz/fuzz_targets/basket_checkout.rs` |
| **Coverage** | Which lines did the tests never run? | `cargo llvm-cov` – CI fails below 100% of the lines in `src/` |
| **Mutation** | Would the tests notice if the code were wrong? | `cargo mutants` – see below |
| **Static checks** | Can we catch mistakes without running anything? | `cargo fmt --check`, `cargo clippy`, `cargo doc` |

### Test doubles in Rust

Rust has no reflection, so a mocking library cannot replace a function at run time. Instead,
the code depends on **traits** – `Clock` (what time is it?) and `OrderRepository`
(`save` / `find`) – and tests pass in whatever implementation they need:

* `InMemoryOrderRepository` is a **fake**: a real, working repository that keeps orders in a
  `Vec`. The BDD scenarios and most tests use it.
* `FixedClock` is a **stub**: it always answers the same time. `tests/common/fixtures.rs`
  also has a `SharedClock` the test can move while a repository holds it.
* `#[cfg_attr(test, mockall::automock)]` on the traits makes mockall generate
  `MockOrderRepository` and `MockClock` – **mocks** that check how they were called, for
  example `repository.expect_save().never()` proves a failed checkout never saves (P14).

### Snapshots with insta

`assert_json_snapshot!(receipt)` writes the receipt as JSON to `tests/snapshots/*.snap` the
first time and compares with that file afterwards. The `.snap` files are committed: they are
the approved output. When a test fails because the output changed:

1. run `cargo insta review` (needs `cargo install cargo-insta`), which shows each change as a
   diff and lets you accept or reject it, or
2. without cargo-insta: `INSTA_UPDATE=always cargo test --test snapshots`, then check the
   changed `.snap` files with `git diff` before committing.

Two snapshots are *inline*: the approved value is written in the test itself (`@r#"..."#`),
and `cargo insta review` updates the test file.

### Mutation testing, and what it found

cargo-mutants makes small changes to `src/` ("mutants": `>=` becomes `>`, `+` becomes `-`,
a function returns a default value) and runs the tests against each one. A mutant that
survives means no test noticed the change.

The score is currently **100%**: of 215 mutants, 182 are caught and 33 are *unviable*
(the change does not even compile – Rust's type checker catches those before any test
runs). Getting there shaped both the tests and the code, and the git history shows each step:

* **The first run** (214 mutants: 175 caught, 5 timeouts, 2 missed) found an untested
  `Debug` implementation, and an *equivalent mutant*: `close()` on the SQLite repository could
  be replaced by "do nothing" without any observable difference, because the database closes
  itself when the repository is dropped anyway. The fix was to delete `close()` – in Rust,
  dropping a value is how resources are released.
* **The five timeouts** were mutants that broke the calendar arithmetic in
  `format_timestamp` so that a `while` loop never ended. A timeout still counts as caught,
  but costs two minutes each; rewriting the loops over bounded ranges (`1970..=9999`,
  `1..=12`) made those mutants fail at once and cut the full run from 30 to 13 minutes.
* **The second run** (1 missed) found a real gap: `<` could become `<=` in the month
  counting, because no test used the first day of a month other than January. The last
  second of January and the 1st of March are now in the test table.

Boundaries tested on both sides (`2499` / `2500`, `4999` / `5000`, `9` / `10`, `99` / `100`,
`0` / `1`) and error **messages** tested next to error codes are what make most of the
other mutants die.

Three things are worth knowing:

* Run it with **`--in-place`**. By default cargo-mutants copies `rust/` to a temporary folder,
  where `../spec/` does not exist; in place, it mutates the real files and puts them back
  afterwards (commit first, just in case).
* `.cargo/mutants.toml` excludes `src/bin/checkout.rs` – it runs in a separate process, which
  the mutated test process cannot observe (its logic is in `command.rs`, which *is* mutated) –
  and `src/test_fixtures.rs`, which is test data.
* It takes about a quarter of an hour for ~215 mutants, because every mutant means a
  rebuild. CI therefore only tests the mutants on the lines a pull request changes
  (`--in-diff`), which takes a minute or two; the full run happens when the workflow is
  started by hand. Results are written to `mutants.out/` (`missed.txt` lists survivors).

### Fuzzing (optional, needs nightly)

`fuzz/` is a separate crate for [cargo-fuzz](https://rust-fuzz.github.io/book/cargo-fuzz.html)
(the main `Cargo.toml` excludes it, so `cargo test` never builds it). Its one target feeds
random bytes to `run_checkout_command` – the same code the CLI runs – and fails on any panic:
property P10, tested with millions of inputs instead of hundreds.

```bash
rustup toolchain install nightly
cargo install cargo-fuzz --locked
cargo +nightly fuzz run basket_checkout -- -max_total_time=60
```

CI does not fuzz – it needs a nightly compiler and is open-ended by nature – but it does
check that the fuzz target still compiles (`cargo check --manifest-path fuzz/Cargo.toml`
works on stable).

## How to add a new test

Say you want to check that `round_half_up` rounds 10% of $0.05 (0.5 cents) up to 1 cent.

1. **Pick the kind of test.** It is one function with made-up numbers: a unit test, so open
   `src/money.rs` and scroll to `mod tests` at the bottom.
2. **Write the test**, following the *arrange / act / assert* pattern:

   ```rust
   #[test]
   fn rounds_half_a_cent_up_to_one_cent() {
       // Arrange
       let amount = 5;

       // Act
       let rounded = round_half_up(amount, 10);

       // Assert
       assert_eq!(rounded, 1);
   }
   ```

   For a table of examples, add a `#[case(...)]` row to an existing `#[rstest]` test instead.
3. **Run it:** `cargo test half_a_cent`.
4. **Make sure it can fail.** Change `1` to `0` and check that the test goes red. A test that
   never fails does not test anything! (Mutation testing automates this idea.)
5. **Run everything** before you open a pull request: `cargo fmt --check`,
   `cargo clippy --all-targets -- -D warnings` and `cargo test`.

For a test that uses only the public API (or the spec files), create a new file in `tests/`
and write `mod common;` at the top if you need the helpers. Want to add a case for *every*
language instead? Add it to [`spec/test-cases/`](../spec/test-cases/): the spec-case, contract
and e2e tests pick it up automatically (update the count in `tests/spec_files.rs`).

## Things to know about this implementation

* **Whole numbers only.** `round_half_up(amount, percent)` is `(amount × percent + 50) / 100`
  with integer division, as the hint in rules §7 suggests. No floating point anywhere near money.
* **Errors are values.** `checkout` returns `Result<Receipt, CheckoutError>`; every step
  returns a `Result` too, and `?` stops at the first error – which is exactly what makes the
  error order of rules §4.1 come out right. A `CheckoutError` has a `code` (for programs) and
  a `message` (for people).
* **Parse, don't validate.** `parse_shipping_method` turns the basket's free-form `shipping`
  text into the enum `ShippingMethod`; after that the compiler knows it is one of the two.
* **Quantities are `i64`,** so "not a whole number" (`1.5`) cannot even reach the checkout:
  the CLI rejects it while reading the JSON, as "not a basket".
* **No date library.** `created_at` is written and read with the standard library only
  (`format_timestamp` / `parse_timestamp` count years and months like a paper calendar).
* **Three repositories, one trait.** Ids are `ORDER-<n>`; the JSON-file store continues after
  the highest existing file, the SQLite store uses the row number. `place_order` = checkout +
  save, and cannot save anything when checkout fails.
* **The CLI finds the catalog by walking up** from its own location, then from the current
  directory, and finally from the folder it was compiled in – so `cargo run`, the built
  binary and the shared runner all work from any folder.

## Ideas for the next contribution

Every row of the Rust column of the matrix is done, so what is left is polish:

* Compare with the TypeScript version: same spec, same kinds of tests, different idioms.
  That comparison is the point of the project.
* [`divan`](https://crates.io/crates/divan) is a smaller alternative to Criterion; a second
  benchmark with it would show the trade-off.
* The CLI could validate its input against `basket.schema.json` too (today serde accepts a
  basket written as a JSON array of its fields, which the schema does not).
