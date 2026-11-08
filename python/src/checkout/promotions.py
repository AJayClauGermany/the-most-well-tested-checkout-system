"""Promotions (rules §4.2): buy X get Y free, and bulk discounts."""

from checkout.models import Promotion


def promotion_discount(promotion: Promotion | None, quantity: int, unit_price: int) -> int:
    """The discount a promotion gives one receipt line, in cents. ``0`` without a promotion.

    Not implemented yet - see tests/unit/test_promotions.py.
    """
    raise NotImplementedError(
        f"promotion_discount({promotion}, {quantity}, {unit_price}) is not implemented yet"
    )
