"""Merging basket items into receipt lines, and pricing a line (rules §4.2)."""

from collections.abc import Sequence

from checkout.models import BasketItem, Product, ReceiptLine


def merge_items(items: Sequence[BasketItem]) -> list[BasketItem]:
    """One item per distinct SKU, in order of first appearance, quantities added up.

    Not implemented yet - see tests/unit/test_lines.py.
    """
    raise NotImplementedError(f"merge_items({items}) is not implemented yet")


def price_line(product: Product, quantity: int) -> ReceiptLine:
    """The receipt line for ``quantity`` of ``product``, including its promotion discount.

    Not implemented yet - see tests/unit/test_lines.py.
    """
    raise NotImplementedError(f"price_line({product}, {quantity}) is not implemented yet")
