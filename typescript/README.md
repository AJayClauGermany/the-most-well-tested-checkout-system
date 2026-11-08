# Checkout system – TypeScript

The TypeScript version of the checkout described in [`spec/checkout-rules.md`](../spec/checkout-rules.md).

> **Status:** test-first. The types and tests are here; `checkout()` itself is **not
> implemented yet**. The 60 shared spec cases are skipped and waiting for you.

## Tools we use (and why)

| Tool | What it does | Why this one |
|------|--------------|--------------|
| [Node.js](https://nodejs.org/) 22 or newer | Runs JavaScript (and, since Node 22.18, TypeScript files directly) | The standard JavaScript runtime outside the browser |
| [TypeScript](https://www.typescriptlang.org/) | Adds types to JavaScript and checks them | Catches many bugs before the code even runs. We use `strict` mode |
| [Vitest](https://vitest.dev/) | Test runner and assertion library | Fast, runs `.ts` files with no extra setup, and has a Jest-style API (`describe` / `it` / `expect`) that most JavaScript developers already know |
| npm | Installs the tools above | Comes with Node.js. `package-lock.json` pins every exact version |

There is no build step: Vitest runs the `.ts` files directly, and `tsc` (the TypeScript
compiler) is only used to *check* types (`--noEmit`).

## Getting started

You need **Node.js 22.12 or newer** (CI uses Node 24). Check with `node --version`.

```bash
cd typescript
npm ci            # install the exact versions from package-lock.json
npm test          # run all tests once
```

You should see something like `11 passed | 60 skipped`.

## Running tests

| I want to… | Command |
|------------|---------|
| run all tests once | `npm test` |
| re-run tests automatically whenever I save a file | `npm run test:watch` |
| run one file | `npm test -- tests/unit/catalog.test.ts` |
| run only the unit tests (a folder) | `npm test -- tests/unit` |
| run tests whose name contains some text | `npm test -- -t "finds a product"` |
| see every test name (including skipped ones) | `npm test -- --reporter=verbose` |
| check the types (also checks the tests) | `npm run typecheck` |

Everything after `--` is passed on to Vitest. Run `npx vitest --help` for more options.

Tip: to run just one test while you work on it, change `it(...)` to `it.only(...)`
(remember to change it back before you commit).

## Folder layout

```
typescript/
├── package.json          ← scripts and pinned tool versions
├── package-lock.json     ← exact versions of everything installed (commit it!)
├── tsconfig.json         ← TypeScript settings (strict mode)
├── src/                  ← the checkout code
│   ├── index.ts          ← public API: import everything from here
│   ├── types.ts          ← Product, Promotion, Catalog, Basket, Receipt, …
│   ├── errors.ts         ← ErrorCode, ERROR_CODES and CheckoutError
│   ├── catalog.ts        ← findProduct(catalog, sku)
│   └── checkout.ts       ← checkout(basket, catalog) – not implemented yet
└── tests/
    ├── unit/             ← small, fast tests of one piece of code, with made-up data
    ├── integration/      ← tests that use real things, like the JSON files in spec/
    │   ├── spec-loader.test.ts  ← checks that the spec files load correctly
    │   └── spec-cases.test.ts   ← the 60 shared cases (skipped for now)
    ├── e2e/              ← end-to-end tests of the command line program (coming soon)
    └── support/          ← helpers for tests, not tests themselves
        └── spec.ts       ← loadCatalog() and loadSpecCases()
```

In the JavaScript world, tests usually live either next to the code
(`checkout.test.ts` beside `checkout.ts`) or in a separate `tests/` folder.
We use a separate folder so the different *kinds* of test are easy to tell apart.
Vitest finds every file whose name ends in `.test.ts`, wherever it is.

The shared data is read straight from the repository's [`spec/`](../spec/) folder
(never copied), so every language is tested against exactly the same cases.

## How to add a new test

Say you want to check that `findProduct` finds `"PEN"` in the standard catalog.

1. **Pick the kind of test.** It only needs a catalog, so it is a unit test:
   open `tests/unit/catalog.test.ts` (or create a new file ending in `.test.ts`).
2. **Write the test**, following the *arrange / act / assert* pattern:

   ```ts
   it("finds the pen", () => {
     // Arrange
     const catalog: Catalog = {
       products: [{ sku: "PEN", name: "Pen", unit_price: 120, promotion: null }],
     };

     // Act
     const product = findProduct(catalog, "PEN");

     // Assert
     expect(product?.unit_price).toBe(120);
   });
   ```

3. **Run it:** `npm test -- -t "finds the pen"`.
4. **Make sure it can fail.** Change `120` to `121` and check that the test goes red.
   A test that never fails does not test anything!
5. **Run everything** before you open a pull request: `npm run typecheck && npm test`.

Want to add a case for *every* language instead? Add it to
[`spec/test-cases/`](../spec/test-cases/) and it is picked up automatically by
`tests/integration/spec-cases.test.ts` (and update the count in `spec-loader.test.ts`).

## Starting the implementation (TDD)

1. In `tests/integration/spec-cases.test.ts`, change `it.skip.each` to `it.each`.
2. `npm test` – all 60 cases fail. That is the "red" step.
3. Write code in `src/checkout.ts` until they pass ("green"), then tidy up ("refactor").
   Small helper functions (rounding, promotions, coupons, shipping, tax) deserve their
   own unit tests in `tests/unit/`.

## Testing ideas we'd love contributions for

This list mirrors the TypeScript column of the
[testing matrix](../docs/testing-matrix.md) – the project's single list of goals – and adds
the practical hints. Tick the matrix when you add one. Note which Vitest major you are on:
Stryker and vitest-cucumber both track it.

**Tier 1 – every language does these**

* **Implement `checkout()`** test-first, with unit tests for each calculation step
  (rounding, merging, promotions, coupons, shipping, tax).
* **Parameterised tests** with `test.each` / `it.each` – the native table style, separate
  from the shared spec cases.
* **BDD** with [`@cucumber/cucumber`](https://github.com/cucumber/cucumber-js) on
  [`spec/features/`](../spec/features/), step definitions in `tests/bdd/`.
  [`@amiceli/vitest-cucumber`](https://github.com/amiceli/vitest-cucumber) stays inside Vitest
  but is a one-maintainer project that bumps a major with every Vitest major.
* **Test doubles** for the order repository and the clock (rules §6): a hand-written
  in-memory fake, then [`vi.fn()`](https://vitest.dev/api/mock.html) to check a failed
  checkout never calls `save()`, and a fixed clock for `created_at`.
* **Integration tests** for a JSON-file repository (`node:fs` + `os.tmpdir()`) and an SQLite
  one ([`node:sqlite`](https://nodejs.org/api/sqlite.html), or
  [better-sqlite3](https://github.com/WiseLibs/better-sqlite3)). Check stored orders against
  [`spec/schema/order.schema.json`](../spec/schema/order.schema.json).
* **Property-based testing** with [fast-check](https://fast-check.dev/): one test per property
  in [`spec/properties.md`](../spec/properties.md), named after its id.
* **Contract tests** with [ajv](https://ajv.js.org/): every receipt must satisfy
  [`spec/schema/receipt.schema.json`](../spec/schema/receipt.schema.json).
* **Snapshot testing** with [`toMatchSnapshot()` / `toMatchInlineSnapshot()`](https://vitest.dev/guide/snapshot.html)
  for whole receipts.
* **End-to-end tests** of the CLI (`src/cli.ts`, rules §8) with `node:child_process`
  ([execa](https://github.com/sindresorhus/execa) is a friendlier wrapper) – see
  [`tests/e2e/README.md`](tests/e2e/README.md) – and then the shared runner in [`e2e/`](../e2e/).
* **Coverage** with [`@vitest/coverage-v8`](https://vitest.dev/guide/coverage.html) and
  thresholds in `vitest.config.ts`.
* **Mutation testing** with [StrykerJS](https://stryker-mutator.io/docs/stryker-js/vitest-runner/)
  and its Vitest runner.
* **Linting and formatting** with [ESLint](https://eslint.org/) + [typescript-eslint](https://typescript-eslint.io/)
  and [Prettier](https://prettier.io/), added to CI ([Biome](https://biomejs.dev/) is the one-binary alternative).

**Tier 2 – idiomatic extras**

* **Type tests** with [`expectTypeOf`](https://vitest.dev/guide/testing-types.html) and
  `vitest --typecheck`, e.g. a `Promotion` must have a valid `type`.
* **Benchmarks** with [`vitest bench`](https://vitest.dev/guide/features.html#benchmarking) – one
  benchmark of `checkout()` is plenty.

**Tier 3 – plumbing**

* `vitest run --reporter=junit --outputFile=report.xml` for a CI test report; Vitest already
  runs files in parallel, so shared state between tests shows up quickly.
