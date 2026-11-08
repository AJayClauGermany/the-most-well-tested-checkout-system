"""Unit tests for a TEST HELPER: tests/support/dollars.py.

Test code can have bugs too. The BDD steps rely on "$12.39" becoming exactly 1239, so
the helper that does it gets its own tests - including the awkward cases.
"""

import pytest
from support.dollars import Money, dollars_to_cents


@pytest.mark.parametrize(
    ("text", "cents"),
    [
        ("$12.39", 1239),
        ("$0.50", 50),
        ("$0.04", 4),
        ("$5", 500),
        ("$5.5", 550),
        ("$0.00", 0),
        ("$159.45", 15945),
        ("12.39", 1239),
        ("  $1.00 ", 100),
    ],
    ids=repr,
)
def test_turns_dollars_into_whole_cents(text: str, cents: int) -> None:
    assert dollars_to_cents(text) == cents


def test_never_goes_through_floating_point() -> None:
    # 0.29 * 100 is 28.999999999999996 as a float - the classic bug. We must get 29.
    assert dollars_to_cents("$0.29") == 29
    assert dollars_to_cents("$1.15") == 115
    assert type(dollars_to_cents("$1.15")) is int


@pytest.mark.parametrize(
    "text", ["", "$", "$1.234", "$-1.00", "1,00", "one dollar", "$1."], ids=repr
)
def test_rejects_anything_that_is_not_a_dollar_amount(text: str) -> None:
    with pytest.raises(ValueError, match="not a dollar amount"):
        dollars_to_cents(text)


def test_money_is_a_custom_type_for_pytest_bdd_step_sentences() -> None:
    money = Money()

    assert money("$12.39") == 1239
    assert Money.pattern == r"\$\d+(?:\.\d{1,2})?"
