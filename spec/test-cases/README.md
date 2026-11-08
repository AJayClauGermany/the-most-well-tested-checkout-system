# Shared test cases (data-driven)

Each JSON file holds a list of cases. A case is a **basket** plus **either** the
`expected` totals **or** an `expected_error` code. All amounts are whole cents and every
case uses the standard [`catalog.json`](../catalog.json).

```json
{
  "id": "promo-coffee-three-bags-one-free",
  "description": "3 coffees: the third one is free",
  "basket": {"items": [{"sku": "COFFEE", "quantity": 3}]},
  "expected": {"subtotal": 2697, "promotion_discount": 899, "coupon_discount": 0, "shipping": 499, "tax": 144, "total": 2441}
}
```

```json
{
  "id": "validation-empty-basket",
  "description": "A basket with no items cannot be checked out",
  "basket": {"items": []},
  "expected_error": "EMPTY_BASKET"
}
```

| File                    | Rule area                               |
|-------------------------|-----------------------------------------|
| `01-subtotal.json`      | Line totals, merging duplicate lines    |
| `02-promotions.json`    | Buy-X-get-Y-free and bulk discounts     |
| `03-coupons.json`       | `SAVE10`, `WELCOME5`, unknown coupons   |
| `04-shipping.json`      | Standard / express, free-shipping threshold |
| `05-tax.json`           | 8% tax and rounding                     |
| `06-validation.json`    | Error codes and the order they are checked in |
| `07-full-orders.json`   | Realistic baskets using many rules      |

## Rules for writing cases

* `id` is unique across **all** files, lower-case-with-dashes, and starts with the area
  (`promo-`, `coupon-`, …). Test runners print it, so make it describe the case.
* `description` explains *why* the answer is what it is.
* Prefer **boundary** values: one below, exactly at, and one above a limit.
* Quantities are always integers (strongly typed languages could not even parse `1.5`).
* `coupon` and `shipping` may be left out; they then mean "no coupon" and `"standard"`.
