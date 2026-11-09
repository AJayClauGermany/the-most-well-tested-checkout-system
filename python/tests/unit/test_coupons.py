"""Unit tests for coupons (rules §4.4)."""

import pytest

from checkout import CheckoutError, ErrorCode, coupon_discount, normalize_coupon_code

# --- normalize_coupon_code -------------------------------------------------------------


@pytest.mark.parametrize(
    ("text", "expected"),
    [
        ("SAVE10", "SAVE10"),
        ("save10", "SAVE10"),
        ("  save10 ", "SAVE10"),
        ("\tWelcome5\n", "WELCOME5"),
        ("nope", "NOPE"),
    ],
    ids=repr,
)
def test_trims_and_upper_cases_the_code(text: str, expected: str) -> None:
    assert normalize_coupon_code(text) == expected


@pytest.mark.parametrize("text", [None, "", "   ", "\t\n"], ids=repr)
def test_treats_missing_empty_and_blank_as_no_coupon(text: str | None) -> None:
    assert normalize_coupon_code(text) is None


# --- coupon_discount ---------------------------------------------------------------


def test_is_0_without_a_coupon() -> None:
    assert coupon_discount(None, 10000) == 0
    assert coupon_discount("", 10000) == 0
    assert coupon_discount("  ", 10000) == 0


def test_save10_takes_10_percent_off_rounded_half_up() -> None:
    assert coupon_discount("SAVE10", 500) == 50
    assert coupon_discount("SAVE10", 225) == 23  # 22.5 -> 23
    assert coupon_discount("SAVE10", 224) == 22  # 22.4 -> 22
    assert coupon_discount("SAVE10", 1798) == 180  # 179.8 -> 180


def test_save10_has_no_minimum_spend() -> None:
    assert coupon_discount("SAVE10", 1) == 0  # 0.1 -> 0
    assert coupon_discount("SAVE10", 5) == 1  # 0.5 -> 1


def test_save10_ignores_case_and_surrounding_spaces() -> None:
    assert coupon_discount(" save10 ", 500) == 50


def test_welcome5_takes_exactly_5_dollars_off() -> None:
    assert coupon_discount("WELCOME5", 2500) == 500
    assert coupon_discount("WELCOME5", 99999) == 500


def test_welcome5_works_at_exactly_the_25_dollar_minimum_and_fails_one_cent_below() -> None:
    assert coupon_discount("WELCOME5", 2500) == 500

    with pytest.raises(CheckoutError) as caught:
        coupon_discount("WELCOME5", 2499)
    assert caught.value.code == ErrorCode.COUPON_MINIMUM_NOT_MET


def test_welcome5_ignores_case_and_surrounding_spaces() -> None:
    assert coupon_discount("welcome5", 2500) == 500


def test_rejects_a_code_it_does_not_know_with_unknown_coupon() -> None:
    with pytest.raises(CheckoutError) as caught:
        coupon_discount("FREESTUFF", 10000)

    assert caught.value.code == ErrorCode.UNKNOWN_COUPON


def test_explains_what_went_wrong_in_the_error_message() -> None:
    with pytest.raises(CheckoutError, match="^Unknown coupon: FREESTUFF$"):
        coupon_discount("freestuff", 10000)

    with pytest.raises(
        CheckoutError, match="^WELCOME5 needs a spend of at least 2500 cents, got 2499$"
    ):
        coupon_discount("WELCOME5", 2499)


@pytest.mark.parametrize("code", ["SAVE100", "WELCOME50", "SAVE1", "SAVE 10"], ids=repr)
def test_does_not_accept_a_code_that_merely_looks_like_a_real_code(code: str) -> None:
    with pytest.raises(CheckoutError) as caught:
        coupon_discount(code, 10000)

    assert caught.value.code == ErrorCode.UNKNOWN_COUPON
