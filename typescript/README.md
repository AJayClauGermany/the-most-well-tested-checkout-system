# Checkout system – TypeScript

The TypeScript version of the checkout described in [`spec/checkout-rules.md`](../spec/checkout-rules.md).

> **Status: done, and the first language to get there.** `checkout()`, saving orders
> (in memory, as JSON files and in SQLite) and the command line program are implemented,
> and every row of the TypeScript column of the [testing matrix](../docs/testing-matrix.md)
> is ticked. It was built strictly test-first: the commit that adds the tests comes before
> the commit that makes them pass – have a look at the git history.

## Tools we use (and why)

| Tool | What it does | Why this one |
|------|--------------|--------------|
| [Node.js](https://nodejs.org/) 22.18 or newer | Runs JavaScript – and `.ts` files directly, by stripping the types | The standard JavaScript runtime outside the browser. We never compile: `node src/cli.ts` just works |
| [TypeScript](https://www.typescriptlang.org/) 7 | Adds types to JavaScript and checks them | Catches many bugs before the code even runs. We use `strict` mode. Version 7 is the new native compiler – fast, but see the note on Biome below |
| [Vitest](https://vitest.dev/) | Test runner and assertion library; also coverage, snapshots, mocks, fake timers, type tests and benchmarks | Fast, runs `.ts` files with no extra setup, and has a Jest-style API (`describe` / `it` / `expect`) most JavaScript developers already know |
| [fast-check](https://fast-check.dev/) | Property-based testing | The standard property-based library for JavaScript; integrates with any runner |
| [ajv](https://ajv.js.org/) | JSON Schema validation (contract tests) | The most used JSON Schema validator; supports draft 2020-12, which the spec schemas use |
| [Cucumber](https://github.com/cucumber/cucumber-js) (`@cucumber/cucumber`) | Runs the Given/When/Then scenarios in `spec/features/` | The reference Gherkin implementation for JavaScript |
| [`node:sqlite`](https://nodejs.org/api/sqlite.html) | SQLite, built into Node.js | A real database with nothing to install |
| [StrykerJS](https://stryker-mutator.io/) | Mutation testing | The mutation testing tool for JavaScript; has a Vitest runner |
| [Biome](https://biomejs.dev/) | Linter and formatter, one binary | The matrix recommends ESLint + Prettier, but typescript-eslint needs the TypeScript *JavaScript* API, which TypeScript 7 (the native compiler) no longer has. Biome does both jobs and does not need it |
| npm | Installs the tools above | Comes with Node.js. `package-lock.json` pins every exact version |

## Getting started

You need **Node.js 22.18 or newer** (CI uses Node 24). Check with `node --version`.

```bash
cd typescript
npm ci            # install the exact versions from package-lock.json
npm test          # run all tests once
```

You should see `22 passed` test files and something like `435 passed` tests.

## Running the tests and checks

| I want to… | Command |
|------------|---------|
| run all tests once | `npm test` |
| re-run tests automatically whenever I save a file | `npm run test:watch` |
| run one folder or file | `npm test -- tests/unit` or `npm test -- tests/unit/money.test.ts` |
| run tests whose name contains some text | `npm test -- -t "half"` |
| see every test name | `npm test -- --reporter=verbose` |
| run the tests with a **coverage** report (fails below 100%) | `npm run test:coverage` → `reports/coverage/index.html` |
| run the **type-level** tests | `npm run test:types` |
| run the **BDD** scenarios from `spec/features/` | `npm run test:bdd` |
| run the **shared e2e runner** against our CLI | `npm run test:e2e:shared` |
| run **mutation testing** (about two minutes) | `npm run test:mutation` → `reports/mutation/index.html` |
| run the **benchmark** | `npm run bench` |
| check the types (code and tests) | `npm run typecheck` |
| lint and check formatting | `npm run lint` (fix formatting with `npm run format`) |
| everything CI runs, except mutation testing | `npm run check` |
| try the command line program | `echo '{"items":[{"sku":"APPLE","quantity":1}]}' \| node src/cli.ts` |

Everything after `--` is passed on to Vitest. Run `npx vitest --help` for more options.
`npm test` also writes a JUnit XML report to `reports/junit.xml` for CI dashboards.

## Folder layout

```
typescript/
├── package.json            ← scripts and pinned tool versions
├── package-lock.json       ← exact versions of everything installed (commit it!)
├── tsconfig.json           ← TypeScript settings (strict mode)
├── vitest.config.ts        ← Vitest: test folders, coverage thresholds, reporters
├── stryker.config.json     ← StrykerJS mutation testing
├── cucumber.json           ← Cucumber: where the features and step definitions are
├── biome.json              ← Biome: lint rules and formatting style
├── src/                    ← the checkout code
│   ├── index.ts            ← public API: import everything from here
│   ├── types.ts            ← Product, Promotion, Catalog, Basket, Receipt, …
│   ├── errors.ts           ← ErrorCode, ERROR_CODES and CheckoutError
│   ├── checkout.ts         ← checkout(basket, catalog): the pipeline from rules §4
│   ├── validation.ts       ← rows 1-4 of rules §4.1
│   ├── catalog.ts          ← findProduct / requireProduct
│   ├── lines.ts            ← merging basket items into receipt lines, pricing a line
│   ├── promotions.ts       ← buy X get Y free, bulk discount
│   ├── coupons.ts          ← SAVE10, WELCOME5
│   ├── shipping.ts         ← standard / express
│   ├── tax.ts              ← 8%
│   ├── money.ts            ← roundHalfUp, with whole numbers only
│   ├── orders.ts           ← Order, Clock, OrderRepository, InMemoryOrderRepository, placeOrder
│   ├── json-file-repository.ts  ← one ORDER-<n>.json file per order
│   ├── sqlite-repository.ts     ← one table in an SQLite file (node:sqlite)
│   ├── spec-folder.ts      ← finds ../spec/ from anywhere inside the repository
│   ├── command.ts          ← the CLI's logic: JSON text in, output + exit code out
│   └── cli.ts              ← the CLI itself: stdin → command.ts → stdout
└── tests/
    ├── unit/               ← one small piece of code at a time, with made-up data
    ├── integration/        ← real things: the spec files, a temp folder, an SQLite file
    ├── contract/           ← receipts and CLI output checked against spec/schema/ with ajv
    ├── property/           ← P1 to P14 from spec/properties.md with fast-check
    ├── snapshot/           ← whole receipts compared with saved snapshots
    ├── e2e/                ← the real program run as a separate process
    ├── bdd/                ← Cucumber step definitions for spec/features/
    ├── types/              ← type-level tests (*.test-d.ts), run by `vitest --typecheck`
    ├── bench/              ← a micro-benchmark (*.bench.ts), run by `vitest bench`
    └── support/            ← helpers for tests, not tests themselves
        ├── spec.ts         ← loadCatalog() and loadSpecCases()
        ├── schemas.ts      ← schemaErrors(schema, value) with ajv
        ├── fixtures.ts     ← a tiny made-up catalog and a sample receipt
        └── dollars.ts      ← "$12.39" → 1239, without floating point
```

Tests live in a separate `tests/` folder (rather than next to the code) so the different
*kinds* of test are easy to tell apart – one folder per kind. Vitest finds every file
whose name ends in `.test.ts`, wherever it is.

The shared data is read straight from the repository's [`spec/`](../spec/) folder (never
copied), so every language is tested against exactly the same cases.

## The kinds of tests, and where to look

| Kind of test | Question it answers | Look at |
|---|---|---|
| **Unit** | Does this one function do the right thing? | `tests/unit/money.test.ts` is a good first read |
| **Parameterised** | Does it work for all these examples? | `it.each([...])` tables in `tests/unit/promotions.test.ts` |
| **Shared spec cases** | Does it match the 60 cases every language must pass? | `tests/integration/spec-cases.test.ts` |
| **BDD** | Does it behave as the Given/When/Then scenarios say? | `tests/bdd/steps.ts` – each sentence from `spec/features/README.md` appears once |
| **Test doubles** | How do I test without the real clock? Was `save()` really never called? | `tests/unit/orders.test.ts`: a fake (`InMemoryOrderRepository`), a stub clock, `vi.fn()` mocks, `vi.useFakeTimers()` |
| **Integration** | Do the repositories really write files and database rows? | `tests/integration/json-file-repository.test.ts`, `sqlite-repository.test.ts` |
| **Property-based** | Is it true for *every* basket, not just my examples? | `tests/property/checkout.properties.test.ts` – P1…P14 |
| **Contract** | Is the output even the right shape? | `tests/contract/receipt-schema.test.ts` |
| **Snapshot** | Did the output change since last time? | `tests/snapshot/receipt.test.ts` and its `__snapshots__/` |
| **End-to-end** | Does the real program work from the outside? | `tests/e2e/cli.test.ts`, then `npm run test:e2e:shared` |
| **Type-level** | Does TypeScript agree the API has this shape? | `tests/types/types.test-d.ts` |
| **Benchmark** | How fast is it? | `tests/bench/checkout.bench.ts` |
| **Coverage** | Which lines did the tests never run? | `npm run test:coverage` – thresholds are 100% in `vitest.config.ts` |
| **Mutation** | Would the tests notice if the code were wrong? | `npm run test:mutation` – see below |
| **Static checks** | Can we catch mistakes without running anything? | `npm run typecheck`, `npm run lint` |

### Mutation testing, and why the code looks the way it does

StrykerJS makes ~290 small changes to `src/` ("mutants": `>=` becomes `>`, `+` becomes `-`,
a condition becomes `true`, a string becomes `""`) and runs the tests against each one. A
mutant that survives means no test noticed the change. The score is currently **100%**, and
CI fails below 95%.

Getting there shaped both the tests and the code:

* Boundaries are tested on both sides (`2499` / `2500`, `4999` / `5000`, `9` / `10`,
  `99` / `100`), because `>` vs `>=` is the classic surviving mutant.
* Error **messages** are tested, not only error codes – otherwise Stryker shows that the
  message could be anything.
* A few lines were rewritten to remove *equivalent mutants* (changes that cannot be
  observed): for example `Math.max(0, ...numbers)` replaced a loop with a redundant
  `undefined` check.
* Four mutants in `spec-folder.ts` are reported as *runtime errors* rather than killed:
  they break the spec loader while test files are still being imported, so no test gets to
  run. Stryker leaves those out of the score.
* `src/cli.ts` is excluded: it runs in a separate process (see `tests/e2e/`), which an
  in-process tool cannot observe. Its logic is in `command.ts`, which *is* mutated.

Two more things worth knowing: Stryker copies the project into a sandbox folder before
mutating, so the code finds `spec/` by walking **up** the folder tree (`src/spec-folder.ts`)
instead of hard-coding `../../spec`; and because TypeScript 7 has no JavaScript API,
`tsconfig.json` is kept out of that sandbox (Stryker would otherwise try to rewrite it with
the compiler API). Both are explained in `stryker.config.json`.

## How to add a new test

Say you want to check that `roundHalfUp` rounds 10% of $0.05 (0.5 cents) up to 1 cent.

1. **Pick the kind of test.** It is one function with made-up numbers: a unit test, so open
   `tests/unit/money.test.ts` (or create a new file ending in `.test.ts`).
2. **Write the test**, following the *arrange / act / assert* pattern:

   ```ts
   it("rounds half a cent up to one cent", () => {
     // Arrange
     const amount = 5;

     // Act
     const rounded = roundHalfUp(amount, 10);

     // Assert
     expect(rounded).toBe(1);
   });
   ```

3. **Run it:** `npm test -- -t "half a cent"`.
4. **Make sure it can fail.** Change `1` to `0` and check that the test goes red.
   A test that never fails does not test anything! (Mutation testing automates this idea.)
5. **Run everything** before you open a pull request: `npm run check`, and
   `npm run test:mutation` if you changed `src/`.

Want to add a case for *every* language instead? Add it to
[`spec/test-cases/`](../spec/test-cases/) and it is picked up automatically by the spec
cases, contract, e2e and shared-runner tests (and update the count in `spec-loader.test.ts`).

## Things to know about this implementation

* **Whole numbers only.** `roundHalfUp(amount, percent)` is `Math.floor((amount × percent + 50) / 100)`,
  as the hint in rules §7 suggests. No floating point anywhere near money.
* **Parse, don't validate.** `parseShippingMethod` turns the basket's free-form `shipping`
  string into the type `"standard" | "express"` or throws; after that the type system knows
  the value is fine.
* **The clock is a function.** `type Clock = () => Date`. Production uses `systemClock`;
  tests pass `() => new Date("2026-01-15T10:30:00Z")`, or a `vi.fn()` to check it was called.
* **Repositories are an interface** (`save` and `find`), with three implementations. The
  in-memory one doubles as the test fake. Ids are `ORDER-<n>`; the JSON-file store takes the
  next number after the highest existing file, the SQLite store uses the row id.
* **`placeOrder` = checkout + save.** It cannot save anything when checkout throws, which is
  what the mock tests and property P14 check.

## Ideas for the next contribution

Tier 1 and Tier 2 of the matrix are complete for TypeScript, so what is left is polish:

* Compare with another language once it is implemented: same spec, same kinds of tests,
  different idioms. That comparison is the point of the project.
* A second BDD runner inside Vitest ([`@amiceli/vitest-cucumber`](https://github.com/amiceli/vitest-cucumber))
  would show the trade-off against the reference Cucumber runner.
* A JSON Schema for the *input* is already checked for the spec files; the CLI could validate
  its stdin against `basket.schema.json` too and report a seventh, CLI-only error.
