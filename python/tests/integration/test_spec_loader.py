"""Integration tests for the spec loader (tests/spec_loader.py).

An *integration test* checks that several pieces work together - here: our
code, the file system, the JSON parser and the real files in spec/. They are a
little slower than unit tests and can break for reasons outside our code
(for example, a typo in a JSON file), which is exactly what we want to catch.

The `catalog` and `spec_cases` parameters are fixtures from tests/conftest.py.
"""

from checkout import BuyXGetYFree


def test_the_standard_catalog_has_six_products(catalog):
    assert len(catalog.products) == 6


def test_coffee_is_buy_two_get_one_free(catalog):
    coffee = catalog.find("COFFEE")

    assert coffee is not None
    assert coffee.unit_price == 899
    assert coffee.promotion == BuyXGetYFree(buy=2, free=1)


def test_all_sixty_spec_cases_are_loaded(spec_cases):
    assert len(spec_cases) == 60


def test_spec_case_ids_are_unique(spec_cases):
    ids = [case.id for case in spec_cases]

    # A set drops duplicates, so the sizes only match if every id is unique.
    assert len(set(ids)) == len(ids)


def test_every_spec_case_expects_either_totals_or_an_error_but_not_both(spec_cases):
    for case in spec_cases:
        has_totals = case.expected is not None
        has_error = case.expected_error is not None
        # The message after the comma is shown when the assert fails.
        assert has_totals != has_error, f"{case.file}: {case.id}"
