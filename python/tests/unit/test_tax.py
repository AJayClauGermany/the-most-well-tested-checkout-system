"""Unit tests for tax (rules §4.6)."""

import pytest

from checkout import TAX_PERCENT, tax_on


def test_tax_is_8_percent() -> None:
    assert TAX_PERCENT == 8


@pytest.mark.parametrize(
    ("amount", "tax"),
    [
        pytest.param(50, 4, id="8-percent-of-0.50-is-exactly-4-cents"),
        pytest.param(120, 10, id="9.6-cents-rounds-up"),
        pytest.param(240, 19, id="19.2-cents-rounds-down"),
        pytest.param(625, 50, id="exactly-50-cents"),
        pytest.param(1706, 136, id="136.48-rounds-down"),
        pytest.param(11880, 950, id="950.4-rounds-down"),
        pytest.param(0, 0, id="no-spend-no-tax"),
    ],
)
def test_taxes_the_discounted_subtotal(amount: int, tax: int) -> None:
    assert tax_on(amount) == tax


def test_rounds_plain_fractions_the_usual_way() -> None:
    # 8% of a whole number of cents can never be exactly half a cent (x * 8 never ends in 50),
    # so the half-up rule is pinned down by the 10% coupon tests; here we check plain rounding.
    assert tax_on(1875) == 150  # exactly 150.00
    assert tax_on(1881) == 150  # 150.48
    assert tax_on(1888) == 151  # 151.04
