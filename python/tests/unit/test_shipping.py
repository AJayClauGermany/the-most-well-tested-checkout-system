"""Unit tests for shipping (rules §4.5)."""

import pytest

from checkout import (
    EXPRESS_SHIPPING_COST,
    FREE_SHIPPING_FROM,
    STANDARD_SHIPPING_COST,
    CheckoutError,
    ErrorCode,
    parse_shipping_method,
    shipping_cost,
)

# --- parse_shipping_method --------------------------------------------------------


def test_the_constants_are_the_ones_from_the_rules() -> None:
    assert STANDARD_SHIPPING_COST == 499
    assert EXPRESS_SHIPPING_COST == 999
    assert FREE_SHIPPING_FROM == 5000


def test_accepts_the_two_known_methods() -> None:
    assert parse_shipping_method("standard") == "standard"
    assert parse_shipping_method("express") == "express"


@pytest.mark.parametrize(
    "method", ["overnight", "Express", "STANDARD", " standard", "express ", ""], ids=repr
)
def test_rejects_anything_else_with_unknown_shipping_method(method: str) -> None:
    with pytest.raises(CheckoutError) as caught:
        parse_shipping_method(method)

    assert caught.value.code == ErrorCode.UNKNOWN_SHIPPING_METHOD


def test_the_error_message_names_the_method_it_did_not_recognise() -> None:
    with pytest.raises(CheckoutError, match="^Unknown shipping method: overnight$"):
        parse_shipping_method("overnight")


# --- shipping_cost -------------------------------------------------------------------


def test_standard_costs_4_99_below_50_dollars() -> None:
    assert shipping_cost("standard", 0) == 499
    assert shipping_cost("standard", 4999) == 499


def test_standard_is_free_from_exactly_50_dollars() -> None:
    assert shipping_cost("standard", 5000) == 0
    assert shipping_cost("standard", 5001) == 0
    assert shipping_cost("standard", 100000) == 0


@pytest.mark.parametrize("discounted_subtotal", [0, 4999, 5000, 100000], ids=str)
def test_express_always_costs_9_99(discounted_subtotal: int) -> None:
    assert shipping_cost("express", discounted_subtotal) == 999
