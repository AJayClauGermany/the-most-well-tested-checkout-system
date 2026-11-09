# End-to-end tests

An *end-to-end* (e2e) test checks the whole program from the outside, exactly as a
user would run it, without looking at any of the code inside.

Rules §8 ([`spec/checkout-rules.md`](../../../spec/checkout-rules.md#8-command-line-interface-for-end-to-end-tests))
describes a small command line program: it reads a basket as JSON from standard input,
prints the receipt (or `{ "error": "<CODE>" }`) as JSON, and exits with code `0` or `1`.
Ours is [`src/checkout/__main__.py`](../../src/checkout/__main__.py), so it runs as
`python -m checkout`:

```bash
echo '{"items":[{"sku":"COFFEE","quantity":3}],"coupon":"SAVE10"}' | python -m checkout
```

[`test_cli.py`](test_cli.py) does exactly that, from a test:

1. starts the program as a separate process with `subprocess.run` (using `sys.executable`,
   the Python that runs the tests, so the same virtual environment is used),
2. writes a basket to its standard input,
3. checks what it printed and its exit code – for every one of the 60 shared cases, plus a
   few things only the program itself can get wrong (one line of output, exit code `2` and a
   message on stderr when the input is not a basket at all).

Because the test only sees stdout, stderr and the exit code, it would pass unchanged for a
program written in any language. That is the idea behind the **shared runner** in
[`e2e/run.sh`](../../../e2e/run.sh) at the root of the repository, which checks all six
languages the same way:

```bash
cd python && ../e2e/run.sh python -m checkout      # with the .venv active
```

The logic behind the program lives in [`src/checkout/command.py`](../../src/checkout/command.py)
and has ordinary in-process unit tests in [`tests/unit/test_command.py`](../unit/test_command.py)
(with `monkeypatch` standing in for stdin and `capsys` capturing stdout). `__main__.py` itself
is only the few lines that connect that logic to the terminal – which is why it is left out of
the coverage and mutation reports (those only see code running inside the test process) and
tested from the outside here instead.
