//! Small, hand-made test data for the unit tests in this crate.
//!
//! Unit tests prefer tiny made-up data over the real catalog: when a test fails you can
//! work out the right answer in your head. The shared spec cases (in `tests/`) use the
//! real catalog from `spec/catalog.json`.
//!
//! This module is only compiled for `cargo test` (see `#[cfg(test)]` in `lib.rs`).
//! The integration tests in `tests/` have their own copy in `tests/common/fixtures.rs`,
//! because they can only see the crate's public API.

use crate::{Basket, BasketItem, Catalog, Product, Promotion, Receipt, ReceiptLine};

/// $0.50, no promotion.
pub fn apple() -> Product {
    product("APPLE", "Apple", 50, None)
}

/// $1.20, no promotion.
pub fn pen() -> Product {
    product("PEN", "Pen", 120, None)
}

/// $8.99, buy 2 get 1 free.
pub fn coffee() -> Product {
    product(
        "COFFEE",
        "Bag of Coffee Beans",
        899,
        Some(Promotion::BuyXGetYFree { buy: 2, free: 1 }),
    )
}

/// $3.49, 10% off from 10.
pub fn notebook() -> Product {
    product(
        "NOTEBOOK",
        "Notebook",
        349,
        Some(Promotion::BulkDiscount {
            min_quantity: 10,
            percent: 10,
        }),
    )
}

/// A $100.00 item: handy for reaching free shipping and coupon minimums quickly.
pub fn tv() -> Product {
    product("TV", "Television", 10000, None)
}

/// Five products, two of them with a promotion.
pub fn tiny_catalog() -> Catalog {
    Catalog::new(vec![apple(), pen(), coffee(), notebook(), tv()])
}

/// A valid receipt for one apple, used wherever "some receipt" is needed.
pub fn one_apple_receipt() -> Receipt {
    Receipt {
        lines: vec![ReceiptLine {
            sku: "APPLE".to_string(),
            name: "Apple".to_string(),
            quantity: 1,
            unit_price: 50,
            line_total: 50,
            line_discount: 0,
        }],
        subtotal: 50,
        promotion_discount: 0,
        coupon_discount: 0,
        shipping: 499,
        tax: 4,
        total: 553,
    }
}

/// One basket line.
pub fn item(sku: &str, quantity: i64) -> BasketItem {
    BasketItem {
        sku: sku.to_string(),
        quantity,
    }
}

/// A basket with these lines, no coupon and standard shipping.
pub fn basket(items: &[(&str, i64)]) -> Basket {
    basket_with(items, None, "standard")
}

/// A basket with these lines, this coupon (or `None`) and this shipping method.
pub fn basket_with(items: &[(&str, i64)], coupon: Option<&str>, shipping: &str) -> Basket {
    Basket {
        items: items
            .iter()
            .map(|(sku, quantity)| item(sku, *quantity))
            .collect(),
        coupon: coupon.map(str::to_string),
        shipping: shipping.to_string(),
    }
}

fn product(sku: &str, name: &str, unit_price: i64, promotion: Option<Promotion>) -> Product {
    Product {
        sku: sku.to_string(),
        name: name.to_string(),
        unit_price,
        promotion,
    }
}
