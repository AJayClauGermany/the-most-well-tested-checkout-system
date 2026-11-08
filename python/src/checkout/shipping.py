"""Shipping (rules §4.5)."""

from typing import Literal

# A `Literal` type: a value of this type can only be one of these two strings.
ShippingMethod = Literal["standard", "express"]

STANDARD_SHIPPING_COST = 499
EXPRESS_SHIPPING_COST = 999
# Standard shipping is free from this discounted subtotal upwards.
FREE_SHIPPING_FROM = 5000


def parse_shipping_method(shipping: str) -> ShippingMethod:
    """Turns the basket's ``shipping`` field into a known method, or raises UNKNOWN_SHIPPING_METHOD.

    Not implemented yet - see tests/unit/test_shipping.py.
    """
    raise NotImplementedError(f"parse_shipping_method({shipping!r}) is not implemented yet")


def shipping_cost(method: ShippingMethod, discounted_subtotal: int) -> int:
    """The shipping cost in cents for this method and this discounted subtotal.

    Not implemented yet - see tests/unit/test_shipping.py.
    """
    raise NotImplementedError(
        f"shipping_cost({method!r}, {discounted_subtotal}) is not implemented yet"
    )
