"""Acceptance tests: the shared spec cases from spec/test-cases/*.json.

An *acceptance test* asks "does the whole checkout do what the spec says?".
It only uses the public entry point, checkout(), so it does not care how the
checkout is written inside.

These tests were written BEFORE the code (test-driven development) and were skipped
until the implementation started. Now they run: 60 cases, the same for all six languages.
"""

import pytest
from spec_loader import ExpectedTotals, SpecCase, load_spec_cases

from checkout import Catalog, CheckoutError, checkout

# The cases are loaded when pytest collects this file, because
# parametrize needs the list before any test runs.
SPEC_CASES = load_spec_cases()


# ids= gives each generated test a readable name, e.g.
#   test_spec_case[promo-coffee-three-bags-one-free]
# so you can run just one case with:  pytest -k promo-coffee-three-bags-one-free
@pytest.mark.parametrize("case", SPEC_CASES, ids=[case.id for case in SPEC_CASES])
def test_spec_case(case: SpecCase, catalog: Catalog) -> None:
    if case.expected_error is not None:
        # This case describes an invalid basket: checkout must fail with that code.
        with pytest.raises(CheckoutError) as caught:
            checkout(case.basket, catalog)
        assert caught.value.code == case.expected_error
    else:
        # This case describes a valid basket: compare the six totals.
        receipt = checkout(case.basket, catalog)
        actual = ExpectedTotals(
            subtotal=receipt.subtotal,
            promotion_discount=receipt.promotion_discount,
            coupon_discount=receipt.coupon_discount,
            shipping=receipt.shipping,
            tax=receipt.tax,
            total=receipt.total,
        )
        # Comparing whole objects means a failure shows ALL the totals that differ.
        assert actual == case.expected
