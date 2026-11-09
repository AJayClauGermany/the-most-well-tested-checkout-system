"""BDD tests with pytest-bdd: the Given/When/Then scenarios in spec/features/*.feature.

Behaviour-driven development (BDD) tests are written as plain-English scenarios that a
non-programmer can read and check. The scenarios are shared by every language and are read
straight from the repository's spec/features/ folder - never copied. What this file adds is
the Python "glue": one step definition for every sentence in spec/features/README.md.

`scenarios(...)` turns every scenario in those files into an ordinary pytest test, so the
usual tools (-k, -v, fixtures) all work. Each scenario gets a fresh `world` fixture that
carries state from step to step.
"""

from collections.abc import Sequence
from dataclasses import dataclass, field, replace
from datetime import UTC, datetime

import pytest
from pytest_bdd import given, parsers, scenarios, then, when
from spec_loader import SPEC_FOLDER, load_catalog
from support.dollars import Money

from checkout import (
    Basket,
    BasketItem,
    Catalog,
    CheckoutError,
    InMemoryOrderRepository,
    Receipt,
    checkout,
    place_order,
)

# Every scenario in every .feature file in the shared folder becomes a test in this module.
scenarios(str(SPEC_FOLDER / "features"))

# "{amount:Money}" in a step sentence matches "$12.39" and hands the step 1239 cents.
MONEY = {"Money": Money()}


@dataclass
class CheckoutWorld:
    """Everything one scenario knows: the basket being built, what came out of the checkout,
    and the (fake) order repository with its replaceable clock."""

    catalog: Catalog = field(default_factory=lambda: Catalog(products=[]))
    basket: Basket = field(default_factory=Basket)
    receipt: Receipt | None = None
    error: CheckoutError | None = None
    # What the repository's clock answers. "Given the clock is fixed at ..." changes it.
    now: datetime = datetime(2026, 1, 1, 0, 0, 0, tzinfo=UTC)
    repository: InMemoryOrderRepository | None = None
    last_order_id: str | None = None

    def receipt_or_fail(self) -> Receipt:
        assert self.error is None, f"checkout failed with {self.error.code}"
        assert self.receipt is not None, "the customer has not checked out yet"
        return self.receipt

    def repository_or_fail(self) -> InMemoryOrderRepository:
        assert self.repository is not None, (
            "no order repository - add 'Given an empty order repository'"
        )
        return self.repository


@pytest.fixture
def world() -> CheckoutWorld:
    return CheckoutWorld()


# --- Given ----------------------------------------------------------------------------------


@given("the standard catalog")
def the_standard_catalog(world: CheckoutWorld) -> None:
    world.catalog = load_catalog()


@given("an empty basket")
def an_empty_basket(world: CheckoutWorld) -> None:
    world.basket = Basket(items=[])


@given("a basket containing:")
def a_basket_containing(world: CheckoutWorld, datatable: Sequence[Sequence[object]]) -> None:
    # The data table arrives as rows of cells; the first row holds the column names.
    header, *rows = datatable
    for row in rows:
        cells = dict(zip(header, row, strict=True))
        world.basket.items.append(
            BasketItem(sku=str(cells["sku"]), quantity=int(str(cells["quantity"])))
        )


@given(parsers.parse('the customer uses the coupon "{code}"'))
def the_customer_uses_the_coupon(world: CheckoutWorld, code: str) -> None:
    # Basket is frozen, so `replace` makes a copy with one field changed.
    world.basket = replace(world.basket, coupon=code)


@given(parsers.parse('the customer chooses "{method}" shipping'))
def the_customer_chooses_shipping(world: CheckoutWorld, method: str) -> None:
    world.basket = replace(world.basket, shipping=method)


@given("an empty order repository")
def an_empty_order_repository(world: CheckoutWorld) -> None:
    # The clock reads `world.now`, so a later "the clock is fixed at" step still applies.
    world.repository = InMemoryOrderRepository(lambda: world.now)


@given(parsers.parse('the clock is fixed at "{timestamp}"'))
def the_clock_is_fixed_at(world: CheckoutWorld, timestamp: str) -> None:
    # "2026-01-15T10:30:00Z" -> an aware datetime. fromisoformat() accepts the Z since 3.11,
    # but we support 3.10, so the Z is spelt out as +00:00 first.
    world.now = datetime.fromisoformat(timestamp.replace("Z", "+00:00"))


# --- When -----------------------------------------------------------------------------------


@when("the customer checks out")
def the_customer_checks_out(world: CheckoutWorld) -> None:
    try:
        world.receipt = checkout(world.basket, world.catalog)
    except CheckoutError as error:
        world.error = error


@when("the customer checks out and saves the order")
def the_customer_checks_out_and_saves_the_order(world: CheckoutWorld) -> None:
    try:
        placed = place_order(world.basket, world.catalog, world.repository_or_fail())
    except CheckoutError as error:
        world.error = error
    else:
        world.receipt = placed.receipt
        world.last_order_id = placed.order_id


# --- Then: the receipt ----------------------------------------------------------------------


@then(parsers.parse("the subtotal is {cents:Money}", extra_types=MONEY))
def the_subtotal_is(world: CheckoutWorld, cents: int) -> None:
    assert world.receipt_or_fail().subtotal == cents


@then(parsers.parse("the promotion discount is {cents:Money}", extra_types=MONEY))
def the_promotion_discount_is(world: CheckoutWorld, cents: int) -> None:
    assert world.receipt_or_fail().promotion_discount == cents


@then(parsers.parse("the coupon discount is {cents:Money}", extra_types=MONEY))
def the_coupon_discount_is(world: CheckoutWorld, cents: int) -> None:
    assert world.receipt_or_fail().coupon_discount == cents


@then(parsers.parse("the shipping cost is {cents:Money}", extra_types=MONEY))
def the_shipping_cost_is(world: CheckoutWorld, cents: int) -> None:
    assert world.receipt_or_fail().shipping == cents


@then(parsers.parse("the tax is {cents:Money}", extra_types=MONEY))
def the_tax_is(world: CheckoutWorld, cents: int) -> None:
    assert world.receipt_or_fail().tax == cents


@then(parsers.parse("the total is {cents:Money}", extra_types=MONEY))
def the_total_is(world: CheckoutWorld, cents: int) -> None:
    assert world.receipt_or_fail().total == cents


@then(parsers.parse("the receipt has {count:d} lines"))
def the_receipt_has_lines(world: CheckoutWorld, count: int) -> None:
    assert len(world.receipt_or_fail().lines) == count


@then(parsers.parse('the receipt line for "{sku}" has a quantity of {quantity:d}'))
def the_receipt_line_for_has_a_quantity_of(world: CheckoutWorld, sku: str, quantity: int) -> None:
    lines = [line for line in world.receipt_or_fail().lines if line.sku == sku]
    assert len(lines) == 1, f"expected exactly one receipt line for {sku}, found {len(lines)}"
    assert lines[0].quantity == quantity


@then(parsers.parse('checkout fails with the error "{code}"'))
def checkout_fails_with_the_error(world: CheckoutWorld, code: str) -> None:
    assert world.receipt is None, "checkout succeeded but was expected to fail"
    assert world.error is not None, "checkout did not fail"
    assert world.error.code == code


# --- Then: saved orders ---------------------------------------------------------------------


@then(parsers.parse('the order is saved with the id "{order_id}"'))
def the_order_is_saved_with_the_id(world: CheckoutWorld, order_id: str) -> None:
    assert world.last_order_id == order_id


@then(parsers.parse('the saved order "{order_id}" has a total of {cents:Money}', extra_types=MONEY))
def the_saved_order_has_a_total_of(world: CheckoutWorld, order_id: str, cents: int) -> None:
    order = world.repository_or_fail().find(order_id)
    assert order is not None, f"order {order_id} not found"
    assert order.receipt.total == cents


@then("the order repository is empty")
def the_order_repository_is_empty(world: CheckoutWorld) -> None:
    assert world.repository_or_fail().count() == 0


@then(parsers.parse('the order "{order_id}" cannot be found'))
def the_order_cannot_be_found(world: CheckoutWorld, order_id: str) -> None:
    assert world.repository_or_fail().find(order_id) is None


@then(parsers.parse('the saved order "{order_id}" was created at "{timestamp}"'))
def the_saved_order_was_created_at(world: CheckoutWorld, order_id: str, timestamp: str) -> None:
    order = world.repository_or_fail().find(order_id)
    assert order is not None, f"order {order_id} not found"
    # Timestamps are compared as strings, exactly as spec/features/README.md says.
    assert order.created_at == timestamp
