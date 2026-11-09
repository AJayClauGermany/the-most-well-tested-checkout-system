"""Property-based tests with Hypothesis (https://hypothesis.readthedocs.io/).

An example test says "THIS basket gives THAT total". A property says something that must be
true for EVERY basket, and Hypothesis generates hundreds of random baskets trying to break
it - then shrinks any failure to the smallest example. The properties are listed once for
all six languages in spec/properties.md; each test here is named after its id (P1 ... P14).

`@given(...)` hands the test random values from the strategies below. `assume(...)` throws
away examples that do not apply (an invalid basket for a property about valid ones).
"""

from datetime import UTC, datetime

import pytest
from hypothesis import assume, given, settings
from hypothesis import strategies as st
from spec_loader import load_catalog
from support.schemas import schema_errors

from checkout import (
    Basket,
    BasketItem,
    CheckoutError,
    ErrorCode,
    InMemoryOrderRepository,
    Receipt,
    checkout,
    normalize_coupon_code,
    place_order,
    round_half_up,
)

CATALOG = load_catalog()
SKUS = [product.sku for product in CATALOG.products]
REAL_COUPONS = ["SAVE10", "WELCOME5"]

# ---------------------------------------------------------------------------
# Strategies (the generators), following "Generating baskets" in spec/properties.md
# ---------------------------------------------------------------------------

# A sensible item: a real SKU, a quantity from 1 to 99.
items = st.builds(BasketItem, sku=st.sampled_from(SKUS), quantity=st.integers(1, 99))


def sloppy(code: str) -> st.SearchStrategy[str]:
    """`code` written any old way: mixed case, spaces around it ("  sAvE10 ")."""
    flags = st.lists(st.booleans(), min_size=len(code), max_size=len(code))
    spaces = st.integers(0, 3)
    return st.tuples(flags, spaces, spaces).map(
        lambda parts: (
            " " * parts[1]
            + "".join(c.upper() if up else c.lower() for c, up in zip(code, parts[0], strict=True))
            + " " * parts[2]
        )
    )


sloppy_coupons = st.sampled_from(REAL_COUPONS).flatmap(sloppy)

# Coupons that can check out: none, blank, a real code, or a real code written sloppily.
coupons = st.one_of(
    st.none(), st.just(""), st.just("  "), st.sampled_from(REAL_COUPONS), sloppy_coupons
)
# Shipping methods that can check out. (A missing field means "standard", so None -> default.)
shippings = st.sampled_from(["standard", "express"])


def baskets(
    coupon: st.SearchStrategy[str | None] = coupons,
    shipping: st.SearchStrategy[str] = shippings,
) -> st.SearchStrategy[Basket]:
    """A whole basket: 1 to 8 items (the same SKU may repeat), a coupon, a shipping method."""
    return st.builds(
        Basket, items=st.lists(items, min_size=1, max_size=8), coupon=coupon, shipping=shipping
    )


def checks_out(basket: Basket) -> bool:
    try:
        checkout(basket, CATALOG)
    except CheckoutError:
        return False
    return True


# Only baskets that check out successfully - what P1 to P9 are about. (The only way one of
# the baskets above can fail is WELCOME5 below $25.00, so little is filtered out.)
valid_baskets = baskets().filter(checks_out)

# Anything at all that still fits the Basket type: nonsense SKUs, silly quantities of the
# wrong type, random strings for coupon and shipping, even no items. For P10 and P14.
any_items = st.builds(
    BasketItem,
    sku=st.one_of(st.sampled_from(SKUS), st.text(max_size=10)),
    quantity=st.one_of(
        st.integers(-5, 150),
        st.floats(),
        st.text(max_size=3),
        st.booleans(),
        st.none(),
    ),
)
any_baskets = st.builds(
    Basket,
    items=st.lists(any_items, max_size=6),
    coupon=st.one_of(
        st.none(), st.sampled_from(REAL_COUPONS), sloppy_coupons, st.text(max_size=12)
    ),
    shipping=st.one_of(shippings, st.text(max_size=12)),
)

# Any real moment this century, with a time zone (a clock always has one).
moments = st.datetimes(
    min_value=datetime(2000, 1, 1),
    max_value=datetime(2099, 12, 31, 23, 59, 59),
    timezones=st.just(UTC),
)


def totals_of(receipt: Receipt) -> tuple[int, int, int, int, int, int]:
    return (
        receipt.subtotal,
        receipt.promotion_discount,
        receipt.coupon_discount,
        receipt.shipping,
        receipt.tax,
        receipt.total,
    )


# ---------------------------------------------------------------------------
# Properties of a successful checkout
# ---------------------------------------------------------------------------


@given(basket=valid_baskets)
def test_p1_total_equals_subtotal_minus_discounts_plus_shipping_plus_tax(basket: Basket) -> None:
    r = checkout(basket, CATALOG)

    assert r.total == r.subtotal - r.promotion_discount - r.coupon_discount + r.shipping + r.tax


@given(basket=valid_baskets)
def test_p2_nothing_is_negative_and_a_discount_never_exceeds_what_it_discounts(
    basket: Basket,
) -> None:
    r = checkout(basket, CATALOG)

    assert min(totals_of(r)) >= 0
    for line in r.lines:
        assert min(line.unit_price, line.line_total, line.line_discount) >= 0
        assert line.line_discount <= line.line_total
    assert r.promotion_discount <= r.subtotal
    assert r.coupon_discount <= r.subtotal - r.promotion_discount


@given(basket=valid_baskets)
def test_p3_the_receipt_lines_explain_the_totals(basket: Basket) -> None:
    r = checkout(basket, CATALOG)

    assert r.subtotal == sum(line.line_total for line in r.lines)
    assert r.promotion_discount == sum(line.line_discount for line in r.lines)
    for line in r.lines:
        assert line.line_total == line.unit_price * line.quantity


@given(basket=valid_baskets, data=st.data())
def test_p4_shuffling_the_items_gives_the_same_totals(basket: Basket, data: st.DataObject) -> None:
    # `data.draw` lets a test ask for more random values that depend on earlier ones.
    shuffled = data.draw(st.permutations(basket.items))

    assert totals_of(checkout(Basket(shuffled, basket.coupon, basket.shipping), CATALOG)) == (
        totals_of(checkout(basket, CATALOG))
    )


@given(basket=valid_baskets, data=st.data())
def test_p5_splitting_one_line_into_two_lines_of_the_same_sku_gives_the_same_receipt(
    basket: Basket, data: st.DataObject
) -> None:
    splittable = [i for i, item in enumerate(basket.items) if item.quantity >= 2]
    assume(splittable)  # a basket of single items cannot be split: skip it
    index = data.draw(st.sampled_from(splittable))
    line = basket.items[index]
    first = data.draw(st.integers(1, line.quantity - 1))

    split = [
        *basket.items[:index],
        BasketItem(line.sku, first),
        *basket.items[index + 1 :],
        BasketItem(line.sku, line.quantity - first),
    ]

    assert checkout(Basket(split, basket.coupon, basket.shipping), CATALOG) == checkout(
        basket, CATALOG
    )


@given(basket=valid_baskets)
def test_p6_the_shipping_table_is_applied_exactly(basket: Basket) -> None:
    r = checkout(basket, CATALOG)
    discounted = r.subtotal - r.promotion_discount - r.coupon_discount

    assert r.shipping in (0, 499, 999)
    assert (r.shipping == 999) == (basket.shipping == "express")
    if basket.shipping == "standard":
        assert (r.shipping == 0) == (discounted >= 5000)


@given(basket=valid_baskets)
def test_p7_tax_is_8_percent_rounded_half_up_with_integer_division(basket: Basket) -> None:
    r = checkout(basket, CATALOG)
    discounted = r.subtotal - r.promotion_discount - r.coupon_discount

    assert r.tax == (discounted * 8 + 50) // 100


@given(basket=baskets(coupon=st.sampled_from(REAL_COUPONS)).filter(checks_out), data=st.data())
def test_p8_writing_the_coupon_in_another_case_or_with_spaces_gives_the_same_receipt(
    basket: Basket, data: st.DataObject
) -> None:
    assert basket.coupon is not None
    sloppy_code = data.draw(sloppy(basket.coupon))
    assert normalize_coupon_code(sloppy_code) == basket.coupon  # same code, written differently

    assert checkout(Basket(basket.items, sloppy_code, basket.shipping), CATALOG) == checkout(
        basket, CATALOG
    )


@given(basket=valid_baskets)
def test_p9_one_line_per_distinct_sku_in_order_of_first_appearance_each_1_to_99(
    basket: Basket,
) -> None:
    r = checkout(basket, CATALOG)

    # dict.fromkeys keeps the first appearance of each SKU and drops later repeats.
    assert [line.sku for line in r.lines] == list(dict.fromkeys(item.sku for item in basket.items))
    for line in r.lines:
        assert isinstance(line.quantity, int)
        assert 1 <= line.quantity <= 99
        assert line.quantity == sum(item.quantity for item in basket.items if item.sku == line.sku)


# ---------------------------------------------------------------------------
# Properties of any input at all
# ---------------------------------------------------------------------------


@settings(max_examples=500)
@given(basket=any_baskets)
def test_p10_any_basket_gives_a_schema_valid_receipt_or_one_of_the_six_error_codes(
    basket: Basket,
) -> None:
    try:
        receipt = checkout(basket, CATALOG)
    except CheckoutError as error:
        # The only acceptable way to fail: one of the six codes. Anything else (TypeError,
        # ZeroDivisionError, ...) would escape this `except` and fail the test.
        assert error.code in ErrorCode
    else:
        assert schema_errors("receipt.schema.json", receipt.to_dict()) == []


@given(amount=st.integers(0, 10_000_000), percent=st.integers(0, 100))
def test_p11_round_half_up_is_within_half_a_cent_and_exact_when_the_value_is_whole(
    amount: int, percent: int
) -> None:
    rounded = round_half_up(amount, percent)
    exact_times_100 = amount * percent  # the exact value is this / 100

    assert isinstance(rounded, int)
    assert abs(rounded * 100 - exact_times_100) <= 50
    if exact_times_100 % 100 == 0:
        assert rounded * 100 == exact_times_100
    if exact_times_100 % 100 == 50:
        assert rounded * 100 == exact_times_100 + 50  # halves go UP


# ---------------------------------------------------------------------------
# Properties of saving orders
# ---------------------------------------------------------------------------


def fixed_clock() -> datetime:
    return datetime(2026, 1, 15, 10, 30, 0, tzinfo=UTC)


@given(baskets_to_save=st.lists(valid_baskets, max_size=12))
def test_p12_after_saving_n_orders_the_ids_are_exactly_order_1_to_order_n(
    baskets_to_save: list[Basket],
) -> None:
    repository = InMemoryOrderRepository(fixed_clock)

    ids = [place_order(basket, CATALOG, repository).order_id for basket in baskets_to_save]

    assert ids == [f"ORDER-{n}" for n in range(1, len(baskets_to_save) + 1)]


@given(basket=valid_baskets, when=moments)
def test_p13_find_of_save_gives_back_the_same_receipt_and_the_clock_s_time(
    basket: Basket, when: datetime
) -> None:
    repository = InMemoryOrderRepository(lambda: when)
    receipt = checkout(basket, CATALOG)

    order = repository.find(repository.save(receipt))

    assert order is not None
    assert order.receipt == receipt
    assert order.created_at == when.replace(microsecond=0).isoformat().replace("+00:00", "Z")


@given(basket=any_baskets.filter(lambda b: not checks_out(b)), random_id=st.text(max_size=12))
def test_p14_a_failed_checkout_saves_nothing_and_find_of_an_unknown_id_is_not_found(
    basket: Basket, random_id: str
) -> None:
    repository = InMemoryOrderRepository(fixed_clock)

    with pytest.raises(CheckoutError):
        place_order(basket, CATALOG, repository)

    assert repository.count() == 0
    assert repository.find(random_id) is None
    assert repository.find("ORDER-1") is None
