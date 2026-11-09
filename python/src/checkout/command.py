"""The command line program's logic (rules §8), kept free of the terminal so it can be unit
tested in-process. ``src/checkout/__main__.py`` is the thin wrapper that calls ``main()``.
"""

from __future__ import annotations

import json
import sys
from dataclasses import dataclass
from pathlib import Path
from typing import Literal

from checkout.engine import checkout
from checkout.errors import CheckoutError
from checkout.models import Basket, Catalog
from checkout.spec_folder import find_spec_folder


def standard_catalog_path() -> Path:
    """Where the standard catalog lives: spec/catalog.json at the root of the repository."""
    return find_spec_folder() / "catalog.json"


def load_standard_catalog() -> Catalog:
    """Reads spec/catalog.json."""
    # json.loads accepts bytes and works out the encoding (UTF-8) by itself.
    return Catalog.from_dict(json.loads(standard_catalog_path().read_bytes()))


class NotABasketError(ValueError):
    """The input was valid JSON, but not the shape of a basket (rules §3)."""

    def __init__(self, reason: str) -> None:
        super().__init__(f"The input is not a basket: {reason}")


def _is_whole_number(value: object) -> bool:
    # JSON `true` arrives as Python's True, which is an int in disguise - not a quantity.
    return isinstance(value, int) and not isinstance(value, bool)


def parse_basket(text: str) -> Basket:
    """Parses the JSON text of a basket.

    Only the *shape* is checked here (the types that basket.schema.json describes), so that
    ``checkout()`` can do the real validation: a quantity of 0 is a basket, just not a valid
    one. Raises ``json.JSONDecodeError`` when the text is not JSON, and ``NotABasketError``
    when the JSON is not a basket - or is more than Python's JSON reader can handle (a number
    with more than 4300 digits, or lists nested thousands of levels deep).
    """
    try:
        parsed = json.loads(text)
    except json.JSONDecodeError:
        raise
    except (ValueError, RecursionError) as error:
        raise NotABasketError(str(error)) from error
    if not isinstance(parsed, dict):
        raise NotABasketError("expected a JSON object")
    items = parsed.get("items")
    if not isinstance(items, list):
        raise NotABasketError("'items' must be a list")
    for item in items:
        if not isinstance(item, dict):
            raise NotABasketError("every item must be an object")
        if not isinstance(item.get("sku"), str):
            raise NotABasketError("every item needs a string 'sku'")
        if not _is_whole_number(item.get("quantity")):
            raise NotABasketError("every item needs a whole-number 'quantity'")
    if not isinstance(parsed.get("coupon"), str | None):
        raise NotABasketError("'coupon' must be a string or null")
    if not isinstance(parsed.get("shipping", "standard"), str):
        raise NotABasketError("'shipping' must be a string")
    return Basket.from_dict(parsed)


def compact_json(value: object) -> str:
    """JSON on one line with no spaces, the way the program prints it.

    >>> compact_json({"error": "EMPTY_BASKET", "n": [1, 2]})
    '{"error":"EMPTY_BASKET","n":[1,2]}'
    """
    return json.dumps(value, separators=(",", ":"))


@dataclass(frozen=True)
class CommandResult:
    """What the program prints and the exit code it ends with."""

    output: str
    exit_code: Literal[0, 1]


def run_checkout_command(text: str, catalog: Catalog) -> CommandResult:
    """Runs one checkout from the JSON text of a basket.

    - a valid basket gives the receipt as JSON and exit code 0,
    - an invalid basket gives ``{"error": "<CODE>"}`` and exit code 1,
    - text that is not a basket at all raises (``main()`` turns that into exit code 2).

    The JSON is printed compactly (no spaces), on one line.
    """
    basket = parse_basket(text)
    try:
        receipt = checkout(basket, catalog)
    except CheckoutError as error:
        return CommandResult(compact_json({"error": error.code.value}), 1)
    return CommandResult(compact_json(receipt.to_dict()), 0)


def main(catalog: Catalog | None = None) -> int:
    """The whole program: stdin -> run_checkout_command -> stdout, returning the exit code.

    ``catalog`` defaults to the standard catalog; tests pass a smaller one. Input that is
    not a basket at all gives exit code 2, nothing on stdout and a message on stderr.
    """
    text = sys.stdin.read()
    try:
        result = run_checkout_command(
            text, catalog if catalog is not None else load_standard_catalog()
        )
    except (json.JSONDecodeError, NotABasketError) as error:
        print(f"checkout: {error}", file=sys.stderr)
        return 2
    print(result.output)
    return result.exit_code
