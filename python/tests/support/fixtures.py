"""Test helper: small, hand-made test data for unit tests.

Unit tests prefer tiny made-up data over the real catalog: when a test fails you can
work out the right answer in your head. The shared spec cases use the real catalog.
"""

from datetime import UTC, datetime

from checkout import BulkDiscount, BuyXGetYFree, Catalog, Product, Receipt, ReceiptLine

APPLE = Product(sku="APPLE", name="Apple", unit_price=50)
PEN = Product(sku="PEN", name="Pen", unit_price=120)
COFFEE = Product(
    sku="COFFEE",
    name="Bag of Coffee Beans",
    unit_price=899,
    promotion=BuyXGetYFree(buy=2, free=1),
)
NOTEBOOK = Product(
    sku="NOTEBOOK",
    name="Notebook",
    unit_price=349,
    promotion=BulkDiscount(min_quantity=10, percent=10),
)
# A $100.00 item: handy for reaching free shipping and coupon minimums quickly.
TV = Product(sku="TV", name="Television", unit_price=10000)

# Five products, two of them with a promotion.
TINY_CATALOG = Catalog(products=[APPLE, PEN, COFFEE, NOTEBOOK, TV])

# A valid receipt for one apple, used wherever "some receipt" is needed.
ONE_APPLE_RECEIPT = Receipt(
    lines=[
        ReceiptLine(
            sku="APPLE", name="Apple", quantity=1, unit_price=50, line_total=50, line_discount=0
        )
    ],
    subtotal=50,
    promotion_discount=0,
    coupon_discount=0,
    shipping=499,
    tax=4,
    total=553,
)

# 10:30 on 15 January 2026, UTC - the timestamp used throughout the spec's examples.
FIXED_TIME = datetime(2026, 1, 15, 10, 30, 0, tzinfo=UTC)


def fixed_clock() -> datetime:
    """A STUB clock: it always answers the same time, so `created_at` is predictable."""
    return FIXED_TIME
