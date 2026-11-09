//! Receipt lines (rules §4.2): merging basket items and pricing a line.

use crate::{BasketItem, Product, ReceiptLine};

/// One item per distinct SKU, in order of first appearance, with the quantities added up
/// (rules §4.2). `[APPLE x1, PEN x2, APPLE x3]` gives `[APPLE x4, PEN x2]`.
///
/// ```
/// use checkout::{BasketItem, merge_items};
///
/// let item = |sku: &str, quantity| BasketItem { sku: sku.to_string(), quantity };
/// let merged = merge_items(&[item("APPLE", 1), item("PEN", 2), item("APPLE", 3)]);
///
/// assert_eq!(merged, vec![item("APPLE", 4), item("PEN", 2)]);
/// ```
pub fn merge_items(items: &[BasketItem]) -> Vec<BasketItem> {
    let _ = items;
    todo!("merge_items() is not implemented yet")
}

/// The receipt line for `quantity` of `product`, including its promotion discount.
pub fn price_line(product: &Product, quantity: i64) -> ReceiptLine {
    let _ = (product, quantity);
    todo!("price_line() is not implemented yet")
}

// ---------------------------------------------------------------------------
// Unit tests.
// ---------------------------------------------------------------------------
#[cfg(test)]
mod tests {
    use super::*;
    use crate::test_fixtures::{apple, coffee, item, notebook};
    use pretty_assertions::assert_eq;

    #[test]
    fn keeps_items_with_different_skus_apart_in_basket_order() {
        let merged = merge_items(&[item("PEN", 2), item("APPLE", 1)]);

        assert_eq!(merged, vec![item("PEN", 2), item("APPLE", 1)]);
    }

    #[test]
    fn adds_up_the_quantities_of_the_same_sku_keeping_the_place_of_its_first_appearance() {
        // The example from rules §4.2.
        let merged = merge_items(&[item("APPLE", 1), item("PEN", 2), item("APPLE", 3)]);

        assert_eq!(merged, vec![item("APPLE", 4), item("PEN", 2)]);
    }

    #[test]
    fn merges_three_lines_of_one_sku_into_one() {
        let merged = merge_items(&[item("COFFEE", 1), item("COFFEE", 1), item("COFFEE", 1)]);

        assert_eq!(merged, vec![item("COFFEE", 3)]);
    }

    #[test]
    fn gives_an_empty_list_for_an_empty_list() {
        assert_eq!(merge_items(&[]), Vec::<BasketItem>::new());
    }

    #[test]
    fn does_not_change_the_items_it_is_given() {
        // `merge_items` borrows the slice, so the compiler guarantees this; the test shows it.
        let items = [item("APPLE", 1), item("APPLE", 1)];

        let merged = merge_items(&items);

        assert_eq!(items, [item("APPLE", 1), item("APPLE", 1)]);
        assert_eq!(merged, vec![item("APPLE", 2)]);
    }

    #[test]
    fn copies_sku_name_and_unit_price_from_the_product_and_multiplies_for_the_line_total() {
        let line = price_line(&apple(), 3);

        assert_eq!(
            line,
            ReceiptLine {
                sku: "APPLE".to_string(),
                name: "Apple".to_string(),
                quantity: 3,
                unit_price: 50,
                line_total: 150,
                line_discount: 0,
            }
        );
    }

    #[test]
    fn includes_the_promotion_discount() {
        let coffee_line = price_line(&coffee(), 3);
        let notebook_line = price_line(&notebook(), 10);

        assert_eq!(
            (coffee_line.line_total, coffee_line.line_discount),
            (2697, 899)
        );
        assert_eq!(
            (notebook_line.line_total, notebook_line.line_discount),
            (3490, 349)
        );
    }

    #[test]
    fn gives_no_discount_below_the_promotions_threshold() {
        assert_eq!(price_line(&coffee(), 2).line_discount, 0);
        assert_eq!(price_line(&notebook(), 9).line_discount, 0);
    }
}
