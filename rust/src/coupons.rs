//! Coupons (rules §4.4).

use crate::{Cents, CheckoutError};

/// The coupon code as the shop knows it: trimmed and upper case.
/// Missing, empty or blank all mean "no coupon" and give `None`.
///
/// ```
/// use checkout::normalize_coupon_code;
///
/// assert_eq!(normalize_coupon_code(Some("  save10 ")), Some("SAVE10".to_string()));
/// assert_eq!(normalize_coupon_code(Some("   ")), None);
/// assert_eq!(normalize_coupon_code(None), None);
/// ```
pub fn normalize_coupon_code(coupon: Option<&str>) -> Option<String> {
    let _ = coupon;
    todo!("normalize_coupon_code() is not implemented yet")
}

/// The coupon discount in cents for `after_promotions` cents of shopping.
///
/// Fails with `UNKNOWN_COUPON` or `COUPON_MINIMUM_NOT_MET` (rules §4.1 rows 6 and 7).
pub fn coupon_discount(
    coupon: Option<&str>,
    after_promotions: Cents,
) -> Result<Cents, CheckoutError> {
    let _ = (coupon, after_promotions);
    todo!("coupon_discount() is not implemented yet")
}

// ---------------------------------------------------------------------------
// Unit tests.
// ---------------------------------------------------------------------------
#[cfg(test)]
mod tests {
    use super::*;
    use crate::ErrorCode;
    use rstest::rstest;

    #[rstest]
    #[case("SAVE10", "SAVE10")]
    #[case("save10", "SAVE10")]
    #[case("  save10 ", "SAVE10")]
    #[case("\tWelcome5\n", "WELCOME5")]
    #[case("nope", "NOPE")]
    fn normalize_trims_and_upper_cases(#[case] input: &str, #[case] expected: &str) {
        assert_eq!(
            normalize_coupon_code(Some(input)),
            Some(expected.to_string())
        );
    }

    #[rstest]
    #[case(None)]
    #[case(Some(""))]
    #[case(Some("   "))]
    #[case(Some("\t\n"))]
    fn normalize_treats_missing_empty_and_blank_as_no_coupon(#[case] input: Option<&str>) {
        assert_eq!(normalize_coupon_code(input), None);
    }

    #[rstest]
    #[case(None)]
    #[case(Some(""))]
    #[case(Some("  "))]
    fn discount_is_0_without_a_coupon(#[case] coupon: Option<&str>) {
        assert_eq!(coupon_discount(coupon, 10000), Ok(0));
    }

    #[rstest]
    #[case(500, 50)]
    #[case(225, 23)] // 22.5 → 23
    #[case(224, 22)] // 22.4 → 22
    #[case(1798, 180)] // 179.8 → 180
    fn save10_takes_10_percent_off_rounded_half_up(
        #[case] after_promotions: Cents,
        #[case] discount: Cents,
    ) {
        assert_eq!(
            coupon_discount(Some("SAVE10"), after_promotions),
            Ok(discount)
        );
    }

    #[test]
    fn save10_has_no_minimum_spend() {
        assert_eq!(coupon_discount(Some("SAVE10"), 1), Ok(0)); // 0.1 → 0
        assert_eq!(coupon_discount(Some("SAVE10"), 5), Ok(1)); // 0.5 → 1
    }

    #[test]
    fn save10_ignores_case_and_surrounding_spaces() {
        assert_eq!(coupon_discount(Some(" save10 "), 500), Ok(50));
    }

    #[test]
    fn welcome5_takes_exactly_5_dollars_off() {
        assert_eq!(coupon_discount(Some("WELCOME5"), 2500), Ok(500));
        assert_eq!(coupon_discount(Some("WELCOME5"), 99999), Ok(500));
    }

    #[test]
    fn welcome5_works_at_exactly_the_25_dollar_minimum_and_fails_one_cent_below_it() {
        // The boundary, tested on both sides: 2499 / 2500.
        assert_eq!(coupon_discount(Some("WELCOME5"), 2500), Ok(500));

        let error = coupon_discount(Some("WELCOME5"), 2499).unwrap_err();
        assert_eq!(error.code, ErrorCode::CouponMinimumNotMet);
    }

    #[test]
    fn welcome5_ignores_case_and_surrounding_spaces() {
        assert_eq!(coupon_discount(Some("welcome5"), 2500), Ok(500));
    }

    #[test]
    fn rejects_a_code_it_does_not_know_with_unknown_coupon() {
        let error = coupon_discount(Some("FREESTUFF"), 10000).unwrap_err();

        assert_eq!(error.code, ErrorCode::UnknownCoupon);
    }

    #[test]
    fn explains_what_went_wrong_in_the_error_message() {
        assert_eq!(
            coupon_discount(Some("freestuff"), 10000)
                .unwrap_err()
                .message,
            "Unknown coupon: FREESTUFF"
        );
        assert_eq!(
            coupon_discount(Some("WELCOME5"), 2499).unwrap_err().message,
            "WELCOME5 needs a spend of at least 2500 cents, got 2499"
        );
    }

    #[rstest]
    #[case("SAVE100")]
    #[case("WELCOME50")]
    #[case("SAVE")]
    #[case("SAVE 10")]
    fn does_not_accept_a_code_that_merely_resembles_a_real_code(#[case] code: &str) {
        let error = coupon_discount(Some(code), 10000).unwrap_err();

        assert_eq!(error.code, ErrorCode::UnknownCoupon);
    }
}
