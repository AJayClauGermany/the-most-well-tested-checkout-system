---
layout: page
title: The checkout rules
---

The full, official rules live in
[`spec/checkout-rules.md`]({{ site.repository_files }}/spec/checkout-rules.md).
This page is the quick tour.

## The idea in one picture

```
basket ──▶ validate ──▶ merge lines ──▶ promotions ──▶ coupon ──▶ shipping ──▶ tax ──▶ receipt
                                                                                          │
                                                                          (optional) save order
```

This is one way to picture how each amount on the receipt depends on the one before
it. The rules define the amounts, not this order of work – a program may compute them
however suits its language. Still, splitting the code along these arrows gives small
pieces that can be tested on their own (**unit tests**), while the whole chain is tested
together (**integration** and **end-to-end tests**).

## The products

| SKU          | Name                | Price  | Promotion                       |
|--------------|---------------------|-------:|---------------------------------|
| `APPLE`      | Apple               |  $0.50 | –                               |
| `BREAD`      | Loaf of Bread       |  $2.25 | –                               |
| `COFFEE`     | Bag of Coffee Beans |  $8.99 | Buy 2, get 1 free               |
| `NOTEBOOK`   | Notebook            |  $3.49 | 10% off when you buy 10 or more |
| `PEN`        | Pen                 |  $1.20 | –                               |
| `HEADPHONES` | Wireless Headphones | $59.99 | –                               |

## The rules in brief

| Step       | Rule |
|------------|------|
| Money      | Always whole **cents** (`899` = $8.99). Never floating point. |
| Validation | No empty baskets, only known products, quantities 1–99. |
| Merging    | The same product on two lines becomes one line. |
| Promotions | Coffee: every 3rd bag free. Notebooks: 10% off 10 or more. |
| Coupons    | `SAVE10` = 10% off. `WELCOME5` = $5 off if you spend at least $25. |
| Shipping   | Standard $4.99 (free from $50). Express $9.99 (never free). |
| Tax        | 8% of the price after discounts. Shipping is not taxed. |
| Rounding   | Half a cent rounds **up** (22.5 → 23). |
| Saving     | Optional. Orders get ids `ORDER-1`, `ORDER-2`, … and a `created_at` time; tests must be able to fix the time. |

## A worked example

3 bags of coffee and 2 apples, with the coupon `SAVE10` and standard shipping:

| Step                         | Working                        | Amount  |
|------------------------------|--------------------------------|--------:|
| Coffee                       | 3 × $8.99                      | $26.97  |
| Apples                       | 2 × $0.50                      |  $1.00  |
| **Subtotal**                 |                                | $27.97  |
| Promotion (1 free coffee)    | − 1 × $8.99                    | −$8.99  |
| After promotions             |                                | $18.98  |
| Coupon `SAVE10`              | 10% of $18.98 = 189.8¢ → 190¢  | −$1.90  |
| Discounted subtotal          |                                | $17.08  |
| Shipping (under $50)         |                                |  $4.99  |
| Tax                          | 8% of $17.08 = 136.64¢ → 137¢  |  $1.37  |
| **Total**                    |                                | **$23.44** |

## Where the tests come from

The same rules are written down as tests in two shared, language-agnostic forms:

* **[Data-driven test cases]({{ site.repository_files }}/spec/test-cases/)** – JSON files:
  a basket goes in, totals (or an error code) come out. Sixty of them, focused on edge cases.
* **[Gherkin features]({{ site.repository_files }}/spec/features/)** – readable
  *Given / When / Then* scenarios for behaviour-driven development (BDD).
* **[Properties]({{ site.repository_files }}/spec/properties.md)** – rules that hold for
  *every* basket, for property-based testing.
* **[JSON Schemas]({{ site.repository_files }}/spec/schema/)** – the exact shape of a basket,
  a receipt and a saved order, for contract tests.

Every language loads these files directly, so all six implementations are held to
exactly the same examples.
