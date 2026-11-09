"""Promotions (rules §4.2): buy X get Y free, and bulk discounts."""

from checkout.models import BuyXGetYFree, Promotion
from checkout.money import round_half_up


def promotion_discount(promotion: Promotion | None, quantity: int, unit_price: int) -> int:
    """The discount a promotion gives one receipt line, in cents. ``0`` without a promotion.

    >>> promotion_discount(BuyXGetYFree(buy=2, free=1), 3, 899)   # the third coffee is free
    899
    >>> from checkout.models import BulkDiscount
    >>> promotion_discount(BulkDiscount(min_quantity=10, percent=10), 10, 349)  # 10% of $34.90
    349
    """
    if promotion is None:
        return 0

    if isinstance(promotion, BuyXGetYFree):
        # For every complete group of `buy + free` items, `free` of them cost nothing.
        group_size = promotion.buy + promotion.free
        complete_groups = quantity // group_size
        free_items = complete_groups * promotion.free
        return free_items * unit_price

    # Otherwise it is a BulkDiscount: `percent`% off the whole line, but only from
    # `min_quantity` items. (A Promotion is one of exactly two classes, so no third case.)
    if quantity < promotion.min_quantity:
        return 0
    return round_half_up(unit_price * quantity, promotion.percent)
