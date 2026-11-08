# End-to-end tests (to do)

An **end-to-end (e2e) test** uses the program exactly the way a real user would,
from the outside, with no knowledge of the code inside.

For this project that means the command line interface from
[rules §8](../../../spec/checkout-rules.md#8-command-line-interface-for-end-to-end-tests):

1. Add a CLI, e.g. `src/checkout/__main__.py`, so `python -m checkout` reads a
   basket (JSON) from standard input and prints the receipt (JSON).
2. Add `test_cli.py` here. Use the standard library's
   [`subprocess.run`](https://docs.python.org/3/library/subprocess.html#subprocess.run)
   to start `python -m checkout` (use `sys.executable` for `python`), send a basket on
   standard input, then check the printed JSON and the exit code (`0` on success,
   `1` with `{"error": "<CODE>"}` on failure).

The shared spec cases (`spec/test-cases/*.json`) make good inputs for these tests too.
