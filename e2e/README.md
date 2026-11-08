# Shared end-to-end runner

An **end-to-end (e2e) test** runs the whole program from the outside, the way a user
would, without looking at the code. Here the "user" is this script, and the "program" is
the command line interface that every language provides
([rules §8](../spec/checkout-rules.md#8-command-line-interface-for-end-to-end-tests)):
basket JSON in on standard input, receipt JSON (or `{"error": "<CODE>"}`) out, exit code
`0` or `1`.

Because all six languages speak that same interface, **one** runner can test them all:

```bash
e2e/run.sh <command that starts the checkout CLI>
```

It feeds all 60 shared cases from [`spec/test-cases/`](../spec/test-cases/) to the
program, one at a time, prints `ok` or `FAIL` per case, and exits non-zero if any failed.
It needs only `bash` and [`jq`](https://jqlang.github.io/jq/).

| Language   | Once its CLI exists, run (from the repository root)                          |
|------------|------------------------------------------------------------------------------|
| Python     | `(cd python && ../e2e/run.sh python -m checkout)`                            |
| TypeScript | `(cd typescript && ../e2e/run.sh node src/cli.ts)` – or `npm run test:e2e:shared` |
| Go         | `e2e/run.sh go run ./go/cmd/checkout`                                        |
| Java       | `e2e/run.sh java -jar java/target/checkout.jar`                              |
| C#         | `e2e/run.sh dotnet run --project csharp/src/Checkout.Cli`                    |
| Rust       | `e2e/run.sh cargo run --quiet --manifest-path rust/Cargo.toml --bin checkout` |

The TypeScript command works today (Node.js 22.18+ runs `.ts` files directly); the others
will be confirmed in each language's README when its CLI lands.

## What it checks – and what it leaves to the languages

The runner compares the **six totals** and the **exit code**. It does not check the
receipt `lines` or validate against [`receipt.schema.json`](../spec/schema/receipt.schema.json);
each language's own e2e and contract tests do that with a proper schema library. Keeping
the shared runner small is the point: it is a black-box smoke test that proves "yes, this
language really implements the same checkout".

## Testing the test runner

Test code can have bugs too. [`selftest.sh`](selftest.sh) checks the runner against two
fake programs in [`selftest/`](selftest/):

* `oracle.sh` looks the expected answer up in the spec and prints it – every case must pass.
* `always-empty-basket.sh` is wrong on purpose – the runner must fail it and exit non-zero.

CI runs `selftest.sh` (and `shellcheck`) whenever this folder or the spec changes.

```bash
e2e/selftest.sh
```
