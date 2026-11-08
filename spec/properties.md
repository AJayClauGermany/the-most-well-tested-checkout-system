# Properties: things that are always true

The [test cases](test-cases/) are **examples**: *this* basket gives *that* total.
A **property** is different: it is a rule that must hold for **every** basket.
Property-based testing tools (Hypothesis, fast-check, rapid, jqwik, FsCheck, proptest, …)
generate hundreds of random baskets and try to find one that breaks a property. When they
do, they shrink it down to the smallest failing example.

Just like the test cases and the Gherkin features, this list is shared, so that all six
languages check the **same** properties. Each property has an id (`P1`, `P2`, …); use it in
the test name so the suites can be compared.

All amounts are whole cents. "Valid basket" means a basket that checks out successfully.
`save` and `find` are the vocabulary from rules §6 for storing and looking up an order,
not required method names.

## Generating baskets

A good basket generator produces:

* 1 to 8 items, each with a SKU from the [standard catalog](catalog.json) and a quantity
  from 1 to 99. The same SKU may appear more than once (that is how merging gets tested).
* An optional coupon: missing, `null`, `"SAVE10"`, `"WELCOME5"`, the same codes with
  different case or surrounding spaces, or a random string.
* An optional shipping method: missing, `"standard"`, `"express"`, or a random string.

Baskets that fail validation are **useful**, not a problem: P1–P9 only apply when checkout
succeeds (skip or filter the rest), and P10 is about the failures.

## Properties of a successful checkout

| Id  | Property | In words |
|-----|----------|----------|
| P1  | `total == subtotal − promotion_discount − coupon_discount + shipping + tax` | The totals add up. |
| P2  | every amount on the receipt `≥ 0`; `promotion_discount ≤ subtotal`; `coupon_discount ≤ subtotal − promotion_discount` | Nothing is ever negative and a discount never exceeds what it discounts. |
| P3  | `subtotal == sum(line_total)` and `promotion_discount == sum(line_discount)` and every `line_total == unit_price × quantity` | The receipt lines explain the totals. |
| P4  | shuffling the items gives the same totals | Order of items never changes the price (line order may differ). |
| P5  | splitting one line into two lines of the same SKU (quantities adding up to the original) gives the **same receipt** | Merging lines is invisible to the customer. |
| P6  | `shipping` is `0`, `499` or `999`; it is `999` if and only if shipping is `express`; with `standard` it is `0` if and only if `subtotal − promotion_discount − coupon_discount ≥ 5000` | The shipping table is applied exactly. |
| P7  | `tax == (discounted_subtotal × 8 + 50) ÷ 100` (integer division) | Tax is 8% rounded half up, and shipping is never taxed. |
| P8  | writing the coupon in another case or with surrounding spaces gives the same receipt | Coupons are case-insensitive and trimmed. |
| P9  | the receipt has one line per **distinct** SKU, in order of first appearance, each with quantity 1–99 | Lines are merged and ordered as rules §4.2 says. |

## Properties of any input at all

| Id  | Property | In words |
|-----|----------|----------|
| P10 | for **any** basket (even nonsense), checkout returns **either** a receipt that satisfies [`schema/receipt.schema.json`](schema/receipt.schema.json) **or** exactly one of the six [error codes](checkout-rules.md#9-error-codes) – it never crashes | Validation is total. This is also the target for fuzzing. |
| P11 | `round_half_up(amount × percent ÷ 100)` is within half a cent of the exact value, and equals it when the exact value is whole | Rounding is correct for every input, not just the examples. |

## Properties of saving orders

| Id  | Property | In words |
|-----|----------|----------|
| P12 | after saving `n` orders the ids are exactly `ORDER-1` … `ORDER-n`, in that order | Ids are sequential per repository. |
| P13 | `find(save(receipt))` returns an order whose `receipt` equals the saved one and whose `created_at` is what the clock said | Nothing is lost on the way in and out. |
| P14 | a failed checkout never saves anything, and `find` of an id that was never returned by `save` is "not found" | Failures leave no trace; lookups do not invent orders. |

## Adding a property

Add a row here with the next free id, then add the test in every language you can
(it is fine to start with one). A property must follow from
[`checkout-rules.md`](checkout-rules.md); if it does not, change the rules first.
