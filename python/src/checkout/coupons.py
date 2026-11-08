"""Coupons (rules §4.4)."""


def normalize_coupon_code(coupon: str | None) -> str | None:
    """The coupon code as the shop knows it: trimmed and upper case; ``None`` for no coupon.

    Not implemented yet - see tests/unit/test_coupons.py.
    """
    raise NotImplementedError(f"normalize_coupon_code({coupon!r}) is not implemented yet")


def coupon_discount(coupon: str | None, after_promotions: int) -> int:
    """The coupon discount in cents for ``after_promotions`` cents of shopping.

    Not implemented yet - see tests/unit/test_coupons.py.
    """
    raise NotImplementedError(
        f"coupon_discount({coupon!r}, {after_promotions}) is not implemented yet"
    )
