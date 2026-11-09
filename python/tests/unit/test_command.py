"""Unit tests for the command line program's logic (rules §8), run in-process.

The real program (`python -m checkout`) is tested from the outside in tests/e2e/. Here
`main()` is called directly, with pytest's `monkeypatch` fixture standing in for standard
input and its `capsys` fixture capturing what was printed.
"""

import io
import json
import sys

import pytest
from support.fixtures import TINY_CATALOG

from checkout import (
    Basket,
    BasketItem,
    CommandResult,
    NotABasketError,
    load_standard_catalog,
    main,
    parse_basket,
    run_checkout_command,
    standard_catalog_path,
)

ONE_APPLE = '{"items":[{"sku":"APPLE","quantity":1}]}'

# --- the standard catalog ------------------------------------------------------------


def test_the_standard_catalog_is_spec_catalog_json_at_the_root_of_the_repository() -> None:
    path = standard_catalog_path()

    assert path.name == "catalog.json"
    assert path.parent.name == "spec"
    # The repository root is the folder that also holds the rules.
    assert (path.parent / "checkout-rules.md").is_file()


def test_loads_the_six_standard_products() -> None:
    catalog = load_standard_catalog()

    assert [product.sku for product in catalog.products] == [
        "APPLE",
        "BREAD",
        "COFFEE",
        "NOTEBOOK",
        "PEN",
        "HEADPHONES",
    ]


# --- parse_basket ----------------------------------------------------------------------


def test_parses_a_full_basket() -> None:
    basket = parse_basket(
        '{"items":[{"sku":"COFFEE","quantity":3},{"sku":"APPLE","quantity":2}],'
        '"coupon":"SAVE10","shipping":"express"}'
    )

    assert basket == Basket(
        items=[BasketItem("COFFEE", 3), BasketItem("APPLE", 2)], coupon="SAVE10", shipping="express"
    )


def test_a_missing_coupon_and_shipping_mean_no_coupon_and_standard() -> None:
    basket = parse_basket(ONE_APPLE)

    assert basket.coupon is None
    assert basket.shipping == "standard"


def test_a_null_coupon_means_no_coupon() -> None:
    assert parse_basket('{"items":[{"sku":"APPLE","quantity":1}],"coupon":null}').coupon is None


def test_keeps_invalid_quantities_so_that_checkout_can_report_them() -> None:
    # 0 and 100 are integers, so this IS a basket - just not a valid one (rules §4.1 row 3).
    basket = parse_basket('{"items":[{"sku":"APPLE","quantity":0},{"sku":"PEN","quantity":100}]}')

    assert basket.items == [BasketItem("APPLE", 0), BasketItem("PEN", 100)]


def test_accepts_surrounding_whitespace_as_a_shell_pipeline_adds_a_newline() -> None:
    assert parse_basket(f"  {ONE_APPLE}\n") == Basket(items=[BasketItem("APPLE", 1)])


def test_raises_json_decode_error_when_the_text_is_not_json() -> None:
    with pytest.raises(json.JSONDecodeError):
        parse_basket("not json")
    with pytest.raises(json.JSONDecodeError):
        parse_basket("")


# (JSON text, why it is not a basket)
NOT_BASKETS = [
    ("null", "expected a JSON object"),
    ("42", "expected a JSON object"),
    ('"text"', "expected a JSON object"),
    ("[]", "expected a JSON object"),
    ("{}", "'items' must be a list"),
    ('{"items": "APPLE"}', "'items' must be a list"),
    ('{"items": [1, 2]}', "every item must be an object"),
    ('{"items": [null]}', "every item must be an object"),
    ('{"items": ["APPLE"]}', "every item must be an object"),
    ('{"items": [{"sku": 5, "quantity": 1}]}', "every item needs a string 'sku'"),
    ('{"items": [{"quantity": 1}]}', "every item needs a string 'sku'"),
    ('{"items": [{"sku": "APPLE"}]}', "every item needs a whole-number 'quantity'"),
    (
        '{"items": [{"sku": "APPLE", "quantity": "1"}]}',
        "every item needs a whole-number 'quantity'",
    ),
    (
        '{"items": [{"sku": "APPLE", "quantity": 1.5}]}',
        "every item needs a whole-number 'quantity'",
    ),
    (
        '{"items": [{"sku": "APPLE", "quantity": true}]}',
        "every item needs a whole-number 'quantity'",
    ),
    (
        '{"items": [{"sku": "APPLE", "quantity": 1}], "coupon": 10}',
        "'coupon' must be a string or null",
    ),
    (
        '{"items": [{"sku": "APPLE", "quantity": 1}], "shipping": null}',
        "'shipping' must be a string",
    ),
    ('{"items": [{"sku": "APPLE", "quantity": 1}], "shipping": 1}', "'shipping' must be a string"),
]


@pytest.mark.parametrize(("text", "reason"), NOT_BASKETS, ids=[text for text, _ in NOT_BASKETS])
def test_raises_not_a_basket_error_when_the_json_is_not_a_basket(text: str, reason: str) -> None:
    with pytest.raises(NotABasketError) as caught:
        parse_basket(text)

    assert str(caught.value) == f"The input is not a basket: {reason}"


def test_not_a_basket_error_is_a_value_error() -> None:
    # Callers that just want "the input was bad" can catch ValueError, like for bad JSON.
    assert issubclass(NotABasketError, ValueError)


# --- run_checkout_command ------------------------------------------------------------


def test_prints_the_receipt_as_json_and_exits_with_0_for_a_valid_basket() -> None:
    result = run_checkout_command(ONE_APPLE, TINY_CATALOG)

    assert result.exit_code == 0
    assert json.loads(result.output) == {
        "lines": [
            {
                "sku": "APPLE",
                "name": "Apple",
                "quantity": 1,
                "unit_price": 50,
                "line_total": 50,
                "line_discount": 0,
            }
        ],
        "subtotal": 50,
        "promotion_discount": 0,
        "coupon_discount": 0,
        "shipping": 499,
        "tax": 4,
        "total": 553,
    }


def test_prints_the_error_code_as_json_and_exits_with_1_for_an_invalid_basket() -> None:
    result = run_checkout_command('{"items":[]}', TINY_CATALOG)

    assert result.exit_code == 1
    assert json.loads(result.output) == {"error": "EMPTY_BASKET"}


def test_reports_every_kind_of_error_the_same_way() -> None:
    result = run_checkout_command(
        '{"items":[{"sku":"APPLE","quantity":1}],"coupon":"NOPE"}', TINY_CATALOG
    )

    assert result == CommandResult(output='{"error":"UNKNOWN_COUPON"}', exit_code=1)


def test_prints_compact_json_on_a_single_line() -> None:
    result = run_checkout_command(ONE_APPLE, TINY_CATALOG)

    assert result.output.startswith('{"lines":[{"sku":"APPLE","name":"Apple","quantity":1,')
    assert result.output.endswith(',"shipping":499,"tax":4,"total":553}')
    assert "\n" not in result.output


def test_does_not_hide_bad_input_behind_an_error_code() -> None:
    # Not a basket at all: that is not one of the six error codes, so the exception passes
    # through (main() turns it into exit code 2).
    with pytest.raises(json.JSONDecodeError):
        run_checkout_command("not json", TINY_CATALOG)
    with pytest.raises(NotABasketError):
        run_checkout_command("{}", TINY_CATALOG)


# --- main ----------------------------------------------------------------------------------


def run_main(stdin_text: str, monkeypatch: pytest.MonkeyPatch) -> int:
    """Runs main() with `stdin_text` as standard input and returns its exit code."""
    # monkeypatch.setattr replaces sys.stdin for this test only; pytest restores it afterwards.
    monkeypatch.setattr(sys, "stdin", io.StringIO(stdin_text))
    return main(TINY_CATALOG)


def test_main_prints_the_receipt_and_returns_0(
    monkeypatch: pytest.MonkeyPatch, capsys: pytest.CaptureFixture[str]
) -> None:
    exit_code = run_main(ONE_APPLE + "\n", monkeypatch)

    captured = capsys.readouterr()
    assert exit_code == 0
    assert json.loads(captured.out)["total"] == 553
    assert captured.out.endswith("\n")
    assert captured.err == ""


def test_main_prints_the_error_and_returns_1(
    monkeypatch: pytest.MonkeyPatch, capsys: pytest.CaptureFixture[str]
) -> None:
    exit_code = run_main('{"items":[]}', monkeypatch)

    captured = capsys.readouterr()
    assert exit_code == 1
    assert captured.out == '{"error":"EMPTY_BASKET"}\n'
    assert captured.err == ""


@pytest.mark.parametrize(
    ("stdin_text", "complaint"),
    [
        ("this is not json", "Expecting value"),
        ("", "Expecting value"),
        ('{"hello": "world"}', "The input is not a basket: 'items' must be a list"),
        # Python refuses to turn more than 4300 digits into an int (a guard against very
        # slow conversions), so this quantity is not a number Python can read at all.
        (
            '{"items":[{"sku":"APPLE","quantity":' + "9" * 5000 + "}]}",
            "Exceeds the limit (4300 digits) for integer string conversion",
        ),
        # JSON nested this deeply is more than Python's JSON reader can follow.
        ("[" * 100_000 + "]" * 100_000, "maximum recursion depth exceeded"),
    ],
    ids=["not-json", "empty", "not-a-basket", "number-with-5000-digits", "nested-100000-deep"],
)
def test_main_returns_2_with_nothing_on_stdout_and_a_message_on_stderr_for_bad_input(
    stdin_text: str,
    complaint: str,
    monkeypatch: pytest.MonkeyPatch,
    capsys: pytest.CaptureFixture[str],
) -> None:
    exit_code = run_main(stdin_text, monkeypatch)

    captured = capsys.readouterr()
    assert exit_code == 2
    assert captured.out == ""
    assert captured.err.startswith("checkout: ")
    assert complaint in captured.err


def test_main_uses_the_catalog_it_is_given(
    monkeypatch: pytest.MonkeyPatch, capsys: pytest.CaptureFixture[str]
) -> None:
    # A TV is only in the tiny made-up catalog, not in the standard one.
    exit_code = run_main('{"items":[{"sku":"TV","quantity":1}]}', monkeypatch)

    assert exit_code == 0
    assert json.loads(capsys.readouterr().out)["subtotal"] == 10000


def test_main_uses_the_standard_catalog_when_none_is_given(
    monkeypatch: pytest.MonkeyPatch, capsys: pytest.CaptureFixture[str]
) -> None:
    # HEADPHONES only exist in the standard catalog, not in the tiny one.
    monkeypatch.setattr(sys, "stdin", io.StringIO('{"items":[{"sku":"HEADPHONES","quantity":1}]}'))

    exit_code = main()

    assert exit_code == 0
    assert json.loads(capsys.readouterr().out)["subtotal"] == 5999
