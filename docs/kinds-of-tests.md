---
layout: page
title: Kinds of tests
---

There is no single "right" way to test. Different kinds of tests answer different
questions, cost different amounts of time, and catch different bugs. This page
describes the kinds of tests this project wants to show. The
[testing matrix](testing-matrix.html) says which tool each language uses for each kind,
and which ones exist so far.

## The testing pyramid

```
            /\
           /  \        End-to-end      few, slow, test everything together
          /----\
         /      \      Integration     some, test parts working together
        /--------\
       /          \    Unit            many, fast, test one small piece
      /____________\
```

Lots of small, fast tests at the bottom; a few big, slow ones at the top.
It is a rule of thumb, not a law – but it is a good place to start.

---

## Unit tests

**Question:** *Does this one small piece of code do the right thing?*

A unit test calls one function (or class) directly and checks the result. No files,
no network, no database. They run in milliseconds, so you can have hundreds.

> *Example:* `round_half_up(225, 10)` returns `23`.

## Data-driven (table / parameterised) tests

**Question:** *Does it work for all these examples?*

The same test code runs once per row of a table of inputs and expected outputs. This
project's shared [`spec/test-cases`]({{ site.repository_files }}/spec/test-cases/) are
exactly this: one table of examples, used by six languages.

> *Example:* for each case in `02-promotions.json`, check the promotion discount.

## Behaviour-driven development (BDD)

**Question:** *Does it behave the way the people who asked for it described?*

Scenarios are written in near-English (*Given / When / Then*) in `.feature` files.
"Step definitions" connect each sentence to code. Non-programmers can read – and even
write – the scenarios.

> *Example:* [`spec/features/coupons.feature`]({{ site.repository_files }}/spec/features/coupons.feature)

## Integration tests

**Question:** *Do these parts work together, including the real outside world?*

Integration tests use real files, databases, or other components instead of
pretending. They are slower and need more setup, so there are fewer of them.

> *Example:* save an order with a file-based repository, then read it back from disk.

## Test doubles: fakes, stubs and mocks

**Question:** *How do I test this piece without its slow or unpredictable neighbours?*

A *test double* stands in for a real dependency. A **fake** is a simple working
version (an in-memory order repository). A **stub** returns canned answers. A **mock**
also records how it was called so the test can check it.

> *Example:* check that a failed checkout never calls `save()` on the repository.

## Property-based tests

**Question:** *Is this always true, for any input?*

Instead of choosing examples yourself, you describe a property and the tool invents
hundreds of random inputs to try to break it – then shrinks any failure down to the
smallest example. The properties this project checks are listed once, for all languages,
in [`spec/properties.md`]({{ site.repository_files }}/spec/properties.md).

> *Example:* property P1 – for **any** valid basket,
> `total = subtotal − discounts + shipping + tax`.

## Contract (schema) tests

**Question:** *Is the output even the right shape?*

A contract says what a piece of data must look like: which fields, which types, which
values are allowed. Here the contract is a JSON Schema,
[`receipt.schema.json`]({{ site.repository_files }}/spec/schema/receipt.schema.json), and a
contract test checks a real receipt against it. Because all six languages share the same
schema, they are held to the same contract – which is what lets one e2e runner test them all.

> *Example:* `shipping` must be exactly `0`, `499` or `999`.

## Fuzzing

**Question:** *Can random garbage crash it?*

A fuzzer throws random (often malformed) input at the program, mutating whatever
made it misbehave last time. It does not know the right answer; it only looks for crashes,
hangs and panics. Property P10 says the checkout must *never* crash: for any input it
returns a receipt or one of six error codes.

## Documentation tests

**Question:** *Do the examples in the documentation still work?*

Some languages let an example in a doc comment run as a test (Python's `doctest`, Go's
`Example` functions, Rust's doctests). The documentation can then never drift away from
the code.

## Snapshot (approval) tests

**Question:** *Did the output change since last time?*

The first run saves the output (for example a printed receipt) to a file. Later runs
compare against it. Great for big outputs that are tedious to assert line by line.

## End-to-end (E2E) tests

**Question:** *Does the whole program work the way a user would use it?*

E2E tests run the real program from the outside. Here, every language will provide the
same command line interface (basket JSON in, receipt JSON out – see
[rules §8]({{ site.repository_files }}/spec/checkout-rules.md#8-command-line-interface-for-end-to-end-tests)),
so a single E2E suite tests all six: [`e2e/run.sh`]({{ site.repository_files }}/e2e/run.sh).

## Mutation testing

**Question:** *Would my tests notice if the code were wrong?*

A mutation testing tool makes small changes ("mutants") to your code – `>=` becomes
`>`, `+` becomes `-` – and runs your tests. If the tests still pass, that mutant
"survived", which shows a gap in your tests.

## Code coverage

**Question:** *Which lines did my tests never run?*

Coverage tools report which lines (and branches) ran during the tests. High coverage
does not prove the tests are good – mutation testing is better at that – but
low coverage definitely shows untested code.

## Benchmarks

**Question:** *How fast is it, and did that change?*

A benchmark runs a piece of code many times and reports the time. The checkout is tiny, so
one benchmark is enough to show the tool; this is not a performance project.

## Architecture tests

**Question:** *Does the code keep to the structure we agreed?*

Rules like "the checkout logic never imports file or console code" can be tested, in
languages with tools for it (ArchUnit for Java, ArchUnitNET for C#). With only a couple of
modules that is nearly overkill – which is why it is only a Tier 2 item in the matrix.

## Static checks

**Question:** *Can we catch mistakes without running anything?*

Type checkers, linters and formatters are not tests, but they catch whole groups of bugs
before any test runs. Every language here runs at least one in CI.

## What we leave out, and why

Browser tests, load tests, security scans, contract tests between services and flaky-test
retries all exist for good reasons – in programs that have screens, traffic, attackers,
other services or non-determinism. The checkout has none of those, so we leave them out
rather than invent a reason. The [testing matrix](testing-matrix.html) keeps the full list.
