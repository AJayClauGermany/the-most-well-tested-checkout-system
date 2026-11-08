"""Unit tests for the two kinds of promotion (rules §4.2), mostly as parameterised tables."""

import pytest

from checkout import BulkDiscount, BuyXGetYFree, promotion_discount

BUY_TWO_GET_ONE = BuyXGetYFree(buy=2, free=1)
TEN_PERCENT_OFF_TEN = BulkDiscount(min_quantity=10, percent=10)


def test_is_0_when_the_product_has_no_promotion() -> None:
    assert promotion_discount(None, 50, 899) == 0


# --- buy X, get Y free ---------------------------------------------------------------


@pytest.mark.parametrize(
    ("quantity", "free"),
    [(1, 0), (2, 0), (3, 1), (4, 1), (5, 1), (6, 2), (99, 33)],
    ids=lambda value: str(value),
)
def test_coffees_at_899_buy_2_get_1_free(quantity: int, free: int) -> None:
    assert promotion_discount(BUY_TWO_GET_ONE, quantity, 899) == free * 899


@pytest.mark.parametrize(
    ("promotion", "quantity", "unit_price", "discount"),
    [
        # Buy 1 get 1 free: every second item is free.
        (BuyXGetYFree(buy=1, free=1), 4, 100, 200),
        (BuyXGetYFree(buy=1, free=1), 5, 100, 200),
        # Buy 3 get 2 free: groups of 5, 2 of each group free.
        (BuyXGetYFree(buy=3, free=2), 4, 100, 0),
        (BuyXGetYFree(buy=3, free=2), 5, 100, 200),
        (BuyXGetYFree(buy=3, free=2), 10, 100, 400),
    ],
    ids=["1+1-x4", "1+1-x5", "3+2-x4", "3+2-x5", "3+2-x10"],
)
def test_uses_the_promotion_s_own_numbers_not_2_and_1(
    promotion: BuyXGetYFree, quantity: int, unit_price: int, discount: int
) -> None:
    assert promotion_discount(promotion, quantity, unit_price) == discount


def test_multiplies_the_free_items_by_the_unit_price() -> None:
    assert promotion_discount(BUY_TWO_GET_ONE, 3, 100) == 100
    assert promotion_discount(BUY_TWO_GET_ONE, 3, 1) == 1


# --- bulk discount ---------------------------------------------------------------------


@pytest.mark.parametrize(
    ("quantity", "discount"),
    [
        pytest.param(9, 0, id="9-one-short-of-the-minimum"),
        pytest.param(10, 349, id="10-exactly-the-minimum-10-percent-of-34.90"),
        pytest.param(11, 384, id="11-383.9-cents-rounds-up"),
        pytest.param(15, 524, id="15-exactly-523.5-cents-rounds-UP"),
        pytest.param(99, 3455, id="99-3455.1-cents-rounds-down"),
    ],
)
def test_notebooks_at_349_ten_percent_off_from_10(quantity: int, discount: int) -> None:
    assert promotion_discount(TEN_PERCENT_OFF_TEN, quantity, 349) == discount


def test_uses_the_promotion_s_own_minimum_and_percentage() -> None:
    quarter_off_from_three = BulkDiscount(min_quantity=3, percent=25)

    assert promotion_discount(quarter_off_from_three, 2, 100) == 0
    assert promotion_discount(quarter_off_from_three, 3, 100) == 75
    assert promotion_discount(quarter_off_from_three, 3, 101) == 76  # 75.75 -> 76
