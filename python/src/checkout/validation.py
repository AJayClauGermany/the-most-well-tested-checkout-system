"""Basket validation (rules §4.1, rows 1 to 4).

Rows 5 to 7 (shipping and coupon) are checked by shipping.py and coupons.py, because
the coupon check needs the price after promotions.
"""

from checkout.catalog import require_product
from checkout.errors import CheckoutError, ErrorCode
from checkout.lines import merge_items
from checkout.models import Basket, Catalog

MIN_QUANTITY = 1
MAX_QUANTITY = 99


def is_valid_quantity(quantity: object) -> bool:
    """True when ``quantity`` is a whole number from 1 to 99 (rules §4.1 row 3).

    Anything that is not an ``int`` is not a quantity: not ``3.0``, not ``"3"``, and not
    ``True`` either (Python counts ``True`` as ``1`` in arithmetic, but ``bool`` is a
    separate type - that is why the extra check).

    >>> is_valid_quantity(1), is_valid_quantity(99)
    (True, True)
    >>> is_valid_quantity(0), is_valid_quantity(100), is_valid_quantity(2.5)
    (False, False, False)
    """
    if not isinstance(quantity, int) or isinstance(quantity, bool):
        return False
    return MIN_QUANTITY <= quantity <= MAX_QUANTITY


def validate_items(basket: Basket, catalog: Catalog) -> None:
    """Raises a CheckoutError for rows 1-4 of rules §4.1, in that order.

    Does nothing when the items are fine.
    """
    # Row 1: no items at all.
    if not basket.items:
        raise CheckoutError(ErrorCode.EMPTY_BASKET, "The basket has no items")

    # Row 2: every SKU must exist. All SKUs are checked before any quantity, so an unknown
    # product on line 3 still wins over a bad quantity on line 1.
    for item in basket.items:
        require_product(catalog, item.sku)

    # Row 3: every quantity must be a whole number from 1 to 99.
    for item in basket.items:
        if not is_valid_quantity(item.quantity):
            raise CheckoutError(
                ErrorCode.INVALID_QUANTITY, f"Invalid quantity for {item.sku}: {item.quantity}"
            )

    # Row 4: the lines of one SKU, added together, may not go over 99 either.
    for line in merge_items(basket.items):
        if line.quantity > MAX_QUANTITY:
            raise CheckoutError(
                ErrorCode.INVALID_QUANTITY, f"Too many {line.sku} in total: {line.quantity}"
            )
