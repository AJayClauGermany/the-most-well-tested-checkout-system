"""Snapshot tests with syrupy: did the output change since last time?

The `snapshot` fixture saves a value to tests/snapshot/__snapshots__/ the first time a test
runs (`pytest --snapshot-update`) and compares against that file on every later run. When
a change is intended, run `pytest --snapshot-update` again and commit the new file - the
diff in the pull request shows exactly what changed.

The hand-checked receipt in spec/examples/receipt.json plays the role of an "approved"
output shared by all six languages; the first test compares against it directly.
"""

import json
from dataclasses import asdict

from spec_loader import SPEC_FOLDER, load_catalog
from syrupy.assertion import SnapshotAssertion
from syrupy.extensions.json import JSONSnapshotExtension

from checkout import Basket, BasketItem, checkout, run_checkout_command

CATALOG = load_catalog()

# The basket behind spec/examples/receipt.json (rules §5), as an object and as JSON text.
EXAMPLE_BASKET = Basket(
    items=[BasketItem("COFFEE", 3), BasketItem("APPLE", 2)], coupon="SAVE10", shipping="standard"
)
EXAMPLE_BASKET_JSON = (
    '{"items": [{"sku": "COFFEE", "quantity": 3}, {"sku": "APPLE", "quantity": 2}],'
    ' "coupon": "SAVE10", "shipping": "standard"}'
)

# The biggest shared case: every product, both promotions, WELCOME5 and express shipping.
FULL_EVERYTHING_BASKET = Basket(
    items=[
        BasketItem("HEADPHONES", 1),
        BasketItem("COFFEE", 6),
        BasketItem("NOTEBOOK", 12),
        BasketItem("PEN", 5),
        BasketItem("BREAD", 1),
        BasketItem("APPLE", 3),
    ],
    coupon="WELCOME5",
    shipping="express",
)


def test_the_example_receipt_matches_the_hand_checked_spec_examples_receipt_json() -> None:
    # An "approval test" against a file that humans checked by hand, shared by every language.
    approved = json.loads((SPEC_FOLDER / "examples" / "receipt.json").read_text(encoding="utf-8"))

    assert checkout(EXAMPLE_BASKET, CATALOG).to_dict() == approved


def test_the_example_receipt_matches_the_saved_snapshot(snapshot: SnapshotAssertion) -> None:
    # The JSON extension stores the snapshot as a readable .json file instead of syrupy's
    # default Python-like format - handy for a receipt that IS JSON.
    assert checkout(EXAMPLE_BASKET, CATALOG).to_dict() == snapshot(
        extension_class=JSONSnapshotExtension
    )


def test_the_example_receipt_is_printed_by_the_command_exactly_like_this(
    snapshot: SnapshotAssertion,
) -> None:
    result = run_checkout_command(EXAMPLE_BASKET_JSON, CATALOG)

    assert result.output == snapshot


def test_the_full_everything_receipt_matches_the_saved_snapshot(
    snapshot: SnapshotAssertion,
) -> None:
    assert checkout(FULL_EVERYTHING_BASKET, CATALOG).to_dict() == snapshot(
        extension_class=JSONSnapshotExtension
    )


def test_one_pen_written_out_as_receipt_lines(snapshot: SnapshotAssertion) -> None:
    receipt = checkout(Basket(items=[BasketItem("PEN", 1)]), CATALOG)

    assert receipt.to_dict()["lines"] == snapshot(extension_class=JSONSnapshotExtension)


def test_an_error_as_the_command_prints_it(snapshot: SnapshotAssertion) -> None:
    result = run_checkout_command('{"items":[]}', CATALOG)

    assert asdict(result) == snapshot
