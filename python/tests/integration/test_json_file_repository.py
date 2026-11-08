"""Integration tests for the JSON-file order repository: real files in a real temporary folder.

Unlike a unit test, this touches the file system. pytest's `tmp_path` fixture gives each
test its own fresh folder and deletes it afterwards. The stored files must follow
spec/schema/order.schema.json - a contract test inside an integration test.
"""

import json
from datetime import datetime, timezone
from pathlib import Path

import pytest
from support.fixtures import ONE_APPLE_RECEIPT, TINY_CATALOG, fixed_clock
from support.schemas import schema_errors

from checkout import Basket, BasketItem, JsonFileOrderRepository, Order, checkout


def test_writes_order_1_json_into_the_folder_on_the_first_save(tmp_path: Path) -> None:
    repository = JsonFileOrderRepository(tmp_path, fixed_clock)

    identifier = repository.save(ONE_APPLE_RECEIPT)

    assert identifier == "ORDER-1"
    assert (tmp_path / "ORDER-1.json").is_file()


def test_writes_a_file_that_follows_order_schema_json_exactly(tmp_path: Path) -> None:
    repository = JsonFileOrderRepository(tmp_path, fixed_clock)
    repository.save(ONE_APPLE_RECEIPT)

    stored = json.loads((tmp_path / "ORDER-1.json").read_text(encoding="utf-8"))

    assert schema_errors("order.schema.json", stored) == []
    assert stored == {
        "id": "ORDER-1",
        "created_at": "2026-01-15T10:30:00Z",
        "receipt": ONE_APPLE_RECEIPT.to_dict(),
    }


def test_writes_human_readable_json_indented_ending_in_a_newline(tmp_path: Path) -> None:
    repository = JsonFileOrderRepository(tmp_path, fixed_clock)
    repository.save(ONE_APPLE_RECEIPT)

    text = (tmp_path / "ORDER-1.json").read_text(encoding="utf-8")

    assert "\n  " in text
    assert text.endswith("\n")


def test_reads_a_saved_order_back_exactly_as_saved(tmp_path: Path) -> None:
    repository = JsonFileOrderRepository(tmp_path, fixed_clock)
    receipt = checkout(Basket(items=[BasketItem("COFFEE", 3)], coupon="SAVE10"), TINY_CATALOG)

    identifier = repository.save(receipt)

    assert repository.find(identifier) == Order(
        id="ORDER-1", created_at="2026-01-15T10:30:00Z", receipt=receipt
    )


def test_numbers_orders_sequentially(tmp_path: Path) -> None:
    repository = JsonFileOrderRepository(tmp_path, fixed_clock)

    ids = [
        repository.save(ONE_APPLE_RECEIPT),
        repository.save(ONE_APPLE_RECEIPT),
        repository.save(ONE_APPLE_RECEIPT),
    ]

    assert ids == ["ORDER-1", "ORDER-2", "ORDER-3"]
    assert sorted(path.name for path in tmp_path.iterdir()) == [
        "ORDER-1.json",
        "ORDER-2.json",
        "ORDER-3.json",
    ]


def test_continues_numbering_when_the_same_folder_is_opened_again_later(tmp_path: Path) -> None:
    first = JsonFileOrderRepository(tmp_path, fixed_clock)
    first.save(ONE_APPLE_RECEIPT)
    first.save(ONE_APPLE_RECEIPT)

    second = JsonFileOrderRepository(tmp_path, fixed_clock)

    assert second.save(ONE_APPLE_RECEIPT) == "ORDER-3"
    found = second.find("ORDER-1")
    assert found is not None
    assert found.id == "ORDER-1"


def write_order_file(folder: Path, identifier: str) -> None:
    order = Order(id=identifier, created_at="2026-01-15T10:30:00Z", receipt=ONE_APPLE_RECEIPT)
    (folder / f"{identifier}.json").write_text(json.dumps(order.to_dict()), encoding="utf-8")


def test_numbers_after_the_highest_existing_order_even_with_gaps_and_2_digit_numbers(
    tmp_path: Path,
) -> None:
    write_order_file(tmp_path, "ORDER-3")
    write_order_file(tmp_path, "ORDER-12")
    repository = JsonFileOrderRepository(tmp_path, fixed_clock)

    # Not ORDER-4: alphabetically "ORDER-3" > "ORDER-12", but 12 is the higher number.
    assert repository.save(ONE_APPLE_RECEIPT) == "ORDER-13"


def test_ignores_files_in_the_folder_that_are_not_orders(tmp_path: Path) -> None:
    (tmp_path / "README.txt").write_text("not an order")
    (tmp_path / "ORDER-9.json.bak").write_text("not an order either")
    (tmp_path / "ORDER-9").write_text("no .json extension, so not an order")
    (tmp_path / "ORDER-09.json").write_text("a leading zero is not an order id")
    (tmp_path / "ORDER-7.json").mkdir()  # a folder with a tempting name
    repository = JsonFileOrderRepository(tmp_path, fixed_clock)

    assert repository.save(ONE_APPLE_RECEIPT) == "ORDER-1"


def test_creates_the_folder_if_it_does_not_exist_yet(tmp_path: Path) -> None:
    nested = tmp_path / "orders" / "2026"
    repository = JsonFileOrderRepository(nested, fixed_clock)

    repository.save(ONE_APPLE_RECEIPT)

    assert (nested / "ORDER-1.json").is_file()


@pytest.mark.parametrize("identifier", ["ORDER-2", "ORDER-0", "nonsense", ""], ids=repr)
def test_answers_not_found_for_ids_that_have_no_file_without_raising(
    tmp_path: Path, identifier: str
) -> None:
    repository = JsonFileOrderRepository(tmp_path, fixed_clock)
    repository.save(ONE_APPLE_RECEIPT)

    assert repository.find(identifier) is None


def test_does_not_let_an_id_escape_the_folder_no_path_tricks(tmp_path: Path) -> None:
    (tmp_path / "secret.json").write_text('{"id": "x", "created_at": "x", "receipt": {}}')
    repository = JsonFileOrderRepository(tmp_path / "orders", fixed_clock)

    assert repository.find("../secret") is None
    assert repository.find("../secret.json") is None


def test_stamps_each_order_with_the_clock_s_time_at_the_moment_of_saving(tmp_path: Path) -> None:
    # A clock that can be moved forward from the test: it reads whatever `now` holds.
    now = datetime(2026, 1, 15, 10, 30, 0, tzinfo=timezone.utc)
    repository = JsonFileOrderRepository(tmp_path, lambda: now)

    first = repository.save(ONE_APPLE_RECEIPT)
    now = datetime(2026, 1, 16, 9, 0, 0, tzinfo=timezone.utc)
    second = repository.save(ONE_APPLE_RECEIPT)

    first_order = repository.find(first)
    second_order = repository.find(second)
    assert first_order is not None and first_order.created_at == "2026-01-15T10:30:00Z"
    assert second_order is not None and second_order.created_at == "2026-01-16T09:00:00Z"


def test_uses_the_real_clock_when_none_is_given(tmp_path: Path) -> None:
    repository = JsonFileOrderRepository(tmp_path)

    found = repository.find(repository.save(ONE_APPLE_RECEIPT))

    assert found is not None
    assert found.created_at.startswith(str(datetime.now(timezone.utc).year))
