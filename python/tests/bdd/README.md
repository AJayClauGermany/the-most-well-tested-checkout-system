# BDD tests (to do)

**Behaviour-driven development (BDD)** tests are written as plain-English
*Given / When / Then* scenarios that a non-programmer can read and check.

The scenarios already exist in [`spec/features/`](../../../spec/features/) and are
shared by every language. What's missing is the Python *glue code* (the "step
definitions") that turns each sentence into a function call.

Suggested tool: [pytest-bdd](https://pytest-bdd.readthedocs.io/), which runs the
scenarios as normal pytest tests (so fixtures from `tests/conftest.py` just work).
[behave](https://behave.readthedocs.io/) is a popular alternative with its own runner.

Tips:

* Load the `.feature` files straight from `spec/features/` - don't copy them.
  `spec_loader.SPEC_FOLDER` already knows where that folder is.
* The allowed sentences are listed in
  [`spec/features/README.md`](../../../spec/features/README.md#step-vocabulary).
* Turn amounts like `$12.39` into cents (`1239`) with whole-number maths,
  never with `float`.
