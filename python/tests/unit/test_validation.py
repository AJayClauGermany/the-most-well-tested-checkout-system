"""Unit tests for basket validation (rules §4.1, rows 1 to 4) and the order errors are
reported in. The check order matters: a basket can be wrong in several ways at once.
"""

import pytest
from support.fixtures import TINY_CATALOG

from checkout import (
    MAX_QUANTITY,
    MIN_QUANTITY,
    Basket,
    BasketItem,
    CheckoutError,
    ErrorCode,
    is_valid_quantity,
    validate_items,
)


def expect_error(basket: Basket, code: ErrorCode) -> CheckoutError:
    """Asserts that validating `basket` fails with exactly this code, and returns the error."""
    with pytest.raises(CheckoutError) as caught:
        validate_items(basket, TINY_CATALOG)
    assert caught.value.code == code
    return caught.value


# --- is_valid_quantity --------------------------------------------------------


def test_the_limits_are_the_ones_from_the_rules() -> None:
    assert MIN_QUANTITY == 1
    assert MAX_QUANTITY == 99


@pytest.mark.parametrize("quantity", [1, 2, 50, 98, 99], ids=str)
def test_accepts_whole_numbers_from_1_to_99(quantity: int) -> None:
    assert is_valid_quantity(quantity) is True


@pytest.mark.parametrize("quantity", [0, -1, 100, 101, 1000], ids=str)
def test_rejects_numbers_outside_1_to_99(quantity: int) -> None:
    assert is_valid_quantity(quantity) is False


@pytest.mark.parametrize(
    "quantity",
    [1.5, 0.999, 99.0001, 3.0, float("nan"), float("inf"), "3", None, True, [3]],
    ids=[
        "1.5",
        "0.999",
        "99.0001",
        "3.0-is-a-float",
        "nan",
        "inf",
        "string",
        "None",
        "bool",
        "list",
    ],
)
def test_rejects_anything_that_is_not_a_whole_number(quantity: object) -> None:
    # Only a real `int` counts. 3.0 is a float, and True is a bool (even though
    # Python treats True as 1 in arithmetic) - neither is a quantity.
    assert is_valid_quantity(quantity) is False


# --- validate_items ------------------------------------------------------------


def test_accepts_a_normal_basket_and_returns_nothing() -> None:
    basket = Basket(items=[BasketItem("APPLE", 1), BasketItem("COFFEE", 99)])

    validate_items(basket, TINY_CATALOG)  # must not raise


def test_row_1_an_empty_basket() -> None:
    expect_error(Basket(items=[]), ErrorCode.EMPTY_BASKET)


def test_row_2_a_sku_that_is_not_in_the_catalog() -> None:
    expect_error(Basket(items=[BasketItem("BANANA", 1)]), ErrorCode.UNKNOWN_PRODUCT)


def test_row_2_skus_are_case_sensitive() -> None:
    expect_error(Basket(items=[BasketItem("apple", 1)]), ErrorCode.UNKNOWN_PRODUCT)


def test_row_2_an_unknown_product_on_any_line_not_only_the_first() -> None:
    basket = Basket(items=[BasketItem("APPLE", 1), BasketItem("PEN", 1), BasketItem("BANANA", 1)])

    expect_error(basket, ErrorCode.UNKNOWN_PRODUCT)


@pytest.mark.parametrize("quantity", [0, -1, 100], ids=str)
def test_row_3_a_quantity_outside_1_to_99(quantity: int) -> None:
    expect_error(Basket(items=[BasketItem("APPLE", quantity)]), ErrorCode.INVALID_QUANTITY)


def test_row_3_a_bad_quantity_on_a_later_line() -> None:
    basket = Basket(items=[BasketItem("APPLE", 1), BasketItem("PEN", 0)])

    expect_error(basket, ErrorCode.INVALID_QUANTITY)


def test_row_4_lines_of_one_sku_that_add_up_to_more_than_99() -> None:
    basket = Basket(items=[BasketItem("APPLE", 50), BasketItem("APPLE", 50)])

    expect_error(basket, ErrorCode.INVALID_QUANTITY)


def test_row_4_lines_of_one_sku_that_add_up_to_exactly_99_are_fine() -> None:
    basket = Basket(items=[BasketItem("APPLE", 50), BasketItem("PEN", 99), BasketItem("APPLE", 49)])

    validate_items(basket, TINY_CATALOG)  # must not raise


def test_row_4_99_of_one_sku_and_99_of_another_is_fine_the_limit_is_per_sku() -> None:
    basket = Basket(items=[BasketItem("APPLE", 99), BasketItem("PEN", 99)])

    validate_items(basket, TINY_CATALOG)  # must not raise


# --- error messages -------------------------------------------------------------
# The code is what programs look at; the message is for the human reading a log.


def test_the_messages_explain_the_problem() -> None:
    assert str(expect_error(Basket(items=[]), ErrorCode.EMPTY_BASKET)) == "The basket has no items"
    assert (
        str(expect_error(Basket(items=[BasketItem("BANANA", 1)]), ErrorCode.UNKNOWN_PRODUCT))
        == "Unknown product: BANANA"
    )
    assert (
        str(expect_error(Basket(items=[BasketItem("APPLE", 0)]), ErrorCode.INVALID_QUANTITY))
        == "Invalid quantity for APPLE: 0"
    )
    too_many = Basket(items=[BasketItem("APPLE", 60), BasketItem("APPLE", 40)])
    assert str(expect_error(too_many, ErrorCode.INVALID_QUANTITY)) == "Too many APPLE in total: 100"


# --- check order -----------------------------------------------------------------


def test_an_unknown_product_beats_a_bad_quantity_even_on_a_later_line() -> None:
    basket = Basket(items=[BasketItem("APPLE", 0), BasketItem("BANANA", 1)])

    expect_error(basket, ErrorCode.UNKNOWN_PRODUCT)


def test_a_bad_single_quantity_beats_a_bad_merged_quantity() -> None:
    # Both rows 3 and 4 apply; row 3 is reported (the message tells them apart).
    basket = Basket(items=[BasketItem("APPLE", 100), BasketItem("APPLE", 1)])

    error = expect_error(basket, ErrorCode.INVALID_QUANTITY)

    assert str(error) == "Invalid quantity for APPLE: 100"


def test_ignores_the_coupon_and_shipping_method_those_are_checked_later_by_checkout() -> None:
    basket = Basket(items=[BasketItem("APPLE", 1)], coupon="NOPE", shipping="overnight")

    validate_items(basket, TINY_CATALOG)  # must not raise
