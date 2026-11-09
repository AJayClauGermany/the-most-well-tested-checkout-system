"""Contract tests: is the output even the right shape?

spec/schema/receipt.schema.json is the contract every language's receipt must honour.
Here the jsonschema library checks real receipts against it. The e2e runner only compares
totals; this is where the shape of the whole receipt (lines included) is held to the
contract - and the CLI's error output to cli-error.schema.json.
"""

import json

import pytest
from spec_loader import (
    SpecCase,
    invalid_spec_cases,
    load_catalog,
    load_spec_cases,
    valid_spec_cases,
)
from support.fixtures import TINY_CATALOG
from support.schemas import schema_errors

from checkout import Basket, BasketItem, checkout, run_checkout_command

CATALOG = load_catalog()
VALID_CASES = valid_spec_cases()
INVALID_CASES = invalid_spec_cases()

# --- receipt.schema.json ----------------------------------------------------------------


@pytest.mark.parametrize("case", VALID_CASES, ids=[case.id for case in VALID_CASES])
def test_the_receipt_for_a_spec_case_follows_the_schema(case: SpecCase) -> None:
    receipt = checkout(case.basket, CATALOG)

    assert schema_errors("receipt.schema.json", receipt.to_dict()) == []


def test_a_receipt_from_a_made_up_catalog_follows_the_schema_too() -> None:
    receipt = checkout(Basket(items=[BasketItem("TV", 1), BasketItem("COFFEE", 6)]), TINY_CATALOG)

    assert schema_errors("receipt.schema.json", receipt.to_dict()) == []


def test_the_schema_check_itself_works_a_broken_receipt_is_rejected() -> None:
    # Testing the test: a checker that accepts everything would make every test above useless.
    receipt = checkout(Basket(items=[BasketItem("APPLE", 1)]), TINY_CATALOG).to_dict()
    broken = {**receipt, "shipping": 123, "extra": True, "tax": -1}

    assert schema_errors("receipt.schema.json", broken) == [
        "/ Additional properties are not allowed ('extra' was unexpected)",
        "/shipping 123 is not one of [0, 499, 999]",
        "/tax -1 is less than the minimum of 0",
    ]


def test_a_receipt_with_a_broken_line_is_rejected_with_the_path_to_the_line() -> None:
    receipt = checkout(Basket(items=[BasketItem("APPLE", 1)]), TINY_CATALOG).to_dict()
    receipt["lines"][0]["quantity"] = 100

    assert schema_errors("receipt.schema.json", receipt) == [
        "/lines/0/quantity 100 is greater than the maximum of 99"
    ]


# --- cli-error.schema.json ---------------------------------------------------------------


@pytest.mark.parametrize("case", INVALID_CASES, ids=[case.id for case in INVALID_CASES])
def test_the_command_s_output_for_an_invalid_spec_case_follows_the_schema(case: SpecCase) -> None:
    result = run_checkout_command(json.dumps(case.basket_json), CATALOG)

    assert result.exit_code == 1
    assert schema_errors("cli-error.schema.json", json.loads(result.output)) == []


def test_an_unknown_error_code_is_rejected_by_the_schema() -> None:
    assert schema_errors("cli-error.schema.json", {"error": "OUT_OF_STOCK"}) == [
        "/error 'OUT_OF_STOCK' is not one of ['EMPTY_BASKET', 'UNKNOWN_PRODUCT', "
        "'INVALID_QUANTITY', 'UNKNOWN_SHIPPING_METHOD', 'UNKNOWN_COUPON', "
        "'COUPON_MINIMUM_NOT_MET']"
    ]


# --- basket.schema.json and order.schema.json --------------------------------------------


def test_every_shared_test_case_s_basket_is_a_valid_basket() -> None:
    # The spec is consistent with itself: every input the cases use fits the input schema.
    for case in load_spec_cases():
        assert schema_errors("basket.schema.json", case.basket_json) == [], case.id


def test_the_order_schema_reaches_into_the_receipt_schema_through_its_ref() -> None:
    # order.schema.json says "$ref": "receipt.schema.json"; a broken receipt inside an
    # otherwise fine order must therefore be reported.
    order = {"id": "ORDER-1", "created_at": "2026-01-15T10:30:00Z", "receipt": {"total": -5}}

    errors = schema_errors("order.schema.json", order)

    assert "/receipt/total -5 is less than the minimum of 0" in errors
    assert any(error.startswith("/receipt 'lines' is a required property") for error in errors)
