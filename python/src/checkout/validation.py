"""Basket validation (rules §4.1, rows 1 to 4).

Rows 5 to 7 (shipping and coupon) are checked by shipping.py and coupons.py, because
the coupon check needs the price after promotions.
"""

from checkout.models import Basket, Catalog

MIN_QUANTITY = 1
MAX_QUANTITY = 99


def is_valid_quantity(quantity: object) -> bool:
    """True when ``quantity`` is a whole number from 1 to 99 (rules §4.1 row 3).

    Not implemented yet - see tests/unit/test_validation.py.
    """
    raise NotImplementedError(f"is_valid_quantity({quantity!r}) is not implemented yet")


def validate_items(basket: Basket, catalog: Catalog) -> None:
    """Raises a CheckoutError for rows 1-4 of rules §4.1, in that order.

    Not implemented yet - see tests/unit/test_validation.py.
    """
    raise NotImplementedError(f"validate_items({basket}, {catalog}) is not implemented yet")
