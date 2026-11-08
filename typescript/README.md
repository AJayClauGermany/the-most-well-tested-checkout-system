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

Each of these shows off a different kind of testing. Pick one and open a pull request!

* **BDD with Gherkin** – run the scenarios in [`spec/features/`](../spec/features/) with
  [`@cucumber/cucumber`](https://github.com/cucumber/cucumber-js) (or
  [`@amiceli/vitest-cucumber`](https://github.com/amiceli/vitest-cucumber) to stay inside Vitest),
  e.g. step definitions in `tests/bdd/`.
* **Property-based testing** – generate thousands of random baskets with
  [fast-check](https://fast-check.dev/) and check rules that must *always* hold
  (e.g. `total = subtotal − promotion_discount − coupon_discount + shipping + tax`,
  discounts never exceed the subtotal, rounding matches rules §7).
* **Unit tests for each calculation step** – rounding, promotions, coupons, shipping, tax.
* **Mocks and fakes** – an in-memory order repository (rules §6) and tests that use
  Vitest's built-in [`vi.fn()`](https://vitest.dev/api/mock.html) to check a failed
  checkout never calls `save()`.
* **Integration tests** – a JSON-file order repository tested against a real temporary folder
  (`node:fs` + `os.tmpdir()`).
* **Snapshot testing** – Vitest's [`toMatchSnapshot()` / `toMatchInlineSnapshot()`](https://vitest.dev/guide/snapshot.html)
  for whole receipts.
* **End-to-end tests** – the command line program from rules §8 (`src/cli.ts`), driven with
  `node:child_process` (see [`tests/e2e/README.md`](tests/e2e/README.md)).
* **Coverage** – [`@vitest/coverage-v8`](https://vitest.dev/guide/coverage.html) and `vitest run --coverage`.
* **Mutation testing** – [StrykerJS](https://stryker-mutator.io/) with its Vitest runner,
  to find tests that pass even when the code is broken.
* **Type tests** – Vitest's [`expectTypeOf`](https://vitest.dev/guide/testing-types.html)
  to check that, for example, a `Promotion` must have a valid `type`.
* **Linting and formatting** – [ESLint](https://eslint.org/) with
  [typescript-eslint](https://typescript-eslint.io/) and [Prettier](https://prettier.io/),
  added to CI.
