# End-to-end tests (coming soon)

An *end-to-end* (e2e) test checks the whole program from the outside, exactly as a
user would run it, without looking at any of the code inside.

Rules §8 ([`spec/checkout-rules.md`](../../../spec/checkout-rules.md#8-command-line-interface-for-end-to-end-tests))
describes a small command line program: it reads a basket as JSON from standard input,
prints the receipt (or `{ "error": "<CODE>" }`) as JSON, and exits with code `0` or `1`.

Once that program exists (for example `src/cli.ts`, run with `node src/cli.ts`), the
tests in this folder can:

1. start it as a separate process with `spawnSync` from `node:child_process`,
2. write a basket to its standard input,
3. check what it printed and its exit code.

Good first e2e test: run every case from `spec/test-cases/` through the program
(see `tests/support/spec.ts` for loading them).
