"""Tax (rules §4.6)."""

TAX_PERCENT = 8


def tax_on(discounted_subtotal: int) -> int:
    """8% of ``discounted_subtotal``, rounded half up. Shipping is never passed in here.

    Not implemented yet - see tests/unit/test_tax.py.
    """
    raise NotImplementedError(f"tax_on({discounted_subtotal}) is not implemented yet")
