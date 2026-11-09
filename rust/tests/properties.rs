//! Property-based tests with `proptest`.
//!
//! An example test says "THIS basket gives THAT total". A property says something that must
//! be true for EVERY basket, and proptest generates hundreds of random baskets trying to
//! break it – then shrinks any failure to the smallest example. The properties are listed
//! once for all six languages in `spec/properties.md`; each test is named after its id.
//!
//! Run only this file with: `cargo test --test properties`

mod common;

use checkout::{
    Basket, BasketItem, Catalog, ErrorCode, FixedClock, InMemoryOrderRepository, OrderRepository,
    Receipt, checkout, format_timestamp, place_order, round_half_up,
};
use common::schemas::schema_errors;
use proptest::prelude::*;
use std::collections::HashSet;
use std::sync::LazyLock;
use std::time::{Duration, UNIX_EPOCH};

/// The standard catalog, loaded once for the whole test program.
static CATALOG: LazyLock<Catalog> = LazyLock::new(common::load_catalog);

// ---------------------------------------------------------------------------
// Generators ("strategies"), following the "Generating baskets" section of properties.md
// ---------------------------------------------------------------------------

/// One of the six real SKUs.
fn real_sku() -> impl Strategy<Value = String> {
    let skus: Vec<String> = CATALOG.products.iter().map(|p| p.sku.clone()).collect();
    prop::sample::select(skus)
}

/// A sensible item: a real SKU, a quantity from 1 to 99.
fn sensible_item() -> impl Strategy<Value = BasketItem> {
    (real_sku(), 1..=99i64).prop_map(|(sku, quantity)| BasketItem { sku, quantity })
}

/// SAVE10 or WELCOME5 written any old way: mixed case, spaces around it.
fn sloppy_coupon() -> impl Strategy<Value = String> {
    (
        prop::sample::select(vec!["SAVE10", "WELCOME5"]),
        prop::collection::vec(any::<bool>(), 8),
        0..4usize,
        0..4usize,
    )
        .prop_map(|(code, upper, before, after)| sloppy_spelling(code, &upper, before, after))
}

/// `code` with each letter upper or lower case as `upper` says, padded with spaces.
fn sloppy_spelling(code: &str, upper: &[bool], before: usize, after: usize) -> String {
    let letters: String = code
        .chars()
        .zip(upper)
        .map(|(letter, &upper)| {
            if upper {
                letter.to_ascii_uppercase()
            } else {
                letter.to_ascii_lowercase()
            }
        })
        .collect();
    format!("{}{letters}{}", " ".repeat(before), " ".repeat(after))
}

/// Missing, a real code, a sloppy real code, or a random string.
fn maybe_coupon() -> impl Strategy<Value = Option<String>> {
    prop_oneof![
        3 => Just(None),
        2 => Just(Some("SAVE10".to_string())),
        2 => Just(Some("WELCOME5".to_string())),
        2 => sloppy_coupon().prop_map(Some),
        1 => "[A-Za-z0-9 ]{0,8}".prop_map(Some),
    ]
}

/// Missing (which the Basket type spells "standard"), a real method, or nonsense.
fn shipping() -> impl Strategy<Value = String> {
    prop::sample::select(vec![
        "standard",
        "standard",
        "express",
        "overnight",
        "Express",
        "",
    ])
    .prop_map(str::to_string)
}

/// A whole basket: 1 to 8 sensible items (the same SKU may repeat), maybe a coupon, maybe
/// a shipping method. Many of these are still invalid (too many of one SKU, a bad coupon).
fn basket() -> impl Strategy<Value = Basket> {
    (
        prop::collection::vec(sensible_item(), 1..=8),
        maybe_coupon(),
        shipping(),
    )
        .prop_map(|(items, coupon, shipping)| Basket {
            items,
            coupon,
            shipping,
        })
}

/// Only baskets that check out successfully – what P1 to P9 are about.
fn valid_basket() -> impl Strategy<Value = Basket> {
    basket().prop_filter("the basket must check out", |b| {
        checkout(b, &CATALOG).is_ok()
    })
}

/// Anything at all that still fits the Basket type: nonsense SKUs, silly quantities,
/// random coupon and shipping strings.
fn any_basket() -> impl Strategy<Value = Basket> {
    let any_item = (
        prop_oneof![real_sku(), "[A-Za-z]{0,6}"],
        prop_oneof![-5..150i64, any::<i64>()],
    )
        .prop_map(|(sku, quantity)| BasketItem { sku, quantity });
    (
        prop::collection::vec(any_item, 0..=6),
        prop::option::of(prop_oneof![
            prop::sample::select(vec!["SAVE10", "WELCOME5", "save10", " welcome5 "])
                .prop_map(str::to_string),
            ".{0,10}",
        ]),
        prop_oneof![
            prop::sample::select(vec!["standard", "express"]).prop_map(str::to_string),
            ".{0,10}",
        ],
    )
        .prop_map(|(items, coupon, shipping)| Basket {
            items,
            coupon,
            shipping,
        })
}

/// The six totals of a receipt (what P4 compares; line order may differ).
fn totals_of(receipt: &Receipt) -> [i64; 6] {
    [
        receipt.subtotal,
        receipt.promotion_discount,
        receipt.coupon_discount,
        receipt.shipping,
        receipt.tax,
        receipt.total,
    ]
}

fn is_checkout_error(code: ErrorCode) -> bool {
    ErrorCode::ALL.contains(&code)
}

// ---------------------------------------------------------------------------
// Properties of a successful checkout
// ---------------------------------------------------------------------------

proptest! {
    #[test]
    fn p1_total_equals_subtotal_minus_discounts_plus_shipping_and_tax(basket in valid_basket()) {
        let r = checkout(&basket, &CATALOG).unwrap();

        prop_assert_eq!(
            r.total,
            r.subtotal - r.promotion_discount - r.coupon_discount + r.shipping + r.tax
        );
    }

    #[test]
    fn p2_nothing_is_negative_and_a_discount_never_exceeds_what_it_discounts(basket in valid_basket()) {
        let r = checkout(&basket, &CATALOG).unwrap();

        for amount in totals_of(&r) {
            prop_assert!(amount >= 0, "negative amount {amount} in {r:?}");
        }
        for line in &r.lines {
            prop_assert!(line.unit_price >= 0);
            prop_assert!(line.line_total >= 0);
            prop_assert!(line.line_discount >= 0);
            prop_assert!(line.line_discount <= line.line_total);
        }
        prop_assert!(r.promotion_discount <= r.subtotal);
        prop_assert!(r.coupon_discount <= r.subtotal - r.promotion_discount);
    }

    #[test]
    fn p3_the_receipt_lines_explain_the_totals(basket in valid_basket()) {
        let r = checkout(&basket, &CATALOG).unwrap();

        prop_assert_eq!(r.subtotal, r.lines.iter().map(|line| line.line_total).sum::<i64>());
        prop_assert_eq!(r.promotion_discount, r.lines.iter().map(|line| line.line_discount).sum::<i64>());
        for line in &r.lines {
            prop_assert_eq!(line.line_total, line.unit_price * line.quantity);
        }
    }

    #[test]
    fn p4_shuffling_the_items_gives_the_same_totals(
        (basket, shuffled_items) in valid_basket().prop_flat_map(|basket| {
            let items = basket.items.clone();
            (Just(basket), Just(items).prop_shuffle())
        })
    ) {
        let shuffled = Basket { items: shuffled_items, ..basket.clone() };

        prop_assert_eq!(
            totals_of(&checkout(&shuffled, &CATALOG).unwrap()),
            totals_of(&checkout(&basket, &CATALOG).unwrap())
        );
    }

    #[test]
    fn p5_splitting_one_line_into_two_lines_of_the_same_sku_gives_the_same_receipt(
        basket in valid_basket().prop_filter("needs a line with quantity >= 2", |b| {
            b.items.iter().any(|item| item.quantity >= 2)
        }),
        pick in any::<prop::sample::Index>(),
        cut in any::<prop::sample::Index>(),
    ) {
        let candidates: Vec<usize> = (0..basket.items.len())
            .filter(|&i| basket.items[i].quantity >= 2)
            .collect();
        let index = *pick.get(&candidates);
        let line = &basket.items[index];
        let first = 1 + cut.index((line.quantity - 1) as usize) as i64; // 1 .. quantity-1

        let mut split = basket.items.clone();
        split[index].quantity = first;
        split.push(BasketItem { sku: line.sku.clone(), quantity: line.quantity - first });
        let split_basket = Basket { items: split, ..basket.clone() };

        prop_assert_eq!(checkout(&split_basket, &CATALOG), checkout(&basket, &CATALOG));
    }

    #[test]
    fn p6_the_shipping_table_is_applied_exactly(basket in valid_basket()) {
        let r = checkout(&basket, &CATALOG).unwrap();

        prop_assert!([0, 499, 999].contains(&r.shipping), "shipping was {}", r.shipping);
        let discounted = r.subtotal - r.promotion_discount - r.coupon_discount;
        if basket.shipping == "express" {
            prop_assert_eq!(r.shipping, 999);
        } else {
            prop_assert_eq!(r.shipping, if discounted >= 5000 { 0 } else { 499 });
        }
    }

    #[test]
    fn p7_tax_is_8_percent_of_the_discounted_subtotal_rounded_half_up(basket in valid_basket()) {
        let r = checkout(&basket, &CATALOG).unwrap();

        let discounted = r.subtotal - r.promotion_discount - r.coupon_discount;
        prop_assert_eq!(r.tax, (discounted * 8 + 50) / 100);
    }

    #[test]
    fn p8_writing_the_coupon_in_another_case_or_with_spaces_gives_the_same_receipt(
        basket in valid_basket().prop_filter("needs a real coupon", |b| {
            matches!(b.coupon.as_deref(), Some("SAVE10") | Some("WELCOME5"))
        }),
        upper in prop::collection::vec(any::<bool>(), 8),
        before in 0..4usize,
        after in 0..4usize,
    ) {
        let code = basket.coupon.clone().unwrap();
        let sloppy = Basket { coupon: Some(sloppy_spelling(&code, &upper, before, after)), ..basket.clone() };

        prop_assert_eq!(checkout(&sloppy, &CATALOG), checkout(&basket, &CATALOG));
    }

    #[test]
    fn p9_one_line_per_distinct_sku_in_order_of_first_appearance_each_with_quantity_1_to_99(
        basket in valid_basket()
    ) {
        let r = checkout(&basket, &CATALOG).unwrap();

        let mut seen = HashSet::new();
        let distinct_in_order: Vec<&str> = basket
            .items
            .iter()
            .map(|item| item.sku.as_str())
            .filter(|sku| seen.insert(*sku))
            .collect();
        let line_skus: Vec<&str> = r.lines.iter().map(|line| line.sku.as_str()).collect();
        prop_assert_eq!(line_skus, distinct_in_order);

        for line in &r.lines {
            prop_assert!((1..=99).contains(&line.quantity), "quantity {}", line.quantity);
            let expected: i64 = basket
                .items
                .iter()
                .filter(|item| item.sku == line.sku)
                .map(|item| item.quantity)
                .sum();
            prop_assert_eq!(line.quantity, expected);
        }
    }
}

// ---------------------------------------------------------------------------
// Properties of any input at all
// ---------------------------------------------------------------------------

proptest! {
    // More cases than the default 256: this is the property fuzzing is built on.
    #![proptest_config(ProptestConfig::with_cases(500))]

    #[test]
    fn p10_for_any_basket_checkout_gives_a_schema_valid_receipt_or_one_of_the_six_error_codes(
        basket in any_basket()
    ) {
        // `checkout` returns a Result, so "never crashes" means: it returns at all. The call
        // happening without a panic is the first half of the property.
        match checkout(&basket, &CATALOG) {
            Ok(receipt) => {
                let as_json = serde_json::to_value(&receipt).unwrap();
                prop_assert_eq!(schema_errors("receipt.schema.json", &as_json), Vec::<String>::new());
            }
            Err(error) => prop_assert!(is_checkout_error(error.code)),
        }
    }

    #[test]
    fn p11_round_half_up_is_within_half_a_cent_and_exact_when_the_value_is_whole(
        amount in 0..=10_000_000i64,
        percent in 0..=100i64,
    ) {
        let rounded = round_half_up(amount, percent);
        let exact_times_100 = amount * percent; // the exact value is this ÷ 100

        prop_assert!((rounded * 100 - exact_times_100).abs() <= 50);
        if exact_times_100 % 100 == 0 {
            prop_assert_eq!(rounded * 100, exact_times_100);
        }
        if exact_times_100 % 100 == 50 {
            prop_assert_eq!(rounded * 100, exact_times_100 + 50); // halves go UP
        }
    }
}

// ---------------------------------------------------------------------------
// Properties of saving orders
// ---------------------------------------------------------------------------

fn fixed_clock() -> Box<FixedClock> {
    Box::new(FixedClock::at(common::FIXED_TIME).unwrap())
}

proptest! {
    #[test]
    fn p12_after_saving_n_orders_the_ids_are_exactly_order_1_to_order_n_in_that_order(
        baskets in prop::collection::vec(valid_basket(), 0..=20)
    ) {
        let mut repository = InMemoryOrderRepository::new(fixed_clock());

        let ids: Vec<String> = baskets
            .iter()
            .map(|basket| place_order(basket, &CATALOG, &mut repository).unwrap().order_id)
            .collect();

        let expected: Vec<String> = (1..=baskets.len()).map(|n| format!("ORDER-{n}")).collect();
        prop_assert_eq!(ids, expected);
    }

    #[test]
    fn p13_find_of_save_gives_back_the_same_receipt_and_the_clocks_time(
        basket in valid_basket(),
        // Any whole second this century.
        unix_seconds in 946_684_800u64..=4_102_444_799u64,
    ) {
        let when = UNIX_EPOCH + Duration::from_secs(unix_seconds);
        let mut repository = InMemoryOrderRepository::new(Box::new(FixedClock::new(when)));
        let receipt = checkout(&basket, &CATALOG).unwrap();

        let id = repository.save(&receipt).unwrap();
        let order = repository.find(&id).unwrap().expect("the order should be found");

        prop_assert_eq!(order.receipt, receipt);
        prop_assert_eq!(order.created_at, format_timestamp(when));
    }

    #[test]
    fn p14_a_failed_checkout_saves_nothing_and_find_of_an_unknown_id_is_not_found(
        basket in any_basket().prop_filter("the basket must be invalid", |b| checkout(b, &CATALOG).is_err()),
        random_id in ".{0,12}",
    ) {
        let mut repository = InMemoryOrderRepository::new(fixed_clock());

        let error = place_order(&basket, &CATALOG, &mut repository).unwrap_err();

        prop_assert!(error.checkout_code().is_some(), "not a checkout error: {error}");
        prop_assert_eq!(repository.count(), 0);
        prop_assert_eq!(repository.find(&random_id), Ok(None));
        prop_assert_eq!(repository.find("ORDER-1"), Ok(None));
    }
}
