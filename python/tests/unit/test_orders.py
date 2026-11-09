"""Unit tests for saving orders (rules §6), and the project's showcase of TEST DOUBLES.

A *test double* stands in for a real dependency:
  * a FAKE is a simple working version - InMemoryOrderRepository is one,
  * a STUB gives canned answers - `fixed_clock` always says 10:30 on 15 January 2026,
  * a MOCK records how it was called - unittest.mock.Mock lets us check that save()
    was never called when checkout fails.
pytest's `monkeypatch` fixture is a fourth tool: it swaps a real function for a fake one
for the duration of one test, and puts the original back afterwards.
"""

import time
from datetime import datetime, timedelta, timezone
from unittest.mock import MagicMock, Mock, create_autospec

import pytest
from support.fixtures import FIXED_TIME, ONE_APPLE_RECEIPT, TINY_CATALOG, fixed_clock

from checkout import (
    Basket,
    BasketItem,
    CheckoutError,
    ErrorCode,
    InMemoryOrderRepository,
    Order,
    OrderRepository,
    PlacedOrder,
    ReceiptLine,
    checkout,
    format_timestamp,
    order_id,
    order_number,
    place_order,
    system_clock,
)
from checkout import orders as orders_module

# --- format_timestamp -----------------------------------------------------------------


def test_writes_utc_with_whole_seconds_and_a_trailing_z() -> None:
    assert format_timestamp(FIXED_TIME) == "2026-01-15T10:30:00Z"


def test_drops_the_microseconds_instead_of_rounding_them() -> None:
    almost_next_second = datetime(2026, 1, 15, 10, 30, 0, 999_999, tzinfo=timezone.utc)

    assert format_timestamp(almost_next_second) == "2026-01-15T10:30:00Z"


def test_converts_other_time_zones_to_utc() -> None:
    berlin = timezone(timedelta(hours=2))
    half_past_twelve_in_berlin = datetime(2026, 1, 15, 12, 30, 0, tzinfo=berlin)

    assert format_timestamp(half_past_twelve_in_berlin) == "2026-01-15T10:30:00Z"


@pytest.mark.skipif(not hasattr(time, "tzset"), reason="time.tzset() only exists on Unix")
def test_converts_to_utc_even_when_the_computer_is_set_to_another_time_zone(
    monkeypatch: pytest.MonkeyPatch,
) -> None:
    # This test machine is probably set to UTC, which would hide a bug that converts to the
    # computer's LOCAL time instead of UTC. So: pretend to be in Japan for one test.
    # "JST-9" is a POSIX time zone string (9 hours ahead of UTC) that needs no tz database.
    monkeypatch.setenv("TZ", "JST-9")
    time.tzset()
    try:
        assert format_timestamp(FIXED_TIME) == "2026-01-15T10:30:00Z"
    finally:
        monkeypatch.undo()  # put the real TZ back ...
        time.tzset()  # ... and tell the C library about it


def test_refuses_a_datetime_without_a_time_zone() -> None:
    # A "naive" datetime could be in any time zone, so there is no right answer.
    with pytest.raises(ValueError, match="time zone"):
        format_timestamp(datetime(2026, 1, 15, 10, 30, 0))


# --- order_id and order_number ----------------------------------------------------


@pytest.mark.parametrize(
    ("number", "identifier"),
    [(1, "ORDER-1"), (2, "ORDER-2"), (10, "ORDER-10"), (1234, "ORDER-1234")],
    ids=["1", "2", "10", "1234"],
)
def test_number_and_id_convert_both_ways(number: int, identifier: str) -> None:
    assert order_id(number) == identifier
    assert order_number(identifier) == number


@pytest.mark.parametrize(
    "identifier",
    [
        "",
        "ORDER-",
        "ORDER-0",
        "ORDER-01",
        "ORDER--1",
        "ORDER-1.5",
        "ORDER-1 ",
        " ORDER-1",
        "MYORDER-1",
        "ORDER-1x",
        "order-1",
        "ORDER-x",
        "1",
        "ORDER_1",
        "ORDER-1\n",
    ],
    ids=repr,
)
def test_does_not_recognise_anything_that_is_not_an_order_id(identifier: str) -> None:
    assert order_number(identifier) is None


# --- system_clock ---------------------------------------------------------------------


def test_the_real_clock_answers_the_current_time_in_utc() -> None:
    before = datetime.now(timezone.utc)

    now = system_clock()

    assert now.tzinfo == timezone.utc
    assert before <= now <= datetime.now(timezone.utc)


# --- Order -----------------------------------------------------------------------------


def test_an_order_can_be_turned_into_json_and_back() -> None:
    order = Order(id="ORDER-1", created_at="2026-01-15T10:30:00Z", receipt=ONE_APPLE_RECEIPT)

    as_json = order.to_dict()

    assert as_json == {
        "id": "ORDER-1",
        "created_at": "2026-01-15T10:30:00Z",
        "receipt": ONE_APPLE_RECEIPT.to_dict(),
    }
    assert Order.from_dict(as_json) == order


# --- InMemoryOrderRepository (a hand-written FAKE) ------------------------------------


def test_starts_empty() -> None:
    repository = InMemoryOrderRepository(fixed_clock)

    assert repository.count() == 0
    assert repository.find("ORDER-1") is None


def test_numbers_orders_1_2_3_in_the_order_they_are_saved() -> None:
    repository = InMemoryOrderRepository(fixed_clock)

    ids = [
        repository.save(ONE_APPLE_RECEIPT),
        repository.save(ONE_APPLE_RECEIPT),
        repository.save(ONE_APPLE_RECEIPT),
    ]

    assert ids == ["ORDER-1", "ORDER-2", "ORDER-3"]
    assert repository.count() == 3


def test_finds_a_saved_order_with_its_receipt_and_the_time_the_stub_clock_said() -> None:
    repository = InMemoryOrderRepository(fixed_clock)

    identifier = repository.save(ONE_APPLE_RECEIPT)

    assert repository.find(identifier) == Order(
        id="ORDER-1", created_at="2026-01-15T10:30:00Z", receipt=ONE_APPLE_RECEIPT
    )


def test_asks_the_clock_at_the_moment_of_saving_once_per_save() -> None:
    # A MOCK clock: Mock(return_value=...) answers like the stub, and remembers every call.
    clock = Mock(return_value=FIXED_TIME)
    repository = InMemoryOrderRepository(clock)

    assert clock.call_count == 0  # creating the repository does not look at the clock
    repository.save(ONE_APPLE_RECEIPT)
    repository.save(ONE_APPLE_RECEIPT)

    assert clock.call_count == 2
    clock.assert_called_with()  # ... and always without arguments


def test_keeps_a_copy_so_changing_the_original_later_does_not_change_the_order() -> None:
    repository = InMemoryOrderRepository(fixed_clock)
    receipt = checkout(Basket(items=[BasketItem("APPLE", 1)]), TINY_CATALOG)

    identifier = repository.save(receipt)
    # Receipt is a frozen dataclass, but its list of lines is still an ordinary list.
    receipt.lines.append(
        ReceiptLine(sku="PEN", name="Pen", quantity=9, unit_price=1, line_total=9, line_discount=0)
    )

    found = repository.find(identifier)
    assert found is not None
    assert found.receipt == ONE_APPLE_RECEIPT


@pytest.mark.parametrize("identifier", ["ORDER-2", "ORDER-0", "ORDER-99", "order-1", ""], ids=repr)
def test_answers_not_found_for_ids_that_were_never_handed_out(identifier: str) -> None:
    repository = InMemoryOrderRepository(fixed_clock)
    repository.save(ONE_APPLE_RECEIPT)

    assert repository.find(identifier) is None


def test_uses_the_real_clock_when_none_is_given(monkeypatch: pytest.MonkeyPatch) -> None:
    # monkeypatch swaps the module's real clock for a stub, for this test only.
    monkeypatch.setattr(
        orders_module, "system_clock", lambda: datetime(2031, 2, 3, 4, 5, 6, 789, timezone.utc)
    )
    repository = InMemoryOrderRepository()

    identifier = repository.save(ONE_APPLE_RECEIPT)

    found = repository.find(identifier)
    assert found is not None
    assert found.created_at == "2031-02-03T04:05:06Z"


def test_without_monkeypatch_the_default_clock_is_roughly_now() -> None:
    repository = InMemoryOrderRepository()

    identifier = repository.save(ONE_APPLE_RECEIPT)

    found = repository.find(identifier)
    assert found is not None
    # Checking the year is enough: this test must never be flaky (rules §6 explains why
    # the fixed clock exists); the exact format is checked with the stub above.
    assert found.created_at.startswith(str(datetime.now(timezone.utc).year))
    assert found.created_at.endswith("Z")


def test_is_an_order_repository() -> None:
    # Protocols are checked by shape, not by inheritance: anything with save() and find()
    # fits. The type checker (mypy) proves it at check time; this line proves it at run time.
    repository: OrderRepository = InMemoryOrderRepository(fixed_clock)

    assert repository.find("ORDER-1") is None


# --- place_order ---------------------------------------------------------------------------


def test_checks_out_and_saves_the_receipt_returning_both_the_receipt_and_the_new_id() -> None:
    repository = InMemoryOrderRepository(fixed_clock)

    placed = place_order(Basket(items=[BasketItem("APPLE", 1)]), TINY_CATALOG, repository)

    assert placed == PlacedOrder(order_id="ORDER-1", receipt=ONE_APPLE_RECEIPT)
    found = repository.find("ORDER-1")
    assert found is not None
    assert found.receipt == placed.receipt


def test_saves_exactly_what_it_returns_checked_with_a_mock_repository() -> None:
    # create_autospec builds a MOCK with the same methods as the real class, so a typo
    # like repository.svae(...) fails loudly instead of silently recording a call.
    repository = create_autospec(InMemoryOrderRepository, instance=True)
    repository.save.return_value = "ORDER-7"

    placed = place_order(Basket(items=[BasketItem("TV", 1)]), TINY_CATALOG, repository)

    repository.save.assert_called_once_with(placed.receipt)
    assert placed.order_id == "ORDER-7"  # the id comes from the repository, not from place_order
    repository.find.assert_not_called()


def test_saves_nothing_when_checkout_fails_p14() -> None:
    # MagicMock is Mock plus support for Python's "magic" methods (len(), ==, with, ...);
    # for a repository either does, and both can prove a call never happened.
    repository = MagicMock(spec=InMemoryOrderRepository)

    with pytest.raises(CheckoutError) as empty:
        place_order(Basket(items=[]), TINY_CATALOG, repository)
    with pytest.raises(CheckoutError) as unknown:
        place_order(Basket(items=[BasketItem("BANANA", 1)]), TINY_CATALOG, repository)

    assert empty.value.code == ErrorCode.EMPTY_BASKET
    assert unknown.value.code == ErrorCode.UNKNOWN_PRODUCT
    repository.save.assert_not_called()


def test_leaves_the_fake_repository_empty_after_a_failed_checkout() -> None:
    repository = InMemoryOrderRepository(fixed_clock)

    with pytest.raises(CheckoutError):
        place_order(Basket(items=[BasketItem("APPLE", 0)]), TINY_CATALOG, repository)

    assert repository.count() == 0
    assert repository.find("ORDER-1") is None
