//! Test code deserves tests too: the BDD steps rely on `dollars_to_cents`, so a bug there
//! would hide bugs in the checkout. These are plain unit tests of a test helper.

mod common;

use common::dollars::dollars_to_cents;
use rstest::rstest;

// A parameterised test: the same check runs once per `#[case]` row.
#[rstest]
#[case("$0.00", 0)]
#[case("$0.50", 50)]
#[case("$12.39", 1239)]
#[case("$5", 500)]
#[case("$5.5", 550)]
#[case("  $8.99 ", 899)]
#[case("0.04", 4)]
#[case("$123456.78", 12345678)]
fn turns_dollars_into_cents(#[case] text: &str, #[case] cents: i64) {
    assert_eq!(dollars_to_cents(text), Ok(cents));
}

#[rstest]
#[case("ten dollars")]
#[case("$1.234")]
#[case("$")]
#[case("")]
#[case("$1.")]
#[case("$-1.00")]
#[case("$1,00")]
fn refuses_text_that_is_not_an_amount(#[case] text: &str) {
    let error = dollars_to_cents(text).unwrap_err();

    assert!(error.starts_with("not a dollar amount"), "{error}");
}
