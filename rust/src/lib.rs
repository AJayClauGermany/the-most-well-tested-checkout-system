//! A small checkout system, used to show off many kinds of software testing.
//!
//! The rules live in `spec/checkout-rules.md` at the root of the repository. The code
//! follows the hint in rules §4.7 and reads like a pipeline of small steps, each in its
//! own module with its own unit tests:
//!
//! | Module                   | Rules  | What it does                                            |
//! |--------------------------|--------|---------------------------------------------------------|
//! | `validation`             | §4.1   | rows 1-4: empty basket, unknown product, bad quantities |
//! | `lines`                  | §4.2   | merges basket items into receipt lines and prices them  |
//! | `promotions`             | §4.2   | buy X get Y free, bulk discount                         |
//! | `coupons`                | §4.4   | `SAVE10`, `WELCOME5` (rows 6-7 of §4.1)                 |
//! | `shipping`               | §4.5   | standard / express (row 5 of §4.1)                      |
//! | `tax`                    | §4.6   | 8%                                                      |
//! | `money`                  | §7     | `round_half_up`, with whole numbers only                |
//! | `orders`                 | §6     | saving orders: `Clock`, `OrderRepository`, `place_order` |
//! | `json_file_repository`   | §6     | one `ORDER-<n>.json` file per order                     |
//! | `sqlite_repository`      | §6     | one table in an SQLite file                             |
//! | `spec_folder`            | –      | finds the repository's `spec/` folder by walking up     |
//! | `command`                | §8     | the command line program's logic, unit testable         |
//!
//! Where do the tests live? (Rust has its own conventions!)
//!
//! * **Unit tests** sit at the bottom of each `src/*.rs` file, inside a
//!   `#[cfg(test)] mod tests { ... }` block. They can even test private functions.
//! * **Integration tests** live in the `tests/` folder. Each file there is compiled as a
//!   separate program that can only use the *public* API, exactly like a real user would.
//! * **Doc tests** are the code examples inside `///` comments, like the one on [`checkout`].

// Each `mod` line pulls in the file with the same name from `src/`.
mod basket;
mod catalog;
mod command;
mod coupons;
mod error;
mod json_file_repository;
mod lines;
mod money;
mod orders;
mod promotions;
mod receipt;
mod shipping;
mod spec_folder;
mod sqlite_repository;
mod tax;
mod validation;

// Small, hand-made test data shared by the unit tests in this crate (compiled only for tests).
#[cfg(test)]
pub(crate) mod test_fixtures;

// Re-export the types so users can write `checkout::Basket` instead of `checkout::basket::Basket`.
pub use basket::{Basket, BasketItem};
pub use catalog::{Catalog, Product, Promotion};
pub use command::{
    CommandError, CommandResult, load_standard_catalog, run_checkout_command, standard_catalog_path,
};
pub use coupons::{coupon_discount, normalize_coupon_code};
pub use error::{CheckoutError, ErrorCode};
pub use json_file_repository::JsonFileOrderRepository;
pub use lines::{merge_items, price_line};
pub use money::round_half_up;
pub use orders::{
    Clock, FixedClock, InMemoryOrderRepository, Order, OrderRepository, PlaceOrderError,
    PlacedOrder, RepositoryError, SystemClock, format_timestamp, order_id, order_number,
    parse_timestamp, place_order,
};
pub use promotions::promotion_discount;
pub use receipt::{Receipt, ReceiptLine};
pub use shipping::{
    EXPRESS_COST, FREE_SHIPPING_FROM, STANDARD_COST, ShippingMethod, parse_shipping_method,
    shipping_cost,
};
pub use spec_folder::{find_spec_folder, locate_spec_folder};
pub use sqlite_repository::SqliteOrderRepository;
pub use tax::{TAX_PERCENT, tax_on};
pub use validation::{MAX_QUANTITY, MIN_QUANTITY, is_valid_quantity, validate_items};

/// An amount of money in whole cents: `199` means $1.99.
///
/// We never use floating point numbers for money (see rules §1).
/// `i64` (a signed number) keeps subtraction simple and can never overflow for our baskets.
pub type Cents = i64;

/// Calculates the receipt for a basket, using the products in `catalog`.
///
/// Returns `Ok(receipt)` on success, or `Err(error)` carrying one of the six
/// [`ErrorCode`]s when the basket is invalid (rules §4.1).
///
/// # Example
///
/// This example is a *doc test*: `cargo test` compiles and runs it, so it can never
/// go out of date.
///
/// ```
/// use checkout::{Basket, BasketItem, Catalog, Product, checkout};
///
/// let catalog = Catalog::new(vec![Product {
///     sku: "APPLE".to_string(),
///     name: "Apple".to_string(),
///     unit_price: 50,
///     promotion: None,
/// }]);
/// let basket = Basket {
///     items: vec![BasketItem { sku: "APPLE".to_string(), quantity: 2 }],
///     coupon: None,
///     shipping: "standard".to_string(),
/// };
///
/// let receipt = checkout(&basket, &catalog).unwrap();
/// assert_eq!(receipt.subtotal, 100);
/// assert_eq!(receipt.shipping, 499);
/// assert_eq!(receipt.tax, 8);
/// assert_eq!(receipt.total, 607);
/// ```
pub fn checkout(basket: &Basket, catalog: &Catalog) -> Result<Receipt, CheckoutError> {
    // The pipeline from the hint in rules §4.7. Every `?` stops at the first error, which is
    // what makes the error order of rules §4.1 come out right.
    validate_items(basket, catalog)?; // rows 1-4
    let method = parse_shipping_method(&basket.shipping)?; // row 5

    let lines = merge_items(&basket.items)
        .iter()
        .map(|item| Ok(price_line(catalog.require(&item.sku)?, item.quantity)))
        .collect::<Result<Vec<ReceiptLine>, CheckoutError>>()?;
    let subtotal: Cents = lines.iter().map(|line| line.line_total).sum();
    let promotion_discount: Cents = lines.iter().map(|line| line.line_discount).sum();
    let after_promotions = subtotal - promotion_discount;

    let coupon_discount = coupon_discount(basket.coupon.as_deref(), after_promotions)?; // rows 6-7
    let discounted_subtotal = after_promotions - coupon_discount;

    let shipping = shipping_cost(method, discounted_subtotal);
    let tax = tax_on(discounted_subtotal);

    Ok(Receipt {
        lines,
        subtotal,
        promotion_discount,
        coupon_discount,
        shipping,
        tax,
        total: discounted_subtotal + shipping + tax,
    })
}

// ---------------------------------------------------------------------------
// Unit tests of checkout() as a whole, with the tiny made-up catalog from test_fixtures.rs.
//
// The 60 shared spec cases (tests/spec_cases.rs) cover the standard catalog; these tests
// check that the pieces fit together, and the things the six totals alone cannot show,
// such as the receipt lines and the order in which errors are reported.
// ---------------------------------------------------------------------------
#[cfg(test)]
mod tests {
    use super::*;
    use crate::test_fixtures::{basket, basket_with, item, tiny_catalog};
    use pretty_assertions::assert_eq;
    use rstest::rstest;

    #[test]
    fn prices_one_apple_50_cents_plus_499_shipping_plus_4_tax() {
        // Arrange
        let catalog = tiny_catalog();

        // Act
        let receipt = checkout(&basket(&[("APPLE", 1)]), &catalog).unwrap();

        // Assert: the whole receipt, lines included.
        let expected = Receipt {
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
        };
        assert_eq!(receipt, expected);
    }

    #[test]
    fn writes_one_line_per_distinct_sku_in_order_of_first_appearance_with_merged_quantities() {
        let catalog = tiny_catalog();
        let basket = basket(&[("PEN", 1), ("APPLE", 2), ("PEN", 3)]);

        let receipt = checkout(&basket, &catalog).unwrap();

        let lines: Vec<(&str, i64)> = receipt
            .lines
            .iter()
            .map(|line| (line.sku.as_str(), line.quantity))
            .collect();
        assert_eq!(lines, vec![("PEN", 4), ("APPLE", 2)]);
    }

    #[test]
    fn works_the_totals_out_from_the_lines() {
        // 3 coffees ($26.97, one free: -$8.99) and 10 notebooks ($34.90, 10% off: -$3.49)
        let catalog = tiny_catalog();
        let basket = basket(&[("COFFEE", 3), ("NOTEBOOK", 10)]);

        let receipt = checkout(&basket, &catalog).unwrap();

        assert_eq!(receipt.subtotal, 2697 + 3490);
        assert_eq!(receipt.promotion_discount, 899 + 349);
        assert_eq!(receipt.coupon_discount, 0);
        assert_eq!(receipt.shipping, 499); // $49.39 after promotions: not free
        assert_eq!(receipt.tax, 395); // 8% of 4939 = 395.12
        assert_eq!(receipt.total, 4939 + 499 + 395);
    }

    #[test]
    fn takes_the_coupon_off_after_promotions_and_before_shipping_and_tax() {
        // TV $100.00 with SAVE10: $90.00 → free shipping, tax 8% of $90.00 = $7.20
        let catalog = tiny_catalog();
        let basket = basket_with(&[("TV", 1)], Some("SAVE10"), "standard");

        let receipt = checkout(&basket, &catalog).unwrap();

        assert_eq!(receipt.subtotal, 10000);
        assert_eq!(receipt.promotion_discount, 0);
        assert_eq!(receipt.coupon_discount, 1000);
        assert_eq!(receipt.shipping, 0);
        assert_eq!(receipt.tax, 720);
        assert_eq!(receipt.total, 9720);
    }

    #[test]
    fn charges_express_shipping_on_top_untaxed() {
        let catalog = tiny_catalog();
        let basket = basket_with(&[("TV", 1)], None, "express");

        let receipt = checkout(&basket, &catalog).unwrap();

        assert_eq!(receipt.shipping, 999);
        assert_eq!(receipt.tax, 800);
        assert_eq!(receipt.total, 10000 + 999 + 800);
    }

    #[test]
    fn a_coupon_can_push_the_order_below_the_free_shipping_threshold() {
        // $50.00 exactly is free shipping; SAVE10 makes it $45.00, so shipping is charged again.
        let catalog = tiny_catalog();
        let basket = basket_with(&[("APPLE", 99), ("PEN", 1)], Some("SAVE10"), "standard");

        let receipt = checkout(&basket, &catalog).unwrap();

        assert_eq!(receipt.subtotal, 5070);
        assert_eq!(receipt.coupon_discount, 507);
        assert_eq!(receipt.shipping, 499);
    }

    // A parameterised test (rstest's `#[case]`): the same test body runs once per row.
    // Each row is a basket that breaks several rules at once; only the FIRST row of the
    // table in rules §4.1 may be reported.
    #[rstest]
    #[case::empty_basket_before_anything_else(
        basket_with(&[], Some("NOPE"), "overnight"),
        ErrorCode::EmptyBasket
    )]
    #[case::unknown_product_before_bad_quantity_shipping_and_coupon(
        basket_with(&[("APPLE", 0), ("BANANA", 1)], Some("NOPE"), "x"),
        ErrorCode::UnknownProduct
    )]
    #[case::bad_quantity_before_shipping_and_coupon(
        basket_with(&[("APPLE", 0)], Some("NOPE"), "x"),
        ErrorCode::InvalidQuantity
    )]
    #[case::merged_quantity_over_99_before_shipping_and_coupon(
        basket_with(&[("APPLE", 60), ("APPLE", 40)], Some("NOPE"), "x"),
        ErrorCode::InvalidQuantity
    )]
    #[case::bad_shipping_before_bad_coupon(
        basket_with(&[("APPLE", 1)], Some("NOPE"), "x"),
        ErrorCode::UnknownShippingMethod
    )]
    #[case::bad_shipping_before_an_unmet_coupon_minimum(
        basket_with(&[("APPLE", 1)], Some("WELCOME5"), "x"),
        ErrorCode::UnknownShippingMethod
    )]
    #[case::unknown_coupon(
        basket_with(&[("APPLE", 1)], Some("NOPE"), "standard"),
        ErrorCode::UnknownCoupon
    )]
    #[case::coupon_minimum_judged_on_the_price_after_promotions(
        // 3 coffees are $26.97 but $17.98 after the free one: below WELCOME5's $25.00.
        basket_with(&[("COFFEE", 3)], Some("WELCOME5"), "standard"),
        ErrorCode::CouponMinimumNotMet
    )]
    fn reports_the_first_failing_row_of_rules_4_1(
        #[case] basket: Basket,
        #[case] expected: ErrorCode,
    ) {
        let error = checkout(&basket, &tiny_catalog()).unwrap_err();

        assert_eq!(error.code, expected);
    }

    #[test]
    fn does_not_need_to_copy_or_change_the_basket() {
        // `checkout` borrows the basket (`&Basket`), so the compiler already guarantees it
        // cannot change it. This test only shows the basket is still usable afterwards.
        let catalog = tiny_catalog();
        let basket = Basket {
            items: vec![item("APPLE", 1), item("APPLE", 1)],
            coupon: Some(" save10 ".to_string()),
            shipping: "standard".to_string(),
        };

        let first = checkout(&basket, &catalog).unwrap();
        let second = checkout(&basket, &catalog).unwrap();

        assert_eq!(first, second);
        assert_eq!(basket.items.len(), 2);
    }
}
