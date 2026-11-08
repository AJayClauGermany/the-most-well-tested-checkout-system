"""A small checkout system. The rules live in spec/checkout-rules.md.

Everything you normally need can be imported straight from the package::

    from checkout import Basket, BasketItem, Catalog, CheckoutError, checkout
"""

from checkout.engine import checkout
from checkout.errors import CheckoutError, ErrorCode
from checkout.models import (
    Basket,
    BasketItem,
    BulkDiscount,
    BuyXGetYFree,
    Catalog,
    Product,
    Promotion,
    Receipt,
    ReceiptLine,
)

__all__ = [
    "Basket",
    "BasketItem",
    "BulkDiscount",
    "BuyXGetYFree",
    "Catalog",
    "CheckoutError",
    "ErrorCode",
    "Product",
    "Promotion",
    "Receipt",
    "ReceiptLine",
    "checkout",
]
