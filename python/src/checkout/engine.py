"""The checkout engine: turns a basket into a receipt."""

from checkout.catalog import require_product
from checkout.coupons import coupon_discount
from checkout.lines import merge_items, price_line
from checkout.models import Basket, Catalog, Receipt
from checkout.shipping import parse_shipping_method, shipping_cost
from checkout.tax import tax_on
from checkout.validation import validate_items


def checkout(basket: Basket, catalog: Catalog) -> Receipt:
    """Prices ``basket`` using the products in ``catalog``.

    It reads like the hint in rules §4.7: validate, merge lines, price, coupon, shipping,
    tax, total. Each step is a small function in its own module with its own unit tests.

    Returns a Receipt, or raises ``CheckoutError`` with one of the six error codes from
    rules §9 when the basket is invalid.
    """
    validate_items(basket, catalog)  # rules §4.1 rows 1-4
    method = parse_shipping_method(basket.shipping)  # row 5

    lines = [
        price_line(require_product(catalog, item.sku), item.quantity)
        for item in merge_items(basket.items)
    ]
    subtotal = sum(line.line_total for line in lines)
    promotion_discount = sum(line.line_discount for line in lines)
    after_promotions = subtotal - promotion_discount

    coupon = coupon_discount(basket.coupon, after_promotions)  # rows 6-7
    discounted_subtotal = after_promotions - coupon

    shipping = shipping_cost(method, discounted_subtotal)
    tax = tax_on(discounted_subtotal)

    return Receipt(
        lines=lines,
        subtotal=subtotal,
        promotion_discount=promotion_discount,
        coupon_discount=coupon,
        shipping=shipping,
        tax=tax,
        total=discounted_subtotal + shipping + tax,
    )
