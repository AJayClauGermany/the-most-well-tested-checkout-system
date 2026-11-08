"""Integration tests for the spec loader (tests/spec_loader.py).

An *integration test* checks that several pieces work together - here: our
code, the file system, the JSON parser and the real files in spec/. They are a
little slower than unit tests and can break for reasons outside our code
(for example, a typo in a JSON file), which is exactly what we want to catch.

The `catalog` and `spec_cases` parameters are fixtures from tests/conftest.py.
"""

from spec_loader import SpecCase, invalid_spec_cases, valid_spec_cases

from checkout import BuyXGetYFree, Catalog


def test_the_standard_catalog_has_six_products(catalog: Catalog) -> None:
    assert len(catalog.products) == 6


def test_coffee_is_buy_two_get_one_free(catalog: Catalog) -> None:
    coffee = catalog.find("COFFEE")

    assert coffee is not None
    assert coffee.unit_price == 899
    assert coffee.promotion == BuyXGetYFree(buy=2, free=1)


def test_all_sixty_spec_cases_are_loaded(spec_cases: list[SpecCase]) -> None:
    assert len(spec_cases) == 60


def test_spec_case_ids_are_unique(spec_cases: list[SpecCase]) -> None:
    ids = [case.id for case in spec_cases]

    # A set drops duplicates, so the sizes only match if every id is unique.
    assert len(set(ids)) == len(ids)


def test_every_spec_case_expects_either_totals_or_an_error_but_not_both(
    spec_cases: list[SpecCase],
) -> None:
    for case in spec_cases:
        has_totals = case.expected is not None
        has_error = case.expected_error is not None
        # The message after the comma is shown when the assert fails.
        assert has_totals != has_error, f"{case.file}: {case.id}"


def test_valid_and_invalid_cases_together_are_all_the_cases(spec_cases: list[SpecCase]) -> None:
    assert len(valid_spec_cases()) + len(invalid_spec_cases()) == len(spec_cases)
    assert all(case.expected is not None for case in valid_spec_cases())
    assert all(case.expected_error is not None for case in invalid_spec_cases())


def test_the_raw_basket_json_is_kept_next_to_the_parsed_basket(spec_cases: list[SpecCase]) -> None:
    example = next(case for case in spec_cases if case.id == "full-receipt-example")

    assert example.basket_json == {
        "items": [{"sku": "COFFEE", "quantity": 3}, {"sku": "APPLE", "quantity": 2}],
        "coupon": "SAVE10",
        "shipping": "standard",
    }
    assert example.basket.coupon == "SAVE10"
