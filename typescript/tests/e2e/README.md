# End-to-end tests

An *end-to-end* (e2e) test checks the whole program from the outside, exactly as a
user would run it, without looking at any of the code inside.

Rules §8 ([`spec/checkout-rules.md`](../../../spec/checkout-rules.md#8-command-line-interface-for-end-to-end-tests))
describes a small command line program: it reads a basket as JSON from standard input,
prints the receipt (or `{ "error": "<CODE>" }`) as JSON, and exits with code `0` or `1`.
Ours is [`src/cli.ts`](../../src/cli.ts). Node.js 22.18+ runs it directly:

```bash
echo '{"items":[{"sku":"COFFEE","quantity":3}],"coupon":"SAVE10"}' | node src/cli.ts
```

[`cli.test.ts`](cli.test.ts) does exactly that, from a test:

1. starts the program as a separate process with `spawnSync` from `node:child_process`,
2. writes a basket to its standard input,
3. checks what it printed and its exit code – for every one of the 60 shared cases, plus a
   few things only the program itself can get wrong (one line of output, exit code `2` and a
   message on stderr when the input is not a basket at all).

Because the test only sees stdout, stderr and the exit code, it would pass unchanged for a
program written in any language. That is the idea behind the **shared runner** in
[`e2e/run.sh`](../../../e2e/run.sh) at the root of the repository, which checks all six
languages the same way:

```bash
npm run test:e2e:shared      # = ../e2e/run.sh node src/cli.ts
```

The logic behind the program lives in [`src/command.ts`](../../src/command.ts) and has
ordinary in-process unit tests in [`tests/unit/command.test.ts`](../unit/command.test.ts).
`cli.ts` itself is only the few lines that connect that logic to stdin, stdout and the exit
code – which is why it is left out of the coverage and mutation reports (those only see code
running inside the test process) and tested from the outside here instead.
