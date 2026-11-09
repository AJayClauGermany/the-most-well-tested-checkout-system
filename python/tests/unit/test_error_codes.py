"""Example unit tests: start reading here if you are new to pytest!

A *unit test* checks one small piece of code (a "unit") on its own: no files,
no network, no database. Unit tests are fast, so you can run hundreds of them
every time you save.

How pytest works, in three rules:
  1. It collects files named ``test_*.py``.
  2. Inside them, it runs every function whose name starts with ``test_``.
  3. A test passes if it finishes without an error. ``assert`` raises an error
     when its condition is false, and pytest then shows you both values.

Most tests have three steps, often called *Arrange, Act, Assert*:
  Arrange - set up what the test needs
  Act     - run the code you are testing
  Assert  - check the result
"""

import pytest

from checkout import CheckoutError, ErrorCode


# A test name should say what *should* happen, like a sentence.
def test_a_checkout_error_remembers_its_error_code() -> None:
    # Arrange
    code = ErrorCode.UNKNOWN_COUPON

    # Act
    error = CheckoutError(code)

    # Assert
    assert error.code == ErrorCode.UNKNOWN_COUPON


def test_a_checkout_error_can_be_raised_and_caught() -> None:
    # pytest.raises checks that the code inside the `with` block raises
    # this exception. If nothing is raised, the test fails.
    with pytest.raises(CheckoutError) as caught:
        raise CheckoutError(ErrorCode.EMPTY_BASKET)

    # `caught.value` is the exception that was raised.
    assert caught.value.code == ErrorCode.EMPTY_BASKET


def test_without_a_message_the_code_is_the_message() -> None:
    error = CheckoutError(ErrorCode.EMPTY_BASKET)

    assert str(error) == "EMPTY_BASKET"


def test_with_a_message_the_message_is_kept_and_the_code_too() -> None:
    error = CheckoutError(ErrorCode.UNKNOWN_PRODUCT, "Unknown product: BANANA")

    assert str(error) == "Unknown product: BANANA"
    assert error.code == ErrorCode.UNKNOWN_PRODUCT


# The spec (spec/checkout-rules.md, section 9) lists the six error codes.
# Every language must use exactly this spelling.
SPEC_ERROR_CODES = [
    "EMPTY_BASKET",
    "UNKNOWN_PRODUCT",
    "INVALID_QUANTITY",
    "UNKNOWN_SHIPPING_METHOD",
    "UNKNOWN_COUPON",
    "COUPON_MINIMUM_NOT_MET",
]


# @pytest.mark.parametrize runs the same test once for each value in the list.
# The output shows one result per value, e.g. test_...[EMPTY_BASKET].
@pytest.mark.parametrize("spelling", SPEC_ERROR_CODES)
def test_each_error_code_is_spelled_exactly_as_in_the_spec(spelling: str) -> None:
    # ErrorCode("...") looks up the member with that value
    # (it raises ValueError if there is none, which would fail the test).
    code = ErrorCode(spelling)

    assert code.value == spelling
    assert str(code) == spelling


def test_there_are_exactly_six_error_codes() -> None:
    assert len(ErrorCode) == 6
