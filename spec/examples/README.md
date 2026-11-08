# Example outputs

Small, hand-checked examples of the checkout's **outputs**. They are the worked example
from [rules §5](../checkout-rules.md#5-the-receipt-the-output): 3 bags of coffee and
2 apples with the coupon `SAVE10`.

| File              | What it is                                   | Checked against                     |
|-------------------|----------------------------------------------|-------------------------------------|
| `receipt.json`    | A receipt (rules §5)                         | [`schema/receipt.schema.json`](../schema/receipt.schema.json) |
| `order.json`      | A saved order (rules §6)                     | [`schema/order.schema.json`](../schema/order.schema.json)     |
| `cli-error.json`  | What the CLI prints on failure (rules §8)    | [`schema/cli-error.schema.json`](../schema/cli-error.schema.json) |

They exist so that CI can prove the schemas accept what the rules describe, and so
that tests which only need *a* valid receipt (snapshot tests, schema tests, file
repository tests) have one to start from.
