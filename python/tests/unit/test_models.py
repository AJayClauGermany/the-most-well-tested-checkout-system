"""Unit tests for the data classes in models.py: reading JSON in and writing it out."""

from support.fixtures import ONE_APPLE_RECEIPT

from checkout import Basket, BasketItem, Receipt, ReceiptLine


def test_a_basket_is_read_from_its_json_form() -> None:
    basket = Basket.from_dict(
        {"items": [{"sku": "COFFEE", "quantity": 3}], "coupon": "SAVE10", "shipping": "express"}
    )

    assert basket == Basket(items=[BasketItem("COFFEE", 3)], coupon="SAVE10", shipping="express")


def test_a_basket_without_coupon_and_shipping_means_no_coupon_and_standard() -> None:
    basket = Basket.from_dict({"items": []})

    assert basket == Basket(items=[], coupon=None, shipping="standard")


def test_a_receipt_s_json_form_uses_the_schema_s_snake_case_names() -> None:
    assert ONE_APPLE_RECEIPT.to_dict() == {
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


def test_a_receipt_read_back_from_json_equals_the_original() -> None:
    receipt = Receipt(
        lines=[
            ReceiptLine("COFFEE", "Coffee", 3, 899, 2697, 899),
            ReceiptLine("APPLE", "Apple", 2, 50, 100, 0),
        ],
        subtotal=2797,
        promotion_discount=899,
        coupon_discount=190,
        shipping=499,
        tax=137,
        total=2344,
    )

    assert Receipt.from_dict(receipt.to_dict()) == receipt


def test_to_dict_gives_a_copy_not_a_view_of_the_receipt() -> None:
    as_dict = ONE_APPLE_RECEIPT.to_dict()

    as_dict["lines"].clear()

    assert len(ONE_APPLE_RECEIPT.lines) == 1
