"""The data the checkout works with.

These classes mirror the JSON in the shared ``spec/`` folder. Every amount of
money is an ``int`` counting cents (``199`` means $1.99) - never a ``float``.

``@dataclass`` writes the boring parts of a class for us (``__init__``,
``__repr__`` and ``==``). ``frozen=True`` means an object cannot be changed
after it is created, which makes the code (and the tests) easier to reason about.
"""

from __future__ import annotations

from dataclasses import dataclass, field
from typing import Any

# ---------------------------------------------------------------------------
# Promotions: a product has at most one of these two kinds.
# ---------------------------------------------------------------------------


@dataclass(frozen=True)
class BuyXGetYFree:
    """For every complete group of ``buy + free`` items, ``free`` of them cost nothing.

    JSON: ``{"type": "buy_x_get_y_free", "buy": 2, "free": 1}``
    """

    buy: int
    free: int


@dataclass(frozen=True)
class BulkDiscount:
    """Take ``percent``% off the line when ``quantity >= min_quantity``.

    JSON: ``{"type": "bulk_discount", "min_quantity": 10, "percent": 10}``
    """

    min_quantity: int
    percent: int


# A Promotion is one of the two classes above.
# Later you can tell them apart with isinstance() or a `match` statement.
Promotion = BuyXGetYFree | BulkDiscount


def promotion_from_dict(data: dict[str, Any] | None) -> Promotion | None:
    """Turns the JSON form of a promotion (or ``null``) into a Promotion (or ``None``)."""
    if data is None:
        return None
    if data["type"] == "buy_x_get_y_free":
        return BuyXGetYFree(buy=data["buy"], free=data["free"])
    if data["type"] == "bulk_discount":
        return BulkDiscount(min_quantity=data["min_quantity"], percent=data["percent"])
    raise ValueError(f"Unknown promotion type: {data['type']!r}")


# ---------------------------------------------------------------------------
# The catalog: what the shop sells.
# ---------------------------------------------------------------------------


@dataclass(frozen=True)
class Product:
    sku: str
    name: str
    unit_price: int  # cents
    promotion: Promotion | None = None


@dataclass(frozen=True)
class Catalog:
    products: list[Product]

    def find(self, sku: str) -> Product | None:
        """Returns the product with exactly this SKU (case sensitive), or ``None``."""
        for product in self.products:
            if product.sku == sku:
                return product
        return None

    @classmethod
    def from_dict(cls, data: dict[str, Any]) -> Catalog:
        """Builds a Catalog from the JSON form used by ``spec/catalog.json``."""
        return cls(
            products=[
                Product(
                    sku=product["sku"],
                    name=product["name"],
                    unit_price=product["unit_price"],
                    promotion=promotion_from_dict(product.get("promotion")),
                )
                for product in data["products"]
            ]
        )


# ---------------------------------------------------------------------------
# The basket: the input to checkout.
# ---------------------------------------------------------------------------


@dataclass(frozen=True)
class BasketItem:
    sku: str
    quantity: int


@dataclass(frozen=True)
class Basket:
    items: list[BasketItem] = field(default_factory=list)
    coupon: str | None = None  # None means "no coupon"
    # A plain string (not an Enum) on purpose: a basket may contain an unknown
    # shipping method, and checkout must be able to reject it.
    shipping: str = "standard"

    @classmethod
    def from_dict(cls, data: dict[str, Any]) -> Basket:
        """Builds a Basket from its JSON form (spec/checkout-rules.md, section 3)."""
        return cls(
            items=[
                BasketItem(sku=item["sku"], quantity=item["quantity"]) for item in data["items"]
            ],
            coupon=data.get("coupon"),
            shipping=data.get("shipping", "standard"),
        )


# ---------------------------------------------------------------------------
# The receipt: the output of checkout.
# ---------------------------------------------------------------------------


@dataclass(frozen=True)
class ReceiptLine:
    sku: str
    name: str
    quantity: int
    unit_price: int  # cents
    line_total: int  # cents: unit_price x quantity
    line_discount: int  # cents: the promotion discount for this line


@dataclass(frozen=True)
class Receipt:
    lines: list[ReceiptLine]
    subtotal: int
    promotion_discount: int
    coupon_discount: int
    shipping: int
    tax: int
    total: int
