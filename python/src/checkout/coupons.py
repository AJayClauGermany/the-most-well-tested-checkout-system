"""Coupons (rules §4.4)."""

from collections.abc import Callable
from dataclasses import dataclass

from checkout.errors import CheckoutError, ErrorCode
from checkout.money import round_half_up


@dataclass(frozen=True)
class Coupon:
    """What the shop knows about one coupon code."""

    # The discount in cents, given the price after promotions.
    discount: Callable[[int], int]
    # The least the customer must spend (after promotions) to use it.
    minimum_spend: int


# The coupons that exist, by their normalized code.
COUPONS: dict[str, Coupon] = {
    "SAVE10": Coupon(
        discount=lambda after_promotions: round_half_up(after_promotions, 10), minimum_spend=0
    ),
    "WELCOME5": Coupon(discount=lambda _after_promotions: 500, minimum_spend=2500),
}


def normalize_coupon_code(coupon: str | None) -> str | None:
    """The coupon code as the shop knows it: trimmed and upper case.

    Missing, empty or blank all mean "no coupon" and give ``None``.

    >>> normalize_coupon_code("  save10 ")
    'SAVE10'
    >>> normalize_coupon_code("   ") is None
    True
    """
    if coupon is None:
        return None
    code = coupon.strip().upper()
    return code if code else None


def coupon_discount(coupon: str | None, after_promotions: int) -> int:
    """The coupon discount in cents for ``after_promotions`` cents of shopping.

    Raises ``CheckoutError`` with UNKNOWN_COUPON or COUPON_MINIMUM_NOT_MET
    (rules §4.1 rows 6 and 7).
    """
    code = normalize_coupon_code(coupon)
    if code is None:
        return 0

    known = COUPONS.get(code)
    if known is None:
        raise CheckoutError(ErrorCode.UNKNOWN_COUPON, f"Unknown coupon: {code}")
    if after_promotions < known.minimum_spend:
        raise CheckoutError(
            ErrorCode.COUPON_MINIMUM_NOT_MET,
            f"{code} needs a spend of at least {known.minimum_spend} cents, got {after_promotions}",
        )
    return known.discount(after_promotions)
