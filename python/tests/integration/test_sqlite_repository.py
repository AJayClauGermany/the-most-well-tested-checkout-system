"""Integration tests for the SQLite order repository, using the standard library's sqlite3.

A real database, but no server to install: SQLite is one file. Each test gets its own
file in pytest's temporary folder (`tmp_path`) and closes the database afterwards.
"""

import sqlite3
from collections.abc import Iterator
from datetime import datetime, timezone
from pathlib import Path

import pytest
from support.fixtures import ONE_APPLE_RECEIPT, TINY_CATALOG, fixed_clock
from support.schemas import schema_errors

from checkout import Basket, BasketItem, Order, SqliteOrderRepository, checkout


@pytest.fixture
def database(tmp_path: Path) -> Path:
    """The path of a database file that does not exist yet."""
    return tmp_path / "orders.db"


@pytest.fixture
def repository(database: Path) -> Iterator[SqliteOrderRepository]:
    """An open repository with the stub clock. The `yield` hands it to the test; what
    follows the `yield` runs afterwards, like a tidy-up step - here, closing the database."""
    repository = SqliteOrderRepository(database, fixed_clock)
    yield repository
    repository.close()


def test_saves_and_finds_an_order(repository: SqliteOrderRepository) -> None:
    receipt = checkout(Basket(items=[BasketItem("COFFEE", 3)], coupon="SAVE10"), TINY_CATALOG)

    identifier = repository.save(receipt)

    assert identifier == "ORDER-1"
    assert repository.find(identifier) == Order(
        id="ORDER-1", created_at="2026-01-15T10:30:00Z", receipt=receipt
    )


def test_gives_back_an_order_that_follows_order_schema_json(
    repository: SqliteOrderRepository,
) -> None:
    identifier = repository.save(ONE_APPLE_RECEIPT)

    found = repository.find(identifier)

    assert found is not None
    assert schema_errors("order.schema.json", found.to_dict()) == []


def test_numbers_orders_sequentially(repository: SqliteOrderRepository) -> None:
    ids = [
        repository.save(ONE_APPLE_RECEIPT),
        repository.save(ONE_APPLE_RECEIPT),
        repository.save(ONE_APPLE_RECEIPT),
    ]

    assert ids == ["ORDER-1", "ORDER-2", "ORDER-3"]


def test_continues_numbering_after_the_database_is_closed_and_opened_again(
    database: Path,
) -> None:
    with SqliteOrderRepository(database, fixed_clock) as first:
        first.save(ONE_APPLE_RECEIPT)
        first.save(ONE_APPLE_RECEIPT)

    # The `with` block closed the first connection; this is a brand new one on the same file.
    with SqliteOrderRepository(database, fixed_clock) as second:
        assert second.save(ONE_APPLE_RECEIPT) == "ORDER-3"
        found = second.find("ORDER-2")
        assert found is not None
        assert found.id == "ORDER-2"


def test_really_writes_to_the_database_file_checked_with_a_separate_connection(
    database: Path,
) -> None:
    with SqliteOrderRepository(database, fixed_clock) as repository:
        repository.save(ONE_APPLE_RECEIPT)

    connection = sqlite3.connect(database)
    try:
        (count,) = connection.execute("SELECT COUNT(*) FROM orders").fetchone()
    finally:
        connection.close()

    assert count == 1


def test_refuses_to_save_or_find_once_it_is_closed_instead_of_silently_doing_nothing(
    database: Path,
) -> None:
    repository = SqliteOrderRepository(database, fixed_clock)
    repository.close()

    with pytest.raises(sqlite3.ProgrammingError, match="closed database"):
        repository.save(ONE_APPLE_RECEIPT)
    with pytest.raises(sqlite3.ProgrammingError, match="closed database"):
        repository.find("ORDER-1")


def test_works_in_memory_too() -> None:
    # ":memory:" is SQLite's name for a throw-away database that lives only in RAM.
    with SqliteOrderRepository(":memory:", fixed_clock) as repository:
        assert repository.save(ONE_APPLE_RECEIPT) == "ORDER-1"
        found = repository.find("ORDER-1")
        assert found is not None
        assert found.receipt == ONE_APPLE_RECEIPT


@pytest.mark.parametrize(
    "identifier",
    ["ORDER-2", "ORDER-0", "ORDER-01", "ORDER-x", "1", "", "ORDER-1; DROP TABLE orders"],
    ids=repr,
)
def test_answers_not_found_without_raising(
    repository: SqliteOrderRepository, identifier: str
) -> None:
    repository.save(ONE_APPLE_RECEIPT)

    assert repository.find(identifier) is None
    assert repository.find("ORDER-1") is not None  # and the table is still there


def test_stamps_each_order_with_the_clock_s_time_at_the_moment_of_saving(
    database: Path,
) -> None:
    now = datetime(2026, 1, 15, 10, 30, 0, tzinfo=timezone.utc)
    with SqliteOrderRepository(database, lambda: now) as repository:
        first = repository.save(ONE_APPLE_RECEIPT)
        now = datetime(2026, 1, 16, 9, 0, 0, tzinfo=timezone.utc)
        second = repository.save(ONE_APPLE_RECEIPT)

        first_order = repository.find(first)
        second_order = repository.find(second)

    assert first_order is not None and first_order.created_at == "2026-01-15T10:30:00Z"
    assert second_order is not None and second_order.created_at == "2026-01-16T09:00:00Z"


def test_uses_the_real_clock_when_none_is_given(database: Path) -> None:
    with SqliteOrderRepository(database) as repository:
        found = repository.find(repository.save(ONE_APPLE_RECEIPT))

    assert found is not None
    assert found.created_at.startswith(str(datetime.now(timezone.utc).year))


def test_the_with_block_closes_the_database_even_when_the_body_fails(database: Path) -> None:
    with pytest.raises(RuntimeError, match="boom"):
        with SqliteOrderRepository(database, fixed_clock) as repository:
            raise RuntimeError("boom")

    with pytest.raises(sqlite3.ProgrammingError, match="closed database"):
        repository.find("ORDER-1")
