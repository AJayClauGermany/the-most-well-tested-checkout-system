"""Unit tests for Catalog.

Notice that these tests build their own tiny catalog instead of reading
spec/catalog.json. That keeps them independent of any file: if the shared
catalog changes one day, these tests still describe how Catalog behaves.
"""

from checkout import BulkDiscount, Catalog, Product

APPLE = Product(sku="APPLE", name="Apple", unit_price=50)
NOTEBOOK = Product(
    sku="NOTEBOOK",
    name="Notebook",
    unit_price=349,
    promotion=BulkDiscount(min_quantity=10, percent=10),
)


def test_find_returns_the_product_with_that_sku():
    catalog = Catalog(products=[APPLE, NOTEBOOK])

    product = catalog.find("NOTEBOOK")

    assert product == NOTEBOOK


def test_find_returns_none_for_an_unknown_sku():
    catalog = Catalog(products=[APPLE, NOTEBOOK])

    product = catalog.find("BANANA")

    # For None, True and False, Python style is `is`, not `==`.
    assert product is None


def test_find_is_case_sensitive():
    # The spec says "apple" is NOT the same product as "APPLE".
    catalog = Catalog(products=[APPLE])

    assert catalog.find("apple") is None
