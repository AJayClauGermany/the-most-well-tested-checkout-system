"""Merging basket items into receipt lines, and pricing a line (rules §4.2)."""

from collections.abc import Sequence

from checkout.models import BasketItem, Product, ReceiptLine
from checkout.promotions import promotion_discount


def merge_items(items: Sequence[BasketItem]) -> list[BasketItem]:
    """One item per distinct SKU, in order of first appearance, quantities added up.

    >>> merge_items([BasketItem("APPLE", 1), BasketItem("PEN", 2), BasketItem("APPLE", 3)])
    [BasketItem(sku='APPLE', quantity=4), BasketItem(sku='PEN', quantity=2)]
    """
    # A dict remembers the order its keys were first added - exactly the order the rules want.
    quantities: dict[str, int] = {}
    for item in items:
        quantities[item.sku] = quantities.get(item.sku, 0) + item.quantity
    return [BasketItem(sku=sku, quantity=quantity) for sku, quantity in quantities.items()]


def price_line(product: Product, quantity: int) -> ReceiptLine:
    """The receipt line for ``quantity`` of ``product``, including its promotion discount."""
    return ReceiptLine(
        sku=product.sku,
        name=product.name,
        quantity=quantity,
        unit_price=product.unit_price,
        line_total=product.unit_price * quantity,
        line_discount=promotion_discount(product.promotion, quantity, product.unit_price),
    )
