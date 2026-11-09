//! Shipping (rules §4.5).

use crate::{Cents, CheckoutError};

/// The two shipping methods the shop offers.
///
/// The basket carries the method as free text (`"overnight"` must be representable so it
/// can be rejected); [`parse_shipping_method`] turns that text into this enum. After that,
/// the type system knows the value is one of the two – "parse, don't validate".
#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum ShippingMethod {
    Standard,
    Express,
}

/// What standard shipping costs below the free-shipping threshold.
pub const STANDARD_COST: Cents = 499;
/// What express shipping always costs.
pub const EXPRESS_COST: Cents = 999;
/// Standard shipping is free from this discounted subtotal upwards.
pub const FREE_SHIPPING_FROM: Cents = 5000;

/// Turns the basket's `shipping` text into a known method.
///
/// Fails with `UNKNOWN_SHIPPING_METHOD` for anything but `"standard"` and `"express"`
/// (rules §4.1 row 5). Names are matched exactly: `"Express"` is unknown.
///
/// ```
/// use checkout::{ShippingMethod, parse_shipping_method};
///
/// assert_eq!(parse_shipping_method("express"), Ok(ShippingMethod::Express));
/// assert!(parse_shipping_method("overnight").is_err());
/// ```
pub fn parse_shipping_method(shipping: &str) -> Result<ShippingMethod, CheckoutError> {
    let _ = shipping;
    todo!("parse_shipping_method() is not implemented yet")
}

/// The shipping cost in cents for this method and this discounted subtotal.
pub fn shipping_cost(method: ShippingMethod, discounted_subtotal: Cents) -> Cents {
    let _ = (method, discounted_subtotal);
    todo!("shipping_cost() is not implemented yet")
}

// ---------------------------------------------------------------------------
// Unit tests.
// ---------------------------------------------------------------------------
#[cfg(test)]
mod tests {
    use super::*;
    use crate::ErrorCode;
    use rstest::rstest;

    #[test]
    fn accepts_the_two_known_methods() {
        assert_eq!(
            parse_shipping_method("standard"),
            Ok(ShippingMethod::Standard)
        );
        assert_eq!(
            parse_shipping_method("express"),
            Ok(ShippingMethod::Express)
        );
    }

    #[rstest]
    #[case("overnight")]
    #[case("Express")]
    #[case("STANDARD")]
    #[case(" standard")]
    #[case("express ")]
    #[case("")]
    fn rejects_anything_else_with_unknown_shipping_method(#[case] method: &str) {
        let error = parse_shipping_method(method).unwrap_err();

        assert_eq!(error.code, ErrorCode::UnknownShippingMethod);
    }

    #[test]
    fn the_error_message_names_the_method_it_did_not_recognise() {
        let error = parse_shipping_method("overnight").unwrap_err();

        assert_eq!(error.message, "Unknown shipping method: overnight");
    }

    #[test]
    fn standard_costs_4_99_below_50_dollars() {
        assert_eq!(shipping_cost(ShippingMethod::Standard, 0), 499);
        assert_eq!(shipping_cost(ShippingMethod::Standard, 4999), 499);
    }

    #[test]
    fn standard_is_free_from_exactly_50_dollars() {
        // The boundary on both sides: 4999 (above) and 5000 (here).
        assert_eq!(shipping_cost(ShippingMethod::Standard, 5000), 0);
        assert_eq!(shipping_cost(ShippingMethod::Standard, 5001), 0);
        assert_eq!(shipping_cost(ShippingMethod::Standard, 100000), 0);
    }

    #[rstest]
    #[case(0)]
    #[case(4999)]
    #[case(5000)]
    #[case(100000)]
    fn express_always_costs_9_99(#[case] discounted_subtotal: Cents) {
        assert_eq!(
            shipping_cost(ShippingMethod::Express, discounted_subtotal),
            999
        );
    }

    #[test]
    fn the_constants_match_the_table_in_rules_4_5() {
        assert_eq!(STANDARD_COST, 499);
        assert_eq!(EXPRESS_COST, 999);
        assert_eq!(FREE_SHIPPING_FROM, 5000);
    }
}
