"""Unit tests for rounding (rules §7).

Rounding is where checkout systems classically go wrong, so it gets its own small
function and its own tests. This file is also the project's first PARAMETERISED test:
``@pytest.mark.parametrize`` runs one test function once per row of a table, and
``ids=`` gives every row a readable name in the output.
"""

import pytest

from checkout import round_half_up

# The table from rules §7: (amount, percent, the exact value, what it must round to).
ROUNDING_TABLE = [
    (224, 10, "22.4", 22),
    (225, 10, "22.5", 23),
    (226, 10, "22.6", 23),
]


@pytest.mark.parametrize(
    ("amount", "percent", "expected"),
    [(amount, percent, expected) for amount, percent, _exact, expected in ROUNDING_TABLE],
    ids=[f"{exact}-rounds-to-{expected}" for _a, _p, exact, expected in ROUNDING_TABLE],
)
def test_rounds_like_the_table_in_the_rules(amount: int, percent: int, expected: int) -> None:
    assert round_half_up(amount, percent) == expected


def test_rounds_just_under_a_half_down() -> None:
    # 2249 x 1% is 22.49 cents: one hundredth of a cent below the half, so it rounds DOWN.
    # Without this example, "+ 51" instead of "+ 50" in round_half_up would go unnoticed
    # (mutation testing found that only a random property test happened to catch it).
    assert round_half_up(2249, 1) == 22
    assert round_half_up(2250, 1) == 23


def test_rounds_an_exact_half_up_where_bankers_rounding_would_go_to_the_even_number() -> None:
    # Python's built-in round() uses "round half to even": round(22.5) == 22.
    # The rules say halves go UP. Both examples together pin the rule down.
    assert round_half_up(225, 10) == 23  # 22.5 -> 23 (half-to-even would say 22)
    assert round_half_up(5235, 10) == 524  # 523.5 -> 524


def test_does_not_round_when_the_result_is_already_whole() -> None:
    assert round_half_up(500, 10) == 50
    assert round_half_up(50, 8) == 4


def test_rounds_9_point_6_cents_up_to_10_and_19_point_2_cents_down_to_19() -> None:
    assert round_half_up(120, 8) == 10
    assert round_half_up(240, 8) == 19


def test_gives_0_for_an_amount_of_0_and_the_whole_amount_for_100_percent() -> None:
    assert round_half_up(0, 10) == 0
    assert round_half_up(899, 100) == 899


def test_gives_0_for_0_percent() -> None:
    assert round_half_up(899, 0) == 0


def test_stays_exact_where_floating_point_maths_drifts() -> None:
    # 0.1 + 0.2 != 0.3 in Python too. 10% of 3 cents is 0.3 - it must round to 0, not 1.
    assert round_half_up(3, 10) == 0
    # 8% of 1050 is exactly 84.0; a float implementation might land on 83.99999.
    assert round_half_up(1050, 8) == 84


def test_returns_a_plain_int_not_a_float() -> None:
    # Money is whole cents (rules §1): the type matters, not just the value.
    assert type(round_half_up(225, 10)) is int
