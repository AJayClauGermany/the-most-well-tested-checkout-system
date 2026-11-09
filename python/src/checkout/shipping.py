"""Shipping (rules §4.5)."""

from typing import Literal

from checkout.errors import CheckoutError, ErrorCode

# A `Literal` type: a value of this type can only be one of these two strings.
ShippingMethod = Literal["standard", "express"]

STANDARD_SHIPPING_COST = 499
EXPRESS_SHIPPING_COST = 999
# Standard shipping is free from this discounted subtotal upwards.
FREE_SHIPPING_FROM = 5000


def parse_shipping_method(shipping: str) -> ShippingMethod:
    """Turns the basket's ``shipping`` field into a known method.

    "Parse, don't validate": after this function the type checker *knows* the value is one
    of the two methods. Raises ``CheckoutError(UNKNOWN_SHIPPING_METHOD)`` for anything else
    (rules §4.1 row 5). Names are matched exactly, so "Express" is unknown.
    """
    if shipping == "standard":
        return "standard"
    if shipping == "express":
        return "express"
    raise CheckoutError(ErrorCode.UNKNOWN_SHIPPING_METHOD, f"Unknown shipping method: {shipping}")


def shipping_cost(method: ShippingMethod, discounted_subtotal: int) -> int:
    """The shipping cost in cents for this method and this discounted subtotal."""
    if method == "express":
        return EXPRESS_SHIPPING_COST
    if discounted_subtotal >= FREE_SHIPPING_FROM:
        return 0
    return STANDARD_SHIPPING_COST
