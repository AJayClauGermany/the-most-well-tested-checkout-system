//! Basket validation (rules §4.1, rows 1 to 4). Rows 5 to 7 (shipping and coupon) are
//! checked by `shipping.rs` and `coupons.rs`, because the coupon check needs the price
//! after promotions.

use crate::{Basket, Catalog, CheckoutError, ErrorCode, merge_items};

/// The smallest quantity a basket line may have (rules §4.1 row 3).
pub const MIN_QUANTITY: i64 = 1;
/// The largest quantity a basket line – and a whole SKU – may have (rules §4.1 rows 3 and 4).
pub const MAX_QUANTITY: i64 = 99;

/// True when `quantity` is from 1 to 99 (rules §4.1 row 3).
///
/// Quantities are `i64`, so "not a whole number" (`1.5`) cannot even be written down –
/// the type system rules it out before this function runs.
///
/// ```
/// use checkout::is_valid_quantity;
///
/// assert!(is_valid_quantity(1));
/// assert!(is_valid_quantity(99));
/// assert!(!is_valid_quantity(0));
/// assert!(!is_valid_quantity(100));
/// ```
pub fn is_valid_quantity(quantity: i64) -> bool {
    (MIN_QUANTITY..=MAX_QUANTITY).contains(&quantity)
}

/// Checks rows 1-4 of rules §4.1, in that order.
///
/// Returns `Ok(())` when the items are fine, or the first `CheckoutError` that applies.
pub fn validate_items(basket: &Basket, catalog: &Catalog) -> Result<(), CheckoutError> {
    // Row 1: no items at all.
    if basket.items.is_empty() {
        return Err(CheckoutError::with_message(
            ErrorCode::EmptyBasket,
            "The basket has no items",
        ));
    }

    // Row 2: every SKU must exist. All SKUs are checked before any quantity, so an unknown
    // product on line 3 still wins over a bad quantity on line 1.
    for item in &basket.items {
        catalog.require(&item.sku)?;
    }

    // Row 3: every quantity must be from 1 to 99.
    for item in &basket.items {
        if !is_valid_quantity(item.quantity) {
            return Err(CheckoutError::with_message(
                ErrorCode::InvalidQuantity,
                format!("Invalid quantity for {}: {}", item.sku, item.quantity),
            ));
        }
    }

    // Row 4: the lines of one SKU, added together, may not go over 99 either.
    for line in merge_items(&basket.items) {
        if line.quantity > MAX_QUANTITY {
            return Err(CheckoutError::with_message(
                ErrorCode::InvalidQuantity,
                format!("Too many {} in total: {}", line.sku, line.quantity),
            ));
        }
    }

    Ok(())
}

// ---------------------------------------------------------------------------
// Unit tests for rows 1 to 4 and the order errors are reported in. The check order
// matters: a basket can be wrong in several ways at once.
// ---------------------------------------------------------------------------
#[cfg(test)]
mod tests {
    use super::*;
    use crate::ErrorCode;
    use crate::test_fixtures::{basket, basket_with, tiny_catalog};
    use rstest::rstest;

    /// Validates `basket` against the tiny catalog and returns the error (the test fails if there is none).
    fn error_for(basket: &Basket) -> CheckoutError {
        validate_items(basket, &tiny_catalog()).expect_err("expected the basket to be invalid")
    }

    #[test]
    fn has_the_limits_from_the_rules() {
        assert_eq!(MIN_QUANTITY, 1);
        assert_eq!(MAX_QUANTITY, 99);
    }

    // Boundaries are tested on BOTH sides: 0 / 1 and 99 / 100.
    #[rstest]
    #[case(1)]
    #[case(2)]
    #[case(50)]
    #[case(98)]
    #[case(99)]
    fn accepts_quantities_from_1_to_99(#[case] quantity: i64) {
        assert!(is_valid_quantity(quantity));
    }

    #[rstest]
    #[case(0)]
    #[case(-1)]
    #[case(100)]
    #[case(101)]
    #[case(1000)]
    #[case(i64::MIN)]
    #[case(i64::MAX)]
    fn rejects_quantities_outside_1_to_99(#[case] quantity: i64) {
        assert!(!is_valid_quantity(quantity));
    }

    #[test]
    fn accepts_a_normal_basket() {
        let basket = basket(&[("APPLE", 1), ("COFFEE", 99)]);

        assert_eq!(validate_items(&basket, &tiny_catalog()), Ok(()));
    }

    #[test]
    fn row_1_an_empty_basket() {
        assert_eq!(error_for(&basket(&[])).code, ErrorCode::EmptyBasket);
    }

    #[test]
    fn row_2_a_sku_that_is_not_in_the_catalog() {
        assert_eq!(
            error_for(&basket(&[("BANANA", 1)])).code,
            ErrorCode::UnknownProduct
        );
    }

    #[test]
    fn row_2_skus_are_case_sensitive() {
        assert_eq!(
            error_for(&basket(&[("apple", 1)])).code,
            ErrorCode::UnknownProduct
        );
    }

    #[test]
    fn row_2_an_unknown_product_on_any_line_not_only_the_first() {
        let basket = basket(&[("APPLE", 1), ("PEN", 1), ("BANANA", 1)]);

        assert_eq!(error_for(&basket).code, ErrorCode::UnknownProduct);
    }

    #[rstest]
    #[case(0)]
    #[case(-1)]
    #[case(100)]
    fn row_3_a_bad_quantity(#[case] quantity: i64) {
        assert_eq!(
            error_for(&basket(&[("APPLE", quantity)])).code,
            ErrorCode::InvalidQuantity
        );
    }

    #[test]
    fn row_3_a_bad_quantity_on_a_later_line() {
        let basket = basket(&[("APPLE", 1), ("PEN", 0)]);

        assert_eq!(error_for(&basket).code, ErrorCode::InvalidQuantity);
    }

    #[test]
    fn row_4_lines_of_one_sku_that_add_up_to_more_than_99() {
        let basket = basket(&[("APPLE", 50), ("APPLE", 50)]);

        assert_eq!(error_for(&basket).code, ErrorCode::InvalidQuantity);
    }

    #[test]
    fn row_4_lines_of_one_sku_that_add_up_to_exactly_99_are_fine() {
        let basket = basket(&[("APPLE", 50), ("PEN", 99), ("APPLE", 49)]);

        assert_eq!(validate_items(&basket, &tiny_catalog()), Ok(()));
    }

    #[test]
    fn row_4_99_of_one_sku_and_99_of_another_is_fine_the_limit_is_per_sku() {
        let basket = basket(&[("APPLE", 99), ("PEN", 99)]);

        assert_eq!(validate_items(&basket, &tiny_catalog()), Ok(()));
    }

    // The code is what programs look at; the message is for the human reading a log.
    #[rstest]
    #[case::empty(basket(&[]), "The basket has no items")]
    #[case::unknown(basket(&[("BANANA", 1)]), "Unknown product: BANANA")]
    #[case::quantity(basket(&[("APPLE", 0)]), "Invalid quantity for APPLE: 0")]
    #[case::merged(basket(&[("APPLE", 60), ("APPLE", 40)]), "Too many APPLE in total: 100")]
    fn error_messages_explain_the_problem(#[case] basket: Basket, #[case] message: &str) {
        assert_eq!(error_for(&basket).message, message);
    }

    #[test]
    fn check_order_an_unknown_product_beats_a_bad_quantity_even_on_a_later_line() {
        let basket = basket(&[("APPLE", 0), ("BANANA", 1)]);

        assert_eq!(error_for(&basket).code, ErrorCode::UnknownProduct);
    }

    #[test]
    fn check_order_a_bad_single_quantity_beats_a_bad_merged_quantity() {
        // Both rows 3 and 4 apply; row 3 is reported, and its message names the line.
        let basket = basket(&[("APPLE", 100), ("APPLE", 1)]);

        let error = error_for(&basket);

        assert_eq!(error.code, ErrorCode::InvalidQuantity);
        assert_eq!(error.message, "Invalid quantity for APPLE: 100");
    }

    #[test]
    fn check_order_ignores_the_coupon_and_shipping_method_those_are_checked_later() {
        let basket = basket_with(&[("APPLE", 1)], Some("NOPE"), "overnight");

        assert_eq!(validate_items(&basket, &tiny_catalog()), Ok(()));
    }
}
