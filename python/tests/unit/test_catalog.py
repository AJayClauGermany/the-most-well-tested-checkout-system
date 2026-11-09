"""Unit tests for Catalog and require_product.

Notice that these tests build their own tiny catalog instead of reading
spec/catalog.json. That keeps them independent of any file: if the shared
catalog changes one day, these tests still describe how Catalog behaves.
"""

import pytest

from checkout import (
    BulkDiscount,
    BuyXGetYFree,
    Catalog,
    CheckoutError,
    ErrorCode,
    Product,
    require_product,
)

APPLE = Product(sku="APPLE", name="Apple", unit_price=50)
NOTEBOOK = Product(
    sku="NOTEBOOK",
    name="Notebook",
    unit_price=349,
    promotion=BulkDiscount(min_quantity=10, percent=10),
)


def test_find_returns_the_product_with_that_sku() -> None:
    catalog = Catalog(products=[APPLE, NOTEBOOK])

    product = catalog.find("NOTEBOOK")

    assert product == NOTEBOOK


def test_find_returns_none_for_an_unknown_sku() -> None:
    catalog = Catalog(products=[APPLE, NOTEBOOK])

    product = catalog.find("BANANA")

    # For None, True and False, Python style is `is`, not `==`.
    assert product is None


def test_find_is_case_sensitive() -> None:
    # The spec says "apple" is NOT the same product as "APPLE".
    catalog = Catalog(products=[APPLE])

    assert catalog.find("apple") is None


def test_require_product_returns_the_product_with_that_sku() -> None:
    catalog = Catalog(products=[APPLE, NOTEBOOK])

    assert require_product(catalog, "APPLE") == APPLE


def test_require_product_raises_unknown_product_with_the_sku_in_the_message() -> None:
    catalog = Catalog(products=[APPLE])

    with pytest.raises(CheckoutError) as caught:
        require_product(catalog, "BANANA")

    assert caught.value.code == ErrorCode.UNKNOWN_PRODUCT
    assert str(caught.value) == "Unknown product: BANANA"


def test_from_dict_reads_the_json_form_of_a_catalog() -> None:
    catalog = Catalog.from_dict(
        {
            "products": [
                {"sku": "APPLE", "name": "Apple", "unit_price": 50, "promotion": None},
                {
                    "sku": "COFFEE",
                    "name": "Coffee",
                    "unit_price": 899,
                    "promotion": {"type": "buy_x_get_y_free", "buy": 2, "free": 1},
                },
                {
                    "sku": "NOTEBOOK",
                    "name": "Notebook",
                    "unit_price": 349,
                    "promotion": {"type": "bulk_discount", "min_quantity": 10, "percent": 10},
                },
            ]
        }
    )

    assert catalog == Catalog(
        products=[
            APPLE,
            Product("COFFEE", "Coffee", 899, BuyXGetYFree(buy=2, free=1)),
            NOTEBOOK,
        ]
    )


def test_from_dict_rejects_a_promotion_type_it_does_not_know() -> None:
    with pytest.raises(ValueError, match="Unknown promotion type: 'happy_hour'"):
        Catalog.from_dict(
            {
                "products": [
                    {"sku": "X", "name": "X", "unit_price": 1, "promotion": {"type": "happy_hour"}}
                ]
            }
        )
