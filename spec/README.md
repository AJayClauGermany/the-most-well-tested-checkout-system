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
└── schema/                ← JSON Schemas that keep the JSON files honest
```

## Why have a shared spec?

* **Write the test once, run it six times.** A case in `test-cases/` is a fact about the
  checkout (*"3 coffees cost $17.98 after promotions"*). Each language loads the same JSON
  and checks its own implementation against it.
* **Tests before code.** This repository is set up so you can think about the
  behaviour first (TDD / BDD) and implement it afterwards.
* **One end-to-end suite for everyone.** Because every language will provide the same
  command line interface ([rules §8](checkout-rules.md#8-command-line-interface-for-end-to-end-tests)),
  a single black-box test suite can drive all of them.

## Two styles of shared test

| Folder         | Style                     | Good for                                   |
|----------------|---------------------------|--------------------------------------------|
| `test-cases/`  | Data-driven / table tests | Lots of input → output examples, edge cases |
| `features/`    | BDD (Given / When / Then) | Readable scenarios a non-programmer can check |

Both describe the same rules, so they overlap on purpose.

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
   python3 spec/check_unique_ids.py
   ```

4. Implementations that already pass the old cases may now fail. That is expected – it is
   exactly how test-driven development works.
