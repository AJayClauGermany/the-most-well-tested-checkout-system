# BDD tests

**Behaviour-driven development (BDD)** tests are written as plain-English
*Given / When / Then* scenarios that a non-programmer can read and check.

The scenarios live in [`spec/features/`](../../../spec/features/) and are shared by every
language. [`test_features.py`](test_features.py) is the Python *glue code*: it loads the
`.feature` files straight from that folder (never a copy) with
[pytest-bdd](https://pytest-bdd.readthedocs.io/)'s `scenarios(...)`, and defines one *step*
function for each sentence listed in
[`spec/features/README.md`](../../../spec/features/README.md#step-vocabulary).

Because pytest-bdd turns every scenario into an ordinary pytest test, everything you know
already works:

```bash
pytest tests/bdd -v                 # one line per scenario
pytest tests/bdd -k "WELCOME5"      # only scenarios whose name contains WELCOME5
```

Things worth noticing in the step definitions:

* A `world` fixture (a small dataclass) carries state from one step to the next; every
  scenario gets a fresh one.
* `$12.39` is turned into `1239` cents by a custom `Money` type for `parsers.parse`
  (see `tests/support/dollars.py`) – with whole-number maths only, never `float`.
* The data table under `Given a basket containing:` arrives as the `datatable` argument:
  a list of rows, the first of which is the header.
