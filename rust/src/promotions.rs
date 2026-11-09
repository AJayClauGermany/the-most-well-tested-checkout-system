//! Promotions (rules §4.2): buy X get Y free, and bulk discounts.

use crate::{Cents, Promotion, round_half_up};

/// The discount a promotion gives one receipt line, in cents (rules §4.2).
/// `0` when the product has no promotion.
///
/// ```
/// use checkout::{Promotion, promotion_discount};
///
/// let buy_two_get_one = Promotion::BuyXGetYFree { buy: 2, free: 1 };
/// assert_eq!(promotion_discount(Some(&buy_two_get_one), 3, 899), 899); // one bag free
/// assert_eq!(promotion_discount(None, 3, 899), 0);
/// ```
pub fn promotion_discount(
    promotion: Option<&Promotion>,
    quantity: i64,
    unit_price: Cents,
) -> Cents {
    // `match` must handle every possibility: no promotion, or one of the two kinds.
    match promotion {
        None => 0,
        Some(Promotion::BuyXGetYFree { buy, free }) => {
            // For every complete group of `buy + free` items, `free` of them cost nothing.
            // Integer division rounds down, which is exactly "complete groups".
            let complete_groups = quantity / (buy + free);
            let free_items = complete_groups * free;
            free_items * unit_price
        }
        Some(Promotion::BulkDiscount {
            min_quantity,
            percent,
        }) => {
            // `percent`% off the whole line, but only from `min_quantity` items.
            if quantity < *min_quantity {
                0
            } else {
                round_half_up(unit_price * quantity, *percent)
            }
        }
    }
}

// ---------------------------------------------------------------------------
// Unit tests for the two kinds of promotion.
// ---------------------------------------------------------------------------
#[cfg(test)]
mod tests {
    use super::*;
    use rstest::rstest;

    fn buy_two_get_one() -> Promotion {
        Promotion::BuyXGetYFree { buy: 2, free: 1 }
    }

    fn ten_percent_off_ten() -> Promotion {
        Promotion::BulkDiscount {
            min_quantity: 10,
            percent: 10,
        }
    }

    #[test]
    fn is_0_when_the_product_has_no_promotion() {
        assert_eq!(promotion_discount(None, 50, 899), 0);
    }

    // Coffee is buy 2 get 1 free at $8.99: a table of quantity → free items.
    #[rstest]
    #[case(1, 0)]
    #[case(2, 0)]
    #[case(3, 1)]
    #[case(4, 1)]
    #[case(5, 1)]
    #[case(6, 2)]
    #[case(99, 33)]
    fn buy_x_get_y_free_gives_free_items_per_complete_group(
        #[case] quantity: i64,
        #[case] free: i64,
    ) {
        assert_eq!(
            promotion_discount(Some(&buy_two_get_one()), quantity, 899),
            free * 899
        );
    }

    #[rstest]
    // Buy 1 get 1 free: every second item is free.
    #[case(1, 1, 4, 200)]
    #[case(1, 1, 5, 200)]
    // Buy 3 get 2 free: groups of 5, 2 of each group free.
    #[case(3, 2, 4, 0)]
    #[case(3, 2, 5, 200)]
    #[case(3, 2, 10, 400)]
    fn buy_x_get_y_free_uses_the_promotions_own_numbers_not_2_and_1(
        #[case] buy: i64,
        #[case] free: i64,
        #[case] quantity: i64,
        #[case] discount: Cents,
    ) {
        let promotion = Promotion::BuyXGetYFree { buy, free };

        assert_eq!(
            promotion_discount(Some(&promotion), quantity, 100),
            discount
        );
    }

    #[test]
    fn buy_x_get_y_free_multiplies_by_the_unit_price() {
        assert_eq!(promotion_discount(Some(&buy_two_get_one()), 3, 100), 100);
        assert_eq!(promotion_discount(Some(&buy_two_get_one()), 3, 1), 1);
    }

    // Notebooks are 10% off from 10 at $3.49. Boundaries 9 / 10 on both sides.
    #[rstest]
    #[case::one_short_of_the_minimum(9, 0)]
    #[case::exactly_the_minimum_10_percent_of_34_90(10, 349)]
    #[case::rounds_383_point_9_cents_up(11, 384)]
    #[case::rounds_exactly_523_point_5_cents_up(15, 524)]
    #[case::ten_percent_of_345_51_is_3455_point_1(99, 3455)]
    fn bulk_discount_takes_the_percentage_off_from_the_minimum_quantity(
        #[case] quantity: i64,
        #[case] discount: Cents,
    ) {
        assert_eq!(
            promotion_discount(Some(&ten_percent_off_ten()), quantity, 349),
            discount
        );
    }

    #[test]
    fn bulk_discount_uses_the_promotions_own_minimum_and_percentage() {
        let quarter_off_from_three = Promotion::BulkDiscount {
            min_quantity: 3,
            percent: 25,
        };

        assert_eq!(promotion_discount(Some(&quarter_off_from_three), 2, 100), 0);
        assert_eq!(
            promotion_discount(Some(&quarter_off_from_three), 3, 100),
            75
        );
        assert_eq!(
            promotion_discount(Some(&quarter_off_from_three), 3, 101),
            76
        ); // 75.75 → 76
    }
}
