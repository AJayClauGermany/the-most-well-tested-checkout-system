# Checkout – Python

The Python version of the checkout described in [`spec/checkout-rules.md`](../spec/checkout-rules.md).

> **Status: done.** `checkout()`, saving orders (in memory, as JSON files and in SQLite) and
> the command line program are implemented, and every row of the Python column of the
> [testing matrix](../docs/testing-matrix.md) is covered. It was built test-first: the commit
> that adds the tests comes before the commit that makes them pass, and a third commit shows
> what mutation testing found afterwards – have a look at the git history.

## Tools we use (and why)

| Tool | What it does | Why this one |
|------|--------------|--------------|
| [Python](https://www.python.org/) 3.10+ | The language | CI uses 3.12. The checkout itself uses only the standard library (`json`, `sqlite3`, `re`, `datetime`) |
| [venv](https://docs.python.org/3/library/venv.html) + [pip](https://pip.pypa.io/) | A private folder of packages; installs the project and its tools | Built into Python. Every tool is pinned with `==` in `pyproject.toml`, so everybody gets the same versions |
| [pytest](https://docs.pytest.org/) | Finds and runs the tests; fixtures, `parametrize`, `tmp_path`, `monkeypatch`, `capsys`, doctests | The most popular Python test runner; tests are plain functions using plain `assert` |
| [pytest-bdd](https://pytest-bdd.readthedocs.io/) | Runs the Given/When/Then scenarios in `spec/features/` | Each scenario becomes an ordinary pytest test, so fixtures and `-k` just work. ([behave](https://behave.readthedocs.io/) is the alternative with its own runner) |
| [`unittest.mock`](https://docs.python.org/3/library/unittest.mock.html) | Mocks: `Mock`, `MagicMock`, `create_autospec` | In the standard library – nothing to install |
| [Hypothesis](https://hypothesis.readthedocs.io/) | Property-based testing | *The* property-based library for Python; shrinks a failure to the smallest example |
| [jsonschema](https://python-jsonschema.readthedocs.io/) + [referencing](https://referencing.readthedocs.io/) | JSON Schema validation (contract tests) | Supports draft 2020-12, which the spec schemas use; `referencing` lets `order.schema.json` find `receipt.schema.json` by file name, with no network |
| [syrupy](https://github.com/syrupy-project/syrupy) | Snapshot tests (the `snapshot` fixture) | The maintained snapshot plugin for pytest; can store snapshots as readable `.json` |
| [pytest-cov](https://pytest-cov.readthedocs.io/) | Code coverage, lines **and** branches | Wraps [coverage.py](https://coverage.readthedocs.io/); `fail_under = 100` makes the run fail when coverage drops |
| [pytest-xdist](https://pytest-xdist.readthedocs.io/) + [pytest-randomly](https://github.com/pytest-dev/pytest-randomly) | Parallel tests (`-n auto`); a random test order on every run | Together they prove the tests do not depend on each other |
| [mutmut](https://mutmut.readthedocs.io/) | Mutation testing | The mutation testing tool for Python. Uses `fork()`, so on Windows run it in WSL |
| [Ruff](https://docs.astral.sh/ruff/) | Linter and formatter | Very fast, one tool for both jobs |
| [mypy](https://mypy.readthedocs.io/) | Static type checker, in `--strict` mode, for the code **and** the tests | The reference type checker for Python |

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

# 4. Run the tests
pytest
```

You should see **572 passed** in a few seconds.

## Running the tests and checks

| I want to… | Command |
|------------|---------|
| run all tests (and the doctests in `src/`) | `pytest` |
| run them in parallel, on every CPU core | `pytest -n auto` |
| see every test name | `pytest -v` |
| run one kind of test | `pytest tests/unit` (or `tests/bdd`, `tests/property`, `tests/e2e`, …) |
| run one file / one test | `pytest tests/unit/test_money.py` / `pytest tests/unit/test_money.py::test_gives_0_for_0_percent` |
| run tests whose name contains a word | `pytest -k coffee` |
| run one shared spec case | `pytest "tests/acceptance/test_spec_cases.py::test_spec_case[promo-coffee-three-bags-one-free]"` |
| stop at the first failure | `pytest -x` |
| run the tests in file order (pytest-randomly normally shuffles them) | `pytest -p no:randomly` |
| repeat the exact order of an earlier run | `pytest --randomly-seed=<seed>` (the seed is printed at the top of every run) |
| run the tests with a **coverage** report (fails below 100%) | `pytest --cov --cov-report=term-missing --cov-report=html:reports/coverage` → `reports/coverage/index.html` |
| write a **JUnit XML** report for CI | `pytest --junitxml=reports/junit.xml` |
| update the **snapshots** after an intended change | `pytest tests/snapshot --snapshot-update` (then review and commit the diff) |
| run the **shared e2e runner** against our CLI | `../e2e/run.sh python -m checkout` (with the venv active) |
| run **mutation testing** (about 30 seconds) | `mutmut run`, then `mutmut results` / `mutmut show <name>` / `mutmut browse` |
| check the types (code and tests) | `mypy` |
| lint and check formatting | `ruff check .` and `ruff format --check .` (fix formatting with `ruff format .`) |
| try the command line program | `echo '{"items":[{"sku":"APPLE","quantity":1}]}' \| python -m checkout` |

pytest-randomly prints the seed it used at the top of every run (`Using --randomly-seed=1234`).
If a test only fails in some orders, it secretly depends on another test: re-run that exact
order with `pytest --randomly-seed=1234` and look for shared state.

Everything CI runs (see [`.github/workflows/python.yml`](../.github/workflows/python.yml)):

```bash
ruff check . && ruff format --check . && mypy
pytest -n auto --junitxml=reports/junit.xml --cov --cov-report=term-missing \
       --cov-report=xml:reports/coverage.xml --cov-report=html:reports/coverage
../e2e/run.sh python -m checkout
mutmut run && mutmut results
```

## Folder layout

```
python/
├── pyproject.toml              ← project, pinned dev tools, and the settings of pytest, coverage,
│                                  mutmut, ruff and mypy - all in one file
├── src/checkout/               ← the code (the "src layout")
│   ├── __init__.py             ← public API: `from checkout import checkout, Basket, ...`
│   ├── models.py               ← Product, Catalog, Basket, Receipt, ... (frozen dataclasses)
│   ├── errors.py               ← ErrorCode and CheckoutError
│   ├── engine.py               ← checkout(basket, catalog): the pipeline from rules §4
│   ├── validation.py           ← rows 1-4 of rules §4.1
│   ├── catalog.py              ← require_product
│   ├── lines.py                ← merging basket items into receipt lines, pricing a line
│   ├── promotions.py           ← buy X get Y free, bulk discount
│   ├── coupons.py              ← SAVE10, WELCOME5
│   ├── shipping.py             ← standard / express
│   ├── tax.py                  ← 8%
│   ├── money.py                ← round_half_up, with whole numbers only (and doctests)
│   ├── orders.py               ← Order, Clock, OrderRepository, InMemoryOrderRepository, place_order
│   ├── json_file_repository.py ← one ORDER-<n>.json file per order
│   ├── sqlite_repository.py    ← one table in an SQLite file (sqlite3)
│   ├── spec_folder.py          ← finds ../spec/ from anywhere inside the repository
│   ├── command.py              ← the CLI's logic: JSON text in, output + exit code out
│   └── __main__.py             ← the CLI itself: `python -m checkout` (stdin → command.py → stdout)
└── tests/
    ├── conftest.py             ← shared fixtures (`catalog`, `spec_cases`) and Hypothesis settings
    ├── spec_loader.py          ← reads ../spec/catalog.json and ../spec/test-cases/
    ├── support/                ← helpers for tests, not tests themselves
    │   ├── fixtures.py         ← a tiny made-up catalog, a sample receipt, a fixed clock
    │   ├── dollars.py          ← "$12.39" → 1239, without floating point
    │   └── schemas.py          ← schema_errors(schema, value) with jsonschema
    ├── unit/                   ← one small piece of code at a time, with made-up data
    ├── acceptance/             ← the 60 shared spec cases, run against checkout()
    ├── integration/            ← real things: the spec files, a temp folder, an SQLite file
    ├── bdd/                    ← pytest-bdd step definitions for ../spec/features/
    ├── property/               ← P1 to P14 from spec/properties.md with Hypothesis
    ├── contract/               ← receipts and CLI output checked against spec/schema/
    ├── snapshot/               ← whole receipts compared with saved snapshots (__snapshots__/)
    └── e2e/                    ← the real program run as a separate process
```

Python conventions worth knowing:

* **pytest finds tests by name:** files called `test_*.py`, functions called `test_*`.
* **`conftest.py`** holds fixtures that every test below it can use, with no import needed.
* **src layout:** the code lives in `src/`, so tests use the *installed* package
  (`pip install -e .`), just like a real user would. This catches packaging mistakes early.
* The spec files are **read from the shared `spec/` folder**, never copied – by the tests
  (`tests/spec_loader.py`) and by the program (`checkout.find_spec_folder()`, which walks up
  the folder tree, so it also works from mutmut's copy in `mutants/`).

## The kinds of tests, and where to look

| Kind of test | Question it answers | Look at |
|---|---|---|
| **Unit** | Does this one function do the right thing? | `tests/unit/test_error_codes.py` is the gentlest first read, then `test_money.py` |
| **Parameterised** | Does it work for all these examples? | `@pytest.mark.parametrize(..., ids=...)` tables in `tests/unit/test_promotions.py`, `test_tax.py`, `test_money.py` |
| **Shared spec cases** | Does it match the 60 cases every language must pass? | `tests/acceptance/test_spec_cases.py` |
| **BDD** | Does it behave as the Given/When/Then scenarios say? | `tests/bdd/test_features.py` – each sentence from `spec/features/README.md` appears once |
| **Test doubles** | How do I test without the real clock? Was `save()` really never called? | `tests/unit/test_orders.py`: a fake (`InMemoryOrderRepository`), a stub clock, `Mock` / `MagicMock` / `create_autospec` with `assert_not_called()`, and `monkeypatch` |
| **Integration** | Do the repositories really write files and database rows? | `tests/integration/test_json_file_repository.py`, `test_sqlite_repository.py` (with `tmp_path`) |
| **Property-based** | Is it true for *every* basket, not just my examples? | `tests/property/test_properties.py` – `test_p1_...` to `test_p14_...` |
| **Contract** | Is the output even the right shape? | `tests/contract/test_schemas.py` |
| **Snapshot** | Did the output change since last time? | `tests/snapshot/test_receipt.py` and its `__snapshots__/` |
| **End-to-end** | Does the real program work from the outside? | `tests/e2e/test_cli.py`, then `../e2e/run.sh python -m checkout` |
| **Doctests** | Do the examples in the documentation still work? | The `>>>` examples in `src/checkout/money.py`, `tax.py`, `coupons.py`, … run with `--doctest-modules` |
| **Coverage** | Which lines and branches did the tests never run? | `pytest --cov` – `fail_under = 100` in `pyproject.toml` |
| **Mutation** | Would the tests notice if the code were wrong? | `mutmut run` – see below |
| **Static checks** | Can we catch mistakes without running anything? | `mypy` (strict), `ruff check .`, `ruff format --check .` |
| **Independence** | Do the tests pass in any order, side by side? | pytest-randomly shuffles every run; `pytest -n auto` runs them in parallel |

### Mutation testing, and what it found

mutmut makes ~550 small changes to `src/checkout/` ("mutants": `>=` becomes `>`, `+` becomes
`-`, `"standard"` becomes `"XXstandardXX"`, an argument disappears) and runs the tests against
each one. A mutant that *survives* means no test noticed the change. The score is currently
**100% (548 of 548 killed)**, and CI fails below 95%. It takes about 30 seconds.

The first run killed 562 of 577 – the third commit in the history shows what the 15 survivors
taught us:

* **Real gaps in the tests.** `format_timestamp` passed every test while converting to the
  computer's *local* time instead of UTC – invisible on a machine set to UTC. A test now runs
  it in the `JST-9` time zone. Two more: the exact indentation of the saved JSON files, and
  `main()` quietly ignoring the catalog it was given.
* **Code that hid a change.** The CLI's error output had its own `json.dumps` for a dict with
  *one* key, so a broken item separator could never show. Both outputs now share
  `compact_json()`.
* **Equivalent mutants** – changes nobody could ever observe: SQL keywords in lower case
  (SQL ignores case), `"UTF-8"` instead of `"utf-8"`. Reading files as bytes removed some;
  the rest are marked `# pragma: no mutate`, with a comment saying why.

Two more things worth knowing:

* Boundaries are tested on both sides (`2499` / `2500`, `4999` / `5000`, `9` / `10`,
  `99` / `100`), because `>` vs `>=` is the classic surviving mutant. Error **messages** are
  tested, not only error codes.
* `src/checkout/__main__.py` is excluded from coverage and mutation: it runs in a separate
  process (see `tests/e2e/`), which in-process tools cannot observe. Its logic is in
  `command.py`, which *is* covered and mutated. The e2e tests are not part of mutmut's test
  run for the same reason (`[tool.mutmut]` in `pyproject.toml`).

## How to add a new test

Say you want to check that `round_half_up` rounds 10% of $0.05 (0.5 cents) up to 1 cent.

1. **Pick the kind of test.** It is one function with made-up numbers: a unit test, so open
   `tests/unit/test_money.py` (or create a new file whose name starts with `test_`).
2. **Write a function** whose name starts with `test_` and says what should happen,
   following the *arrange / act / assert* pattern:

   ```python
   def test_rounds_half_a_cent_up_to_one_cent() -> None:
       # Arrange
       amount = 5

       # Act
       rounded = round_half_up(amount, 10)

       # Assert
       assert rounded == 1
   ```

3. **Run it:** `pytest -k half_a_cent`.
4. **Make sure it can fail.** Change `1` to `0` and check that the test goes red.
   A test that never fails does not test anything! (Mutation testing automates this idea.)
5. **Run everything** before you open a pull request (the block under
   *Running the tests and checks*), and `mutmut run` if you changed `src/`.

Want the same test for many inputs? Use
[`@pytest.mark.parametrize`](https://docs.pytest.org/en/stable/how-to/parametrize.html) with
`ids=` – see `tests/unit/test_tax.py`. Want to add a case for *every* language instead? Add it
to [`spec/test-cases/`](../spec/test-cases/) and it is picked up automatically by the
acceptance, contract, e2e and shared-runner tests (and update the count in
`tests/integration/test_spec_loader.py`).

## Things to know about this implementation

* **Whole numbers only.** `round_half_up(amount, percent)` is `(amount * percent + 50) // 100`,
  as the hint in rules §7 suggests. Python's built-in `round()` would be wrong here: it uses
  banker's rounding (`round(22.5) == 22`).
* **Frozen dataclasses.** `Basket`, `Receipt` and friends cannot be changed after they are
  created. `Receipt.to_dict()` (via `dataclasses.asdict`) gives the snake_case JSON of
  `receipt.schema.json`; `Receipt.from_dict()` reads it back.
* **Parse, don't validate.** `parse_shipping_method` turns the basket's free-form `shipping`
  string into the type `Literal["standard", "express"]` or raises; after that mypy knows the
  value is fine.
* **The clock is a function.** `Clock = Callable[[], datetime]`. Production uses
  `system_clock`; tests pass `fixed_clock` (always 10:30 UTC on 15 January 2026), a `Mock` to
  count the calls, or `monkeypatch` to replace `system_clock` itself.
* **Repositories are a `Protocol`** (`save` and `find`), with three implementations. The
  in-memory one doubles as the test fake. Ids are `ORDER-<n>`; the JSON-file store takes the
  next number after the highest existing file, the SQLite store uses the row id.
  `SqliteOrderRepository` works in a `with` block, which closes the database file.
* **`place_order` = checkout + save.** It cannot save anything when checkout raises, which is
  what the mock tests and property P14 check.
* **The CLI has three exit codes:** `0` with the receipt, `1` with `{"error":"<CODE>"}`, and
  `2` with nothing on stdout and a message on stderr when the input is not a basket at all
  (not JSON, no `items` list, a `quantity` that is not a whole number, …).

## Ideas for the next contribution

Every Tier 1, Tier 2 and Tier 3 row of the Python column is done, so what is left is polish:

* Compare with another language once it is implemented: same spec, same kinds of tests,
  different idioms. That comparison is the point of the project.
* A second BDD runner, [behave](https://behave.readthedocs.io/), on the same feature files
  would show the trade-off against pytest-bdd.
* [pyright](https://microsoft.github.io/pyright/) next to mypy, to compare two type checkers
  on the same strict code.
* The CLI could validate its stdin against `basket.schema.json` with jsonschema instead of the
  hand-written shape checks in `command.py`.
