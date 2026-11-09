//! Small, hand-made test data for the integration tests, plus a clock the tests can move.
//!
//! Unit tests prefer tiny made-up data over the real catalog: when a test fails you can
//! work out the right answer in your head. The shared spec cases use the real catalog.
//! (The unit tests inside `src/` have their own copy in `src/test_fixtures.rs`, because
//! integration tests can only see the crate's public API.)

use checkout::{Basket, BasketItem, Catalog, Clock, Product, Promotion, Receipt, ReceiptLine};
use std::sync::{Arc, Mutex};
use std::time::SystemTime;

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

/// The basket behind `spec/examples/receipt.json` (rules §5): 3 coffees, 2 apples, SAVE10.
pub fn example_basket() -> Basket {
    basket_with(&[("COFFEE", 3), ("APPLE", 2)], Some("SAVE10"), "standard")
}

/// The biggest shared case: every product, both promotions, WELCOME5 and express shipping.
pub fn full_everything_basket() -> Basket {
    basket_with(
        &[
            ("HEADPHONES", 1),
            ("COFFEE", 6),
            ("NOTEBOOK", 12),
            ("PEN", 5),
            ("BREAD", 1),
            ("APPLE", 3),
        ],
        Some("WELCOME5"),
        "express",
    )
}

/// A clock the test can move while a repository is already holding it: the repository
/// and the test share the same `Arc<Mutex<SystemTime>>`.
///
/// This is a *stub* that can change its canned answer – handy for "the clock is fixed at"
/// steps and for "stamps each order with the time at the moment of saving" tests.
#[derive(Debug, Clone)]
pub struct SharedClock(Arc<Mutex<SystemTime>>);

impl SharedClock {
    pub fn starting_at(time: SystemTime) -> Self {
        SharedClock(Arc::new(Mutex::new(time)))
    }

    /// Moves the clock. Every repository holding a clone of this clock sees the new time.
    pub fn set(&self, time: SystemTime) {
        *self.0.lock().unwrap() = time;
    }
}

impl Clock for SharedClock {
    fn now(&self) -> SystemTime {
        *self.0.lock().unwrap()
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
