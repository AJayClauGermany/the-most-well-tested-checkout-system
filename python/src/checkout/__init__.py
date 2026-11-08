"""A small checkout system. The rules live in spec/checkout-rules.md.

Everything you normally need can be imported straight from the package::

    from checkout import Basket, BasketItem, Catalog, CheckoutError, checkout

The checkout itself is a pipeline of small steps (rules §4), each in its own module with
its own unit tests: validation -> lines -> promotions -> coupons -> shipping -> tax.
Saving orders (rules §6) lives in orders.py and the two file-based repositories; the
command line program (rules §8) in command.py and __main__.py.
"""

from checkout.catalog import require_product
from checkout.command import (
    CommandResult,
    NotABasketError,
    load_standard_catalog,
    main,
    parse_basket,
    run_checkout_command,
    standard_catalog_path,
)
from checkout.coupons import coupon_discount, normalize_coupon_code
from checkout.engine import checkout
from checkout.errors import CheckoutError, ErrorCode
from checkout.json_file_repository import JsonFileOrderRepository
from checkout.lines import merge_items, price_line
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
from checkout.money import round_half_up
from checkout.orders import (
    Clock,
    InMemoryOrderRepository,
    Order,
    OrderRepository,
    PlacedOrder,
    format_timestamp,
    order_id,
    order_number,
    place_order,
    system_clock,
)
from checkout.promotions import promotion_discount
from checkout.shipping import (
    EXPRESS_SHIPPING_COST,
    FREE_SHIPPING_FROM,
    STANDARD_SHIPPING_COST,
    ShippingMethod,
    parse_shipping_method,
    shipping_cost,
)
from checkout.spec_folder import find_spec_folder
from checkout.sqlite_repository import SqliteOrderRepository
from checkout.tax import TAX_PERCENT, tax_on
from checkout.validation import MAX_QUANTITY, MIN_QUANTITY, is_valid_quantity, validate_items

__all__ = [
    "EXPRESS_SHIPPING_COST",
    "FREE_SHIPPING_FROM",
    "MAX_QUANTITY",
    "MIN_QUANTITY",
    "STANDARD_SHIPPING_COST",
    "TAX_PERCENT",
    "Basket",
    "BasketItem",
    "BulkDiscount",
    "BuyXGetYFree",
    "Catalog",
    "CheckoutError",
    "Clock",
    "CommandResult",
    "ErrorCode",
    "InMemoryOrderRepository",
    "JsonFileOrderRepository",
    "NotABasketError",
    "Order",
    "OrderRepository",
    "PlacedOrder",
    "Product",
    "Promotion",
    "Receipt",
    "ReceiptLine",
    "ShippingMethod",
    "SqliteOrderRepository",
    "checkout",
    "coupon_discount",
    "find_spec_folder",
    "format_timestamp",
    "is_valid_quantity",
    "load_standard_catalog",
    "main",
    "merge_items",
    "normalize_coupon_code",
    "order_id",
    "order_number",
    "parse_basket",
    "parse_shipping_method",
    "place_order",
    "price_line",
    "promotion_discount",
    "require_product",
    "round_half_up",
    "run_checkout_command",
    "shipping_cost",
    "standard_catalog_path",
    "system_clock",
    "tax_on",
    "validate_items",
]
