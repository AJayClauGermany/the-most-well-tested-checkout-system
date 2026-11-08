# Checkout – Python

The Python version of the checkout system. The rules it must follow are in
[`spec/checkout-rules.md`](../spec/checkout-rules.md).

> **Status:** test-first baseline. The data types and the tests are ready, but
> `checkout()` itself is not written yet – it raises `NotImplementedError`.
> The 60 shared spec cases are waiting (skipped) in
> [`tests/acceptance/test_spec_cases.py`](tests/acceptance/test_spec_cases.py).
> Making them pass is the next step!

## Tools we use, and why

| Tool | What it does | Why this one |
|------|--------------|--------------|
| [Python](https://www.python.org/) 3.10+ | The language | CI uses 3.12 |
| [venv](https://docs.python.org/3/library/venv.html) | A private folder of packages for this project | Built into Python; keeps this project's packages away from the rest of your computer |
| [pip](https://pip.pypa.io/) + `pyproject.toml` | Installs the project and its tools | The standard way; versions are pinned with `==` so everyone gets the same tools |
| [pytest](https://docs.pytest.org/) | Finds and runs the tests | The most popular Python test runner; tests are plain functions using plain `assert` |
| [Ruff](https://docs.astral.sh/ruff/) | Linter (finds suspicious code) and formatter (tidies layout) | Very fast, one tool for both jobs |

## Getting started

You need **Python 3.10 or newer** (`python3 --version`). From this `python/` folder:

```bash
# 1. Create a virtual environment (once)
python3 -m venv .venv

# 2. Activate it (every time you open a new terminal)
source .venv/bin/activate          # macOS / Linux
# .venv\Scripts\activate           # Windows (PowerShell / cmd)

# 3. Install the project (editable, so your code changes are picked up) and the dev tools
python -m pip install -e ".[dev]"
```

## Running the tests

| I want to run…                     | Command |
|------------------------------------|---------|
| all tests                          | `pytest` |
| all tests, one line per test       | `pytest -v` |
| only the unit tests                | `pytest tests/unit` |
| only the integration tests         | `pytest tests/integration` |
| one file                           | `pytest tests/unit/test_catalog.py` |
| one test                           | `pytest tests/unit/test_catalog.py::test_find_is_case_sensitive` |
| tests whose name contains a word   | `pytest -k coffee` |
| one shared spec case               | `pytest "tests/acceptance/test_spec_cases.py::test_spec_case[promo-coffee-three-bags-one-free]"` |
| stop at the first failure          | `pytest -x` |

Right now you should see **17 passed, 60 skipped**. The skipped ones are the
shared spec cases waiting for `checkout()` to be written.

The same checks CI runs on every pull request:

```bash
ruff check .            # lint
ruff format --check .   # formatting (run `ruff format .` to fix it)
pytest                  # tests
```

## Folder layout

```
python/
├── pyproject.toml            ← project settings, dependencies, pytest + ruff config
├── src/checkout/             ← the code (the "src layout")
│   ├── __init__.py           ← public API: `from checkout import checkout, Basket, ...`
│   ├── models.py             ← Product, Catalog, Basket, Receipt, ... (dataclasses)
│   ├── errors.py             ← ErrorCode and CheckoutError
│   └── engine.py             ← checkout(basket, catalog) – not implemented yet
└── tests/
    ├── conftest.py           ← shared fixtures (e.g. `catalog`), found by pytest automatically
    ├── spec_loader.py        ← helper that reads ../spec/catalog.json and ../spec/test-cases/
    ├── unit/                 ← small, fast tests of one piece of code, no files or network
    ├── integration/          ← tests where pieces work together (here: reading the real spec files)
    ├── acceptance/           ← the shared spec cases, run against checkout()
    ├── e2e/                  ← (to do) run the command line program as a user would
    └── bdd/                  ← (to do) Given/When/Then scenarios from ../spec/features/
```

Python conventions worth knowing:

* **pytest finds tests by name:** files called `test_*.py`, functions called `test_*`.
* **`conftest.py`** holds fixtures that every test below it can use, with no import needed.
* **src layout:** the code lives in `src/`, so tests use the *installed* package
  (`pip install -e .`), just like a real user would. This catches packaging mistakes early.
* The spec files are **read from the shared `spec/` folder**, never copied.
  `spec_loader.find_spec_folder()` finds it from its own location, so `pytest` works
  from any folder.

## How to add a new test

Say you start implementing and want to test that "buy 2, get 1 free" makes 3 coffees cost
the price of 2.

1. **Pick the folder.** Testing one small function with hand-made data? → `tests/unit/`.
2. **Create a file** whose name starts with `test_`, e.g. `tests/unit/test_promotions.py`.
3. **Write a function** whose name starts with `test_` and says what should happen:

   ```python
   from checkout import Basket, BasketItem, BuyXGetYFree, Catalog, Product, checkout

   COFFEE = Product(sku="COFFEE", name="Coffee", unit_price=899, promotion=BuyXGetYFree(buy=2, free=1))


   def test_the_third_coffee_is_free():
       # Arrange
       catalog = Catalog(products=[COFFEE])
       basket = Basket(items=[BasketItem(sku="COFFEE", quantity=3)])

       # Act
       receipt = checkout(basket, catalog)

       # Assert
       assert receipt.promotion_discount == 899
   ```

4. **Run it:** `pytest tests/unit/test_promotions.py`. Watch it fail first (red), then write
   just enough code to make it pass (green), then tidy up (refactor).
5. **Run everything** (`ruff check . && ruff format --check . && pytest`) before opening a
   pull request.

Want the same test for many inputs? Use
[`@pytest.mark.parametrize`](https://docs.pytest.org/en/stable/how-to/parametrize.html) –
see [`tests/unit/test_error_codes.py`](tests/unit/test_error_codes.py) for an example.

## Testing ideas we'd love contributions for

This list mirrors the Python column of the
[testing matrix](../docs/testing-matrix.md) – the project's single list of goals – and adds
the practical hints. Tick the matrix when you add one.

**Tier 1 – every language does these**

* **Implement `checkout()`** test-first: remove the skip in `tests/acceptance/test_spec_cases.py`
  and make the 60 cases pass, adding unit tests for each step (validation, merging,
  promotions, coupons, shipping, tax) as you go.
* **Parameterised tests** with `@pytest.mark.parametrize` (use `ids=` so each row has a
  readable name) – the native table style, separate from the shared spec cases.
* **BDD** with [pytest-bdd](https://pytest-bdd.readthedocs.io/) on [`spec/features/`](../spec/features/);
  [behave](https://behave.readthedocs.io/) is the alternative with its own runner.
  See [`tests/bdd/README.md`](tests/bdd/README.md).
* **Test doubles** for the order repository and the clock (rules §6): a hand-written
  in-memory fake first, then [`unittest.mock`](https://docs.python.org/3/library/unittest.mock.html)
  to check that a failed checkout never calls `save()`, and a fixed clock so `created_at`
  is predictable.
* **Integration tests** for a JSON-file and an [`sqlite3`](https://docs.python.org/3/library/sqlite3.html)
  order repository, using pytest's `tmp_path` fixture. Check the stored JSON against
  [`spec/schema/order.schema.json`](../spec/schema/order.schema.json).
* **Property-based testing** with [Hypothesis](https://hypothesis.readthedocs.io/): one test
  per property in [`spec/properties.md`](../spec/properties.md), named after its id (`P1`, …).
* **Contract tests** with [`jsonschema`](https://python-jsonschema.readthedocs.io/): every
  receipt must satisfy [`spec/schema/receipt.schema.json`](../spec/schema/receipt.schema.json).
* **Snapshot testing** of whole receipts with [syrupy](https://github.com/syrupy-project/syrupy).
* **End-to-end tests** of `python -m checkout` (rules §8) with
  [`subprocess`](https://docs.python.org/3/library/subprocess.html) – see
  [`tests/e2e/README.md`](tests/e2e/README.md) – and then the shared runner:
  `cd python && ../e2e/run.sh python -m checkout`.
* **Code coverage** with [pytest-cov](https://pytest-cov.readthedocs.io/), with a `fail_under`
  threshold so CI notices when it drops.
* **Mutation testing** with [mutmut](https://mutmut.readthedocs.io/) (uses `fork()`: on Windows,
  run it in WSL). Try swapping the rounding for Python's built-in `round()` – banker's
  rounding! – and see which spec cases catch it.
* **Static type checking** with [mypy](https://mypy.readthedocs.io/) or
  [pyright](https://microsoft.github.io/pyright/) – the code already has type hints.

**Tier 2 – idiomatic extras**

* **Doctests** with [`doctest`](https://docs.python.org/3/library/doctest.html): examples in
  docstrings (e.g. of `round_half_up`) that are also tests; run with `pytest --doctest-modules`.

**Tier 3 – plumbing**

* `pytest --junitxml=report.xml` for a CI test report; [pytest-xdist](https://pytest-xdist.readthedocs.io/)
  (`-n auto`) and [pytest-randomly](https://github.com/pytest-dev/pytest-randomly) to prove tests
  are independent of each other; `--strict-markers` and `--strict-config` in `pyproject.toml`.
