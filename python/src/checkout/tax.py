"""Tax (rules §4.6)."""

from checkout.money import round_half_up

TAX_PERCENT = 8


def tax_on(discounted_subtotal: int) -> int:
    """8% of ``discounted_subtotal``, rounded half up. Shipping is never passed in here.

    >>> tax_on(50)    # 8% of $0.50 is exactly 4 cents
    4
    >>> tax_on(120)   # 9.6 cents rounds up to 10
    10
    """
    return round_half_up(discounted_subtotal, TAX_PERCENT)
