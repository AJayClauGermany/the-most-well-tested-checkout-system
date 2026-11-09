---
layout: home
title: Home
---

**Given a basket of products, calculate the final price using quantities, discounts,
shipping and tax, then optionally save the resulting order.**

That is the whole program. It is small on purpose, so that we can spend our energy on
the interesting part: **testing it in as many ways as possible**, in six languages:

| Python | TypeScript | Go | Java | C# | Rust |
|:------:|:----------:|:--:|:----:|:--:|:----:|

## Who is this for?

Anyone getting started with programming who has heard "you should write tests" and
wonders *what that actually looks like*. The code is written to be **read**, not to run a
real shop. If you can write a function and a loop, you can follow along.

## What you will find here

* **[Getting started](getting-started.html)** – clone the repository and run your first tests.
* **[The checkout rules](checkout-rules.html)** – what the program must do. Every language follows the same rules.
* **[Kinds of tests](kinds-of-tests.html)** – unit, integration, BDD, property-based, end-to-end … what they are and when to use them.
* **[Languages](languages.html)** – how each language sets up its tests, side by side.
* **[Contributing](contributing.html)** – add a test, a tool, or an explanation.

## Tests first

The rules and the shared test cases came first; the checkout came after – in all six
languages. Each language folder's git history shows the commit with the whole test suite
(red) before the commit that makes it pass (green), followed by the commits where coverage
and mutation testing found gaps in the tests. Write a failing test, make it pass, tidy up,
repeat.
