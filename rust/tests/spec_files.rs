//! Integration tests for the spec loader in `tests/common/mod.rs`.
//!
//! These make sure the shared spec files load correctly *before* we rely on them in
//! `spec_cases.rs`. If someone breaks a JSON file, these tests tell you so clearly.
//!
//! Run only this file with: `cargo test --test spec_files`

// Pulls in tests/common/mod.rs.
mod common;

use checkout::Promotion;
use std::collections::HashSet;

#[test]
fn standard_catalog_has_six_products() {
    let catalog = common::load_catalog();

    assert_eq!(catalog.products.len(), 6);
}

#[test]
fn coffee_is_buy_two_get_one_free() {
    let catalog = common::load_catalog();

    // `.expect(...)` unwraps the `Option`, or fails the test with this message if it is `None`.
    let coffee = catalog
        .find("COFFEE")
        .expect("COFFEE should be in the catalog");

    assert_eq!(
        coffee.promotion,
        Some(Promotion::BuyXGetYFree { buy: 2, free: 1 })
    );
}

#[test]
fn exactly_sixty_cases_are_loaded() {
    let cases = common::load_spec_cases();

    assert_eq!(cases.len(), 60);
}

#[test]
fn case_ids_are_unique() {
    let cases = common::load_spec_cases();
    let mut seen = HashSet::new();

    for case in &cases {
        // `insert` returns false when the id was already in the set.
        let is_new = seen.insert(case.id.as_str());
        // The text after the condition is shown when the assertion fails.
        assert!(is_new, "duplicate case id {:?} in {}", case.id, case.file);
    }
}

#[test]
fn every_case_has_either_expected_totals_or_an_expected_error() {
    let cases = common::load_spec_cases();

    for case in &cases {
        let has_totals = case.expected.is_some();
        let has_error = case.expected_error.is_some();

        // `!=` on two booleans means "exactly one of them is true" (exclusive or).
        assert!(
            has_totals != has_error,
            "case {:?} in {} must have exactly one of `expected` and `expected_error`",
            case.id,
            case.file
        );
    }
}
