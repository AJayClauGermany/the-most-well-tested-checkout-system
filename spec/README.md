# Language-agnostic specification

Everything in this folder describes **what** the checkout system does, without
saying **how** any particular language does it. Every implementation in this
repository (Python, TypeScript, Go, Java, C#, Rust) reads from here.

```
spec/
├── checkout-rules.md      ← the rules, in plain English (start here!)
├── catalog.json           ← the standard product catalog
├── test-cases/            ← data-driven test cases (JSON): basket in, totals or error out
├── features/              ← behaviour-driven scenarios (Gherkin / Cucumber)
├── properties.md          ← rules that hold for EVERY basket (property-based testing)
├── examples/              ← a hand-checked receipt, saved order and CLI error
└── schema/                ← JSON Schemas for every JSON shape: inputs, outputs and test files
```

The shared **end-to-end runner** that drives each language's command line program lives
one level up, in [`e2e/`](../e2e/).

## Why have a shared spec?

* **Write the test once, run it six times.** A case in `test-cases/` is a fact about the
  checkout (*"3 coffees cost $17.98 after promotions"*). Each language loads the same JSON
  and checks its own implementation against it.
* **Tests before code.** This repository is set up so you can think about the
  behaviour first (TDD / BDD) and implement it afterwards.
* **One end-to-end suite for everyone.** Because every language will provide the same
  command line interface ([rules §8](checkout-rules.md#8-command-line-interface-for-end-to-end-tests)),
  a single black-box test suite can drive all of them.

## Four styles of shared test

| What            | Style                      | Good for                                   |
|-----------------|----------------------------|--------------------------------------------|
| `test-cases/`   | Data-driven / table tests  | Lots of input → output examples, edge cases |
| `features/`     | BDD (Given / When / Then)  | Readable scenarios a non-programmer can check |
| `properties.md` | Property-based tests       | Rules that must hold for *any* basket, found by random generation |
| `schema/`       | Contract (schema) tests    | "Is this receipt even the right shape?" – checked by every language and by the e2e runner |

They all describe the same rules, so they overlap on purpose: the point of this project is
to see the *same* behaviour tested in many different ways.

## Changing the spec

1. Change [`checkout-rules.md`](checkout-rules.md) first.
2. Add or update cases in `test-cases/` and/or `features/`.
   Work the numbers out by hand (show your working in the `description` if it helps!).
3. CI validates every JSON file against [`schema/`](schema/) on your pull request.
   To check locally:

   ```bash
   pip install check-jsonschema
   check-jsonschema --schemafile spec/schema/test-cases.schema.json spec/test-cases/*.json
   check-jsonschema --schemafile spec/schema/catalog.schema.json spec/catalog.json
   check-jsonschema --schemafile spec/schema/receipt.schema.json spec/examples/receipt.json
   check-jsonschema --schemafile spec/schema/order.schema.json spec/examples/order.json
   check-jsonschema --schemafile spec/schema/cli-error.schema.json spec/examples/cli-error.json
   python3 spec/check_unique_ids.py
   ```

4. Implementations that already pass the old cases may now fail. That is expected – it is
   exactly how test-driven development works.
