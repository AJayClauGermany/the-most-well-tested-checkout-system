# Checkout Rules

This is the **single source of truth** for how the checkout system behaves.
Every language implementation follows these rules, and every shared test case
in [`test-cases/`](test-cases/) and [`features/`](features/) is derived from them.

> In one sentence: *given a basket of products, calculate the final price using
> quantities, discounts, shipping, and tax, then optionally save the resulting order.*

The rules are deliberately small. They are not meant to model a real shop.
They are meant to give us **interesting things to test**: arithmetic, rounding,
boundaries, error handling, data loading and storage.

---

## 1. Money is whole cents

All amounts are **integers counting cents**. `199` means `$1.99`.

We never use floating point numbers (`0.1 + 0.2` is not `0.3` in most languages!).
When a percentage produces a fraction of a cent, we round it using
[round half up](#7-rounding).

## 2. The catalog

A catalog is a list of products. Each product has:

| Field        | Example         | Meaning                                    |
|--------------|-----------------|--------------------------------------------|
| `sku`        | `"COFFEE"`      | Unique product code ("stock keeping unit") |
| `name`       | `"Bag of Coffee Beans"` | Human readable name                |
| `unit_price` | `899`           | Price of one item, in cents                |
| `promotion`  | see below       | An optional promotion, or `null`           |

The checkout receives the catalog as an input (it is not hard-coded). This makes it
easy for tests to use a tiny, made-up catalog. The **standard catalog** used by all
shared test cases is [`catalog.json`](catalog.json):

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

## 4. The calculation, step by step

The checkout always runs these steps **in this order**. Each step is small enough
to be unit tested on its own.

### Step 1 – Validate the basket

If anything is wrong, checkout stops and reports **one** error code.
When a basket has several problems, report the **first** one found, checking in
exactly this order:

| # | Check                                                            | Error code                 |
|---|------------------------------------------------------------------|----------------------------|
| 1 | The basket has at least one item                                 | `EMPTY_BASKET`             |
| 2 | Every item's `sku` exists in the catalog (checked top to bottom) | `UNKNOWN_PRODUCT`          |
| 3 | Every item's `quantity` is a whole number from 1 to 99 (top to bottom) | `INVALID_QUANTITY`   |
| 4 | After merging lines (step 2), each product's quantity is at most 99 | `INVALID_QUANTITY`      |
| 5 | The shipping method is `standard` or `express`                   | `UNKNOWN_SHIPPING_METHOD`  |
| 6 | The coupon (if any) is a known coupon                            | `UNKNOWN_COUPON`           |
| 7 | The coupon's minimum spend (if any) is met                       | `COUPON_MINIMUM_NOT_MET`   |

Check 7 needs the price after promotions, so in practice it happens during step 4.
That is fine: the *order* of the checks above is what matters.

SKUs are matched **exactly** – they are case sensitive, so `"apple"` is an unknown product.

### Step 2 – Merge duplicate lines

If the same SKU appears more than once, the lines are combined into a single line
whose quantity is the sum. The merged line keeps the position of the **first**
appearance.

`[APPLE x1, PEN x2, APPLE x3]` becomes `[APPLE x4, PEN x2]`.

### Step 3 – Price each line and apply promotions

For each line:

```
line_total = unit_price × quantity
```

then the product's promotion (if any) gives a **line discount**:

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

Then add everything up:

```
subtotal           = sum of all line_total
promotion_discount = sum of all line_discount
```

### Step 4 – Apply the coupon

At most one coupon can be used. Coupon codes are **not case sensitive** and
spaces around them are ignored (`" save10 "` works the same as `"SAVE10"`).
An empty or blank coupon (`""` or `"   "`) counts as no coupon.

The coupon works on the amount **after promotions**:

```
after_promotions = subtotal − promotion_discount
```

| Code       | Effect                     | Minimum spend (`after_promotions`) |
|------------|----------------------------|------------------------------------|
| `SAVE10`   | 10% off `after_promotions` | none                               |
| `WELCOME5` | $5.00 off                  | $25.00 (`2500`)                    |

```
SAVE10:   coupon_discount = round_half_up(after_promotions × 10 ÷ 100)
WELCOME5: coupon_discount = 500
no coupon: coupon_discount = 0
```

If `WELCOME5` is used and `after_promotions` is below `2500`, checkout fails with
`COUPON_MINIMUM_NOT_MET`.

```
discounted_subtotal = after_promotions − coupon_discount
```

### Step 5 – Shipping

| Method     | Cost   | Free when…                          |
|------------|--------|-------------------------------------|
| `standard` | $4.99  | `discounted_subtotal` ≥ $50.00 (`5000`) |
| `express`  | $9.99  | never                               |

Shipping method names are matched exactly (lower case).

### Step 6 – Tax

Tax is **8%** of the discounted subtotal. Shipping is **not** taxed.

```
tax = round_half_up(discounted_subtotal × 8 ÷ 100)
```

### Step 7 – Total

```
total = discounted_subtotal + shipping + tax
```

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

Lines appear in the order of first appearance in the basket (after merging).

The exact shape of a receipt is written down as a JSON Schema in
[`schema/receipt.schema.json`](schema/receipt.schema.json). It is the **contract** every
language's output must honour, and tests can check real receipts against it.

## 6. Saving an order (optional)

After a successful checkout the caller may choose to **save** the receipt as an order.

* Saving is done through an **order repository** with two operations:
  * `save(receipt)` → returns the id of the new order
  * `find(order_id)` → returns the saved order, or "not found"
* Order ids are strings. The first saved order is `"ORDER-1"`, the next `"ORDER-2"`, and so on
  (per repository).
* A saved order remembers **when** it was saved, as `created_at`: the current time in UTC,
  written as ISO 8601 with whole seconds and a trailing `Z`, e.g. `"2026-01-15T10:30:00Z"`.
* A failed checkout never saves anything.

A saved order therefore looks like this (full shape in
[`schema/order.schema.json`](schema/order.schema.json)):

```json
{ "id": "ORDER-1", "created_at": "2026-01-15T10:30:00Z", "receipt": { "...": "..." } }
```

### The clock

The repository gets the time from a **clock** that is handed to it when it is created
(an argument, a constructor parameter, an interface – whatever is natural in the language).
The real program passes the system clock. Tests pass a **fixed clock** that always answers
the same time, so `created_at` is predictable.

The clock is the only part of the checkout whose answer changes from one run to the next.
That makes it the textbook reason to use a *test double*: a test that compares `created_at`
against the real time would be flaky, so we replace the clock instead.

### Kinds of repository

Every language starts with an **in-memory** repository (a fake that is also useful in
production code for demos). To show **integration testing**, two more are planned, each
storing orders in the shape above:

* a **JSON file** repository (one file per repository or per order – your choice), tested
  against a real temporary folder;
* an **SQLite** repository, using the language's usual SQLite library, tested against a
  temporary database file.

A real database server behind Docker is deliberately **not** a goal: everything here must
run on a learner's laptop with nothing but the language's own tools installed.

## 7. Rounding

`round_half_up` rounds to the nearest whole cent, and exact halves go **up**:

| Exact value | Rounded |
|-------------|---------|
| 22.4        | 22      |
| 22.5        | 23      |
| 22.6        | 23      |

Because every amount is a non-negative integer, this can be done with whole-number
maths only: `round_half_up(amount × percent ÷ 100) = (amount × percent + 50) ÷ 100`
using integer (floor) division.

⚠️ Watch out: some languages' built-in `round()` uses *banker's rounding*
(halves go to the nearest even number, so `22.5 → 22`). That is a classic bug,
and the shared test cases include examples that catch it.

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
