//! Money helpers (rules §1 and §7). All amounts are whole cents.

use crate::Cents;

/// `percent`% of `amount`, rounded half up to a whole cent (rules §7).
///
/// Done with whole numbers only, as the hint in rules §7 suggests: `amount × percent` is the
/// exact answer in hundredths of a cent; adding 50 hundredths and dividing by 100 (rounding
/// down) is the same as rounding half up. No floating point, so no `0.1 + 0.2` surprises.
///
/// ```
/// use checkout::round_half_up;
///
/// assert_eq!(round_half_up(225, 10), 23); // 22.5 rounds UP
/// assert_eq!(round_half_up(224, 10), 22); // 22.4 rounds down
/// assert_eq!(round_half_up(120, 8), 10);  // 9.6 rounds up
/// ```
pub fn round_half_up(amount: Cents, percent: i64) -> Cents {
    (amount * percent + 50) / 100
}

// ---------------------------------------------------------------------------
// Unit tests. Rounding is where checkout systems classically go wrong, so it gets its own
// small function and its own tests.
// ---------------------------------------------------------------------------
#[cfg(test)]
mod tests {
    use super::*;
    use rstest::rstest;

    // The table from rules §7 as a parameterised test: one `#[case]` row, one test.
    // rstest names them `rounds_like_the_table_in_rules_7::case_1_twenty_two_point_four` etc.
    #[rstest]
    #[case::twenty_two_point_four(224, 10, 22)]
    #[case::twenty_two_point_five(225, 10, 23)]
    #[case::twenty_two_point_six(226, 10, 23)]
    fn rounds_like_the_table_in_rules_7(
        #[case] amount: Cents,
        #[case] percent: i64,
        #[case] expected: Cents,
    ) {
        assert_eq!(round_half_up(amount, percent), expected);
    }

    #[test]
    fn rounds_an_exact_half_up_where_bankers_rounding_would_go_to_the_even_number() {
        // 10% of $52.35 is exactly 523.5 cents. Round-half-to-even would also say 524 here,
        // so the 22.5 → 23 row above (where half-to-even gives 22) pins the rule down with it.
        assert_eq!(round_half_up(5235, 10), 524);
    }

    #[test]
    fn does_not_round_when_the_result_is_already_whole() {
        assert_eq!(round_half_up(500, 10), 50);
        assert_eq!(round_half_up(50, 8), 4);
    }

    #[test]
    fn rounds_9_point_6_cents_up_to_10_and_19_point_2_cents_down_to_19() {
        assert_eq!(round_half_up(120, 8), 10);
        assert_eq!(round_half_up(240, 8), 19);
    }

    #[test]
    fn gives_0_for_an_amount_of_0_and_the_whole_amount_for_100_percent() {
        assert_eq!(round_half_up(0, 10), 0);
        assert_eq!(round_half_up(899, 100), 899);
    }

    #[test]
    fn gives_0_for_0_percent() {
        assert_eq!(round_half_up(899, 0), 0);
    }

    #[test]
    fn stays_exact_for_amounts_where_floating_point_maths_drifts() {
        // 10% of 3 cents is 0.3: it must round to 0, not 1.
        assert_eq!(round_half_up(3, 10), 0);
        // 8% of 1050 = 84.0 exactly; a float implementation might land on 83.99999.
        assert_eq!(round_half_up(1050, 8), 84);
        // Half a cent rounds up to one whole cent.
        assert_eq!(round_half_up(5, 10), 1);
    }
}
