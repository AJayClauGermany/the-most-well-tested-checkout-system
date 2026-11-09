//! Tax (rules §4.6).

use crate::Cents;

/// Tax is 8% of the discounted subtotal.
pub const TAX_PERCENT: i64 = 8;

/// 8% of `discounted_subtotal`, rounded half up. Shipping is never passed in here.
///
/// ```
/// use checkout::tax_on;
///
/// assert_eq!(tax_on(50), 4);    // 8% of $0.50 is exactly 4 cents
/// assert_eq!(tax_on(120), 10);  // 9.6 cents rounds up
/// ```
pub fn tax_on(discounted_subtotal: Cents) -> Cents {
    let _ = discounted_subtotal;
    todo!("tax_on() is not implemented yet")
}

// ---------------------------------------------------------------------------
// Unit tests.
// ---------------------------------------------------------------------------
#[cfg(test)]
mod tests {
    use super::*;
    use rstest::rstest;

    #[test]
    fn is_8_percent() {
        assert_eq!(TAX_PERCENT, 8);
    }

    #[rstest]
    #[case::exactly_4_cents(50, 4)]
    #[case::rounds_9_point_6_up(120, 10)]
    #[case::rounds_19_point_2_down(240, 19)]
    #[case::exactly_50_cents(625, 50)]
    #[case::rounds_136_point_48_down(1706, 136)]
    #[case::rounds_950_point_4_down(11880, 950)]
    #[case::no_spend_no_tax(0, 0)]
    fn taxes_the_discounted_subtotal(#[case] amount: Cents, #[case] tax: Cents) {
        assert_eq!(tax_on(amount), tax);
    }

    #[test]
    fn rounds_plain_fractions_the_usual_way() {
        // 8% of a whole number of cents can never be exactly half a cent (x × 8 never ends
        // in 50), so the half-up rule is pinned down by the 10% coupon tests; here we check
        // plain rounding.
        assert_eq!(tax_on(1875), 150); // exactly 150.00
        assert_eq!(tax_on(1881), 150); // 150.48
        assert_eq!(tax_on(1888), 151); // 151.04
    }
}
