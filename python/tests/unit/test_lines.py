"""Unit tests for merging basket items into receipt lines and pricing them (rules §4.2)."""

from support.fixtures import APPLE, COFFEE, NOTEBOOK

from checkout import BasketItem, ReceiptLine, merge_items, price_line

# --- merge_items ---------------------------------------------------------------


def test_keeps_items_with_different_skus_apart_in_basket_order() -> None:
    merged = merge_items([BasketItem("PEN", 2), BasketItem("APPLE", 1)])

    assert merged == [BasketItem("PEN", 2), BasketItem("APPLE", 1)]


def test_adds_up_the_quantities_of_the_same_sku_keeping_the_place_of_its_first_appearance() -> None:
    # The example from rules §4.2: [APPLE x1, PEN x2, APPLE x3] gives [APPLE x4, PEN x2].
    merged = merge_items([BasketItem("APPLE", 1), BasketItem("PEN", 2), BasketItem("APPLE", 3)])

    assert merged == [BasketItem("APPLE", 4), BasketItem("PEN", 2)]


def test_merges_three_lines_of_one_sku_into_one() -> None:
    merged = merge_items(
        [BasketItem("COFFEE", 1), BasketItem("COFFEE", 1), BasketItem("COFFEE", 1)]
    )

    assert merged == [BasketItem("COFFEE", 3)]


def test_gives_an_empty_list_for_an_empty_list() -> None:
    assert merge_items([]) == []


def test_does_not_change_the_basket_s_own_items() -> None:
    items = [BasketItem("APPLE", 1), BasketItem("APPLE", 1)]

    merge_items(items)

    assert items == [BasketItem("APPLE", 1), BasketItem("APPLE", 1)]


def test_accepts_a_tuple_as_well_as_a_list() -> None:
    assert merge_items((BasketItem("APPLE", 1), BasketItem("APPLE", 2))) == [BasketItem("APPLE", 3)]


# --- price_line ------------------------------------------------------------------


def test_copies_sku_name_and_unit_price_from_the_product_and_multiplies_for_the_total() -> None:
    assert price_line(APPLE, 3) == ReceiptLine(
        sku="APPLE", name="Apple", quantity=3, unit_price=50, line_total=150, line_discount=0
    )


def test_includes_the_promotion_discount() -> None:
    coffee = price_line(COFFEE, 3)
    notebooks = price_line(NOTEBOOK, 10)

    assert (coffee.line_total, coffee.line_discount) == (2697, 899)
    assert (notebooks.line_total, notebooks.line_discount) == (3490, 349)


def test_gives_no_discount_below_the_promotion_s_threshold() -> None:
    assert price_line(COFFEE, 2).line_discount == 0
    assert price_line(NOTEBOOK, 9).line_discount == 0
