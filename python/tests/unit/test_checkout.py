"""Unit tests of checkout() as a whole, with the tiny made-up catalog from tests/support/.

The 60 shared spec cases (tests/acceptance/) cover the standard catalog; these tests
check that the pieces fit together, and the things the totals alone cannot show, such
as the receipt lines and the order in which errors are reported.
"""

import pytest
from support.fixtures import TINY_CATALOG

from checkout import Basket, BasketItem, CheckoutError, ErrorCode, Receipt, ReceiptLine, checkout


def expect_error(basket: Basket, code: ErrorCode) -> None:
    with pytest.raises(CheckoutError) as caught:
        checkout(basket, TINY_CATALOG)
    assert caught.value.code == code


def test_prices_one_apple_50_plus_499_shipping_plus_4_tax_is_553() -> None:
    receipt = checkout(Basket(items=[BasketItem("APPLE", 1)]), TINY_CATALOG)

    assert receipt == Receipt(
        lines=[
            ReceiptLine(
                sku="APPLE", name="Apple", quantity=1, unit_price=50, line_total=50, line_discount=0
            )
        ],
        subtotal=50,
        promotion_discount=0,
        coupon_discount=0,
        shipping=499,
        tax=4,
        total=553,
    )


def test_writes_one_line_per_distinct_sku_in_order_of_first_appearance_with_merged_quantities() -> (
    None
):
    basket = Basket(items=[BasketItem("PEN", 1), BasketItem("APPLE", 2), BasketItem("PEN", 3)])

    receipt = checkout(basket, TINY_CATALOG)

    assert [(line.sku, line.quantity) for line in receipt.lines] == [("PEN", 4), ("APPLE", 2)]


def test_works_the_totals_out_from_the_lines() -> None:
    # 3 coffees ($26.97, one free: -$8.99) and 10 notebooks ($34.90, 10% off: -$3.49)
    basket = Basket(items=[BasketItem("COFFEE", 3), BasketItem("NOTEBOOK", 10)])

    receipt = checkout(basket, TINY_CATALOG)

    assert receipt.subtotal == 2697 + 3490
    assert receipt.promotion_discount == 899 + 349
    assert receipt.coupon_discount == 0
    assert receipt.shipping == 499  # $49.39 after promotions: not free
    assert receipt.tax == 395  # 8% of 4939 = 395.12
    assert receipt.total == 4939 + 499 + 395


def test_takes_the_coupon_off_after_promotions_and_before_shipping_and_tax() -> None:
    # TV $100.00 with SAVE10: $90.00 -> free shipping, tax 8% of $90.00 = $7.20
    basket = Basket(items=[BasketItem("TV", 1)], coupon="SAVE10")

    receipt = checkout(basket, TINY_CATALOG)

    assert (receipt.subtotal, receipt.promotion_discount, receipt.coupon_discount) == (
        10000,
        0,
        1000,
    )
    assert (receipt.shipping, receipt.tax, receipt.total) == (0, 720, 9720)


def test_charges_express_shipping_on_top_untaxed() -> None:
    basket = Basket(items=[BasketItem("TV", 1)], shipping="express")

    receipt = checkout(basket, TINY_CATALOG)

    assert (receipt.shipping, receipt.tax, receipt.total) == (999, 800, 10000 + 999 + 800)


def test_does_not_change_the_basket_it_is_given() -> None:
    basket = Basket(items=[BasketItem("APPLE", 1), BasketItem("APPLE", 1)], coupon=" save10 ")

    checkout(basket, TINY_CATALOG)

    assert basket == Basket(
        items=[BasketItem("APPLE", 1), BasketItem("APPLE", 1)], coupon=" save10 "
    )


def test_a_receipt_can_be_turned_into_json_and_back() -> None:
    receipt = checkout(Basket(items=[BasketItem("COFFEE", 3)], coupon="SAVE10"), TINY_CATALOG)

    as_json = receipt.to_dict()

    # snake_case keys, exactly as spec/schema/receipt.schema.json spells them
    assert set(as_json) == {
        "lines",
        "subtotal",
        "promotion_discount",
        "coupon_discount",
        "shipping",
        "tax",
        "total",
    }
    assert as_json["lines"][0] == {
        "sku": "COFFEE",
        "name": "Bag of Coffee Beans",
        "quantity": 3,
        "unit_price": 899,
        "line_total": 2697,
        "line_discount": 899,
    }
    assert Receipt.from_dict(as_json) == receipt


# --- reports the first failing row of rules §4.1 ------------------------------------


def test_empty_basket_before_anything_else() -> None:
    expect_error(Basket(items=[], coupon="NOPE", shipping="overnight"), ErrorCode.EMPTY_BASKET)


def test_unknown_product_before_bad_quantity_shipping_and_coupon() -> None:
    basket = Basket(
        items=[BasketItem("APPLE", 0), BasketItem("BANANA", 1)], coupon="NOPE", shipping="x"
    )

    expect_error(basket, ErrorCode.UNKNOWN_PRODUCT)


def test_bad_quantity_before_shipping_and_coupon() -> None:
    basket = Basket(items=[BasketItem("APPLE", 0)], coupon="NOPE", shipping="x")

    expect_error(basket, ErrorCode.INVALID_QUANTITY)


def test_merged_quantity_over_99_before_shipping_and_coupon() -> None:
    basket = Basket(
        items=[BasketItem("APPLE", 60), BasketItem("APPLE", 40)], coupon="NOPE", shipping="x"
    )

    expect_error(basket, ErrorCode.INVALID_QUANTITY)


def test_bad_shipping_before_bad_coupon() -> None:
    basket = Basket(items=[BasketItem("APPLE", 1)], coupon="NOPE", shipping="x")

    expect_error(basket, ErrorCode.UNKNOWN_SHIPPING_METHOD)


def test_bad_shipping_before_an_unmet_coupon_minimum() -> None:
    basket = Basket(items=[BasketItem("APPLE", 1)], coupon="WELCOME5", shipping="x")

    expect_error(basket, ErrorCode.UNKNOWN_SHIPPING_METHOD)


def test_unknown_coupon() -> None:
    expect_error(Basket(items=[BasketItem("APPLE", 1)], coupon="NOPE"), ErrorCode.UNKNOWN_COUPON)


def test_coupon_minimum_not_met_judged_on_the_price_after_promotions() -> None:
    # 3 coffees are $26.97 but $17.98 after the free one: below WELCOME5's $25.00.
    basket = Basket(items=[BasketItem("COFFEE", 3)], coupon="WELCOME5")

    expect_error(basket, ErrorCode.COUPON_MINIMUM_NOT_MET)
