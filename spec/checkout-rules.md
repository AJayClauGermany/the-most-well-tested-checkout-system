# Checkout Rules

This is the **single source of truth** for how the checkout system behaves.
Every language implementation follows these rules, and every shared test case
in [`test-cases/`](test-cases/), [`features/`](features/) and
[`properties.md`](properties.md) is derived from them.

> In one sentence: *given a basket of products, calculate the final price using
> quantities, discounts, shipping, and tax, then optionally save the resulting order.*

The rules are deliberately small. They are not meant to model a real shop.
They are meant to give us **interesting things to test**: arithmetic, rounding,
boundaries, error handling, data loading and storage.

## How to read these rules

These rules describe **what** the checkout does – what goes in, what comes out – and
never **how** a program should be written. That is on purpose: six languages implement
them, and each should do it in the way that is natural for that language. The tests can
only see inputs, outputs and error codes, so that is all the rules talk about.

* **Formulas define values, not an order of computation.** `tax = …` says what the tax
  *is*. Whether a program works it out first, last or lazily is its own business.
* **Named operations are vocabulary.** When the rules say `save` or `find`, they name a
  capability so the properties and scenarios can refer to it. A language may call it
  `Save`, `store`, `persist` or anything else idiomatic, and shape it as a function, a
  method, a trait or an interface.
* **Shapes are described as JSON** because that is what the shared test files and the
  command line interface use. Inside a program, a basket or a receipt is whatever type
  is natural: a dictionary, a record, a struct, a class.
* **Hints are not rules.** Paragraphs marked 💡 *Hint* show one way to get something
  right. You may ignore them if your tests pass.

If a sentence here can only be satisfied by writing the code one particular way, it is a
bug in the rules – please open an issue.

---

## 1. Money is whole cents

All amounts are **integers counting cents**. `199` means `$1.99`. Every amount on a
receipt is a non-negative whole number. When a percentage produces a fraction of a cent,
the result is rounded using [round half up](#7-rounding).

⚠️ Floating point numbers are the classic way to get this wrong: `0.1 + 0.2` is not
`0.3` in most languages. The shared test cases include amounts that expose this, so
whatever a program does internally, its cents must come out exact.

## 2. The catalog

A catalog is a list of products. Each product has:

| Field        | Example         | Meaning                                    |
|--------------|-----------------|--------------------------------------------|
| `sku`        | `"COFFEE"`      | Unique product code ("stock keeping unit") |
| `name`       | `"Bag of Coffee Beans"` | Human readable name                |
| `unit_price` | `899`           | Price of one item, in cents                |
| `promotion`  | see below       | An optional promotion, or `null`           |

The catalog is an **input** to the checkout, so tests may supply any catalog they like,
including a tiny made-up one. The **standard catalog** used by all shared test cases is
[`catalog.json`](catalog.json):

| SKU          | Name                | Price  | Promotion                          |
|--------------|---------------------|-------:|------------------------------------|
| `APPLE`      | Apple               |  $0.50 | –                                  |
| `BREAD`      | Loaf of Bread       |  $2.25 | –                                  |
| `COFFEE`     | Bag of Coffee Beans |  $8.99 | Buy 2, get 1 free                  |
| `NOTEBOOK`   | Notebook            |  $3.49 | 10% off when you buy 10 or more    |
| `PEN`        | Pen                 |  $1.20 | –                                  |
| `HEADPHONES` | Wireless Headphones | $59.99 | –                                  |

## 3. The basket (the input)

```json
{
  "items": [
    { "sku": "COFFEE", "quantity": 3 },
    { "sku": "APPLE", "quantity": 2 }
  ],
  "coupon": "SAVE10",
  "shipping": "standard"
}
```

| Field      | Required? | Meaning                                                   |
|------------|-----------|-----------------------------------------------------------|
| `items`    | yes       | A list of `{ sku, quantity }` lines                        |
| `coupon`   | no        | A coupon code, or `null` / missing for no coupon          |
| `shipping` | no        | `"standard"` or `"express"`. Missing means `"standard"`   |

This JSON is the form used by the test files and the command line interface
([`schema/basket.schema.json`](schema/basket.schema.json)). Inside a program the basket
can be any natural type.

## 4. The calculation

A checkout either **fails** with exactly one [error code](#9-error-codes) or **succeeds**
with a [receipt](#5-the-receipt-the-output). This section defines when it fails, and what
every amount on the receipt is when it succeeds.

### 4.1 When checkout fails

A basket is **invalid**, and checkout fails, if any of the following is true. When more
than one is true, the error code reported is the **first in this table** that applies.

| # | The basket is invalid when…                                                    | Error code                 |
|---|--------------------------------------------------------------------------------|----------------------------|
| 1 | it has no items                                                                | `EMPTY_BASKET`             |
| 2 | an item's `sku` is not in the catalog (the earliest such item decides)         | `UNKNOWN_PRODUCT`          |
| 3 | an item's `quantity` is not a whole number from 1 to 99                        | `INVALID_QUANTITY`         |
| 4 | the quantities of one SKU across all its lines add up to more than 99          | `INVALID_QUANTITY`         |
| 5 | the shipping method is not `standard` or `express`                             | `UNKNOWN_SHIPPING_METHOD`  |
| 6 | the coupon (if any) is not a known coupon                                      | `UNKNOWN_COUPON`           |
| 7 | the coupon's minimum spend (if any) is not met – see [4.4](#44-the-coupon)     | `COUPON_MINIMUM_NOT_MET`   |

Rows 2 and 3 look at items in basket order, so a basket with an unknown product on line 1
and a bad quantity on line 3 reports `UNKNOWN_PRODUCT`, while an unknown product on line 3
and a bad quantity on line 1 *also* reports `UNKNOWN_PRODUCT` – row 2 wins over row 3
regardless of line position.

SKUs are matched **exactly** – they are case sensitive, so `"apple"` is an unknown product.

### 4.2 Receipt lines

The receipt has **one line per distinct SKU** in the basket, in order of the SKU's first
appearance. A line's `quantity` is the sum of the quantities of every basket item with that
SKU.

`[APPLE x1, PEN x2, APPLE x3]` gives the lines `[APPLE x4, PEN x2]`.

For each line:

```
line_total = unit_price × quantity
```

and the product's promotion (if any) gives the line a **discount**, otherwise
`line_discount = 0`:

* **Buy X, get Y free** (`buy_x_get_y_free`) – for every complete group of
  `buy + free` items, `free` of them cost nothing.

  ```
  free_items    = (quantity ÷ (buy + free), rounded down) × free
  line_discount = free_items × unit_price
  ```

  Coffee is "buy 2, get 1 free": 3 bags → 1 free, 5 bags → 1 free, 6 bags → 2 free.

* **Bulk discount** (`bulk_discount`) – when `quantity ≥ min_quantity`, take
  `percent`% off the line total.

  ```
  line_discount = round_half_up(line_total × percent ÷ 100)   (only if quantity ≥ min_quantity)
  ```

  Notebooks are "10% off 10 or more": 9 notebooks → no discount, 10 → $3.49 off.

### 4.3 Subtotal and promotions

```
subtotal           = sum of all line_total
promotion_discount = sum of all line_discount
after_promotions   = subtotal − promotion_discount
```

### 4.4 The coupon

At most one coupon applies. Coupon codes are **not case sensitive** and spaces around
them are ignored (`" save10 "` means the same as `"SAVE10"`). An empty or blank coupon
(`""` or `"   "`) means no coupon.

| Code       | `coupon_discount`                             | Minimum spend (`after_promotions`) |
|------------|-----------------------------------------------|------------------------------------|
| `SAVE10`   | `round_half_up(after_promotions × 10 ÷ 100)`  | none                               |
| `WELCOME5` | `500`                                         | `2500` ($25.00)                    |
| no coupon  | `0`                                           |                                    |

A coupon whose minimum spend is not met makes the basket invalid with
`COUPON_MINIMUM_NOT_MET` (row 7 in [4.1](#41-when-checkout-fails)).

```
discounted_subtotal = after_promotions − coupon_discount
```

### 4.5 Shipping

| Method     | `shipping` | Free (`0`) when…                       |
|------------|------------|----------------------------------------|
| `standard` | `499`      | `discounted_subtotal ≥ 5000` ($50.00)   |
| `express`  | `999`      | never                                  |

Shipping method names are matched exactly (lower case).

### 4.6 Tax

Tax is **8%** of the discounted subtotal. Shipping is **not** taxed.

```
tax = round_half_up(discounted_subtotal × 8 ÷ 100)
```

### 4.7 Total

```
total = discounted_subtotal + shipping + tax
```

💡 *Hint:* the definitions above happen to read well as a pipeline – validate, merge
lines, price, coupon, shipping, tax, total – and splitting the code that way gives you
small, unit-testable pieces. It is a good first design, not a requirement.

## 5. The receipt (the output)

```json
{
  "lines": [
    { "sku": "COFFEE", "name": "Bag of Coffee Beans", "quantity": 3,
      "unit_price": 899, "line_total": 2697, "line_discount": 899 },
    { "sku": "APPLE", "name": "Apple", "quantity": 2,
      "unit_price": 50, "line_total": 100, "line_discount": 0 }
  ],
  "subtotal": 2797,
  "promotion_discount": 899,
  "coupon_discount": 190,
  "shipping": 499,
  "tax": 137,
  "total": 2344
}
```

Lines appear in the order defined in [4.2](#42-receipt-lines).

This JSON is the **contract** every language's output must honour when it is written
out (in the command line interface, in saved orders, in snapshot tests). It is written
down as a JSON Schema in [`schema/receipt.schema.json`](schema/receipt.schema.json).
Inside a program the receipt can be any natural type with the same information.

## 6. Saving an order (optional)

After a successful checkout the caller may choose to **save** the receipt as an order.
Saving has these observable behaviours:

* Saving a receipt produces an **order id**. Ids are strings: the first order saved into
  a given store is `"ORDER-1"`, the next `"ORDER-2"`, and so on.
* A saved order can be **looked up** by its id. The lookup returns the id, the receipt
  exactly as saved, and `created_at`: the time it was saved, in UTC, as ISO 8601 with
  whole seconds and a trailing `Z` (`"2026-01-15T10:30:00Z"`).
* Looking up an id that was never produced by a save is **"not found"** – a normal
  answer, not a crash.
* A **failed** checkout saves nothing.

A saved order, written out, looks like this
([`schema/order.schema.json`](schema/order.schema.json)):

```json
{ "id": "ORDER-1", "created_at": "2026-01-15T10:30:00Z", "receipt": { "...": "..." } }
```

In this document and in [`properties.md`](properties.md) the two operations are called
`save(receipt) → id` and `find(id) → order | not found`, and the scenarios in
[`features/`](features/) call the store "the order repository". Those are names for the
*capabilities*, not required type or method names.

### Time must be controllable

`created_at` is the only thing in the checkout whose value changes from one run to the
next. The rules therefore require that **the source of the current time can be replaced
in tests** with one that always answers the same time. How is up to the language: a
clock argument, a function parameter, an injected interface, a global test hook.

This is the textbook reason for a *test double*: a test that compared `created_at`
against the real time would be flaky, so the time source is replaced instead.

### Where orders are kept

The rules do not say. An order store may keep orders in memory, in a file, in a database
– anything that gives the behaviours above. Which stores each language shows, and why,
is a testing decision and lives in the
[testing matrix](../docs/testing-matrix.md), not here.

## 7. Rounding

`round_half_up` rounds to the nearest whole cent, and exact halves go **up**:

| Exact value | Rounded |
|-------------|---------|
| 22.4        | 22      |
| 22.5        | 23      |
| 22.6        | 23      |

⚠️ Watch out: some languages' built-in `round()` uses *banker's rounding*
(halves go to the nearest even number, so `22.5 → 22`). That is a classic bug,
and the shared test cases include examples that catch it.

💡 *Hint:* because every amount is a non-negative integer, this can be done with
whole-number maths only: `round_half_up(amount × percent ÷ 100) = (amount × percent + 50) ÷ 100`
using integer (floor) division. Any other method that gives the same answers is fine.

## 8. Command line interface (for end-to-end tests)

So that one end-to-end test suite can check **every** language, each implementation
will eventually ship a small command line program that:

1. reads a basket (the JSON from section 3) from **standard input**,
2. on success prints the receipt (section 5) as JSON to **standard output** and exits with code `0`,
3. on failure prints `{ "error": "<ERROR_CODE>" }` to standard output and exits with code `1`.

The program uses the standard catalog ([`catalog.json`](catalog.json)).

On failure the printed JSON follows [`schema/cli-error.schema.json`](schema/cli-error.schema.json);
on success it follows [`schema/receipt.schema.json`](schema/receipt.schema.json).

The [`e2e/`](../e2e/) folder at the root of the repository holds a **shared end-to-end
runner** that feeds every shared test case to any such program and compares the answers,
so one script can check all six languages.

## 9. Error codes

| Code                       | Meaning                                         |
|----------------------------|-------------------------------------------------|
| `EMPTY_BASKET`             | The basket has no items                         |
| `UNKNOWN_PRODUCT`          | An item's SKU is not in the catalog             |
| `INVALID_QUANTITY`         | A quantity is not a whole number from 1 to 99   |
| `UNKNOWN_SHIPPING_METHOD`  | Shipping is not `standard` or `express`         |
| `UNKNOWN_COUPON`           | The coupon code does not exist                  |
| `COUPON_MINIMUM_NOT_MET`   | The coupon needs a bigger spend                 |

Each language reports errors in its own natural way (an exception, an error value,
a `Result`, …) but always uses these exact codes.
