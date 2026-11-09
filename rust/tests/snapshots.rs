//! Snapshot tests with `insta`: did the output change since last time?
//!
//! `assert_json_snapshot!(value)` saves the value as JSON to `tests/snapshots/*.snap` on the
//! first run and compares against that file on later runs. A snapshot written inline
//! (`@r#"..."#`) keeps the approved value right here in the test. When a change is
//! intended, review and accept it with `cargo insta review` (or `cargo insta accept`);
//! without cargo-insta installed, `INSTA_UPDATE=always cargo test` rewrites the files.
//!
//! The hand-checked receipt in `spec/examples/receipt.json` plays the role of an
//! "approved" output shared by all six languages.
//!
//! Run only this file with: `cargo test --test snapshots`

mod common;

use checkout::{checkout, run_checkout_command};
use common::fixtures::{basket, example_basket, full_everything_basket};
use insta::{assert_json_snapshot, assert_snapshot};
use pretty_assertions::assert_eq;

#[test]
fn the_example_receipt_matches_the_hand_checked_spec_examples_receipt_json_lines_included() {
    let catalog = common::load_catalog();
    let approved = common::load_example("receipt.json");

    let receipt = checkout(&example_basket(), &catalog).unwrap();

    assert_eq!(serde_json::to_value(&receipt).unwrap(), approved);
}

#[test]
fn the_example_receipt_matches_the_saved_snapshot() {
    let catalog = common::load_catalog();

    let receipt = checkout(&example_basket(), &catalog).unwrap();

    // Compares with tests/snapshots/snapshots__the_example_receipt_matches_the_saved_snapshot.snap
    assert_json_snapshot!(receipt);
}

#[test]
fn the_example_receipt_is_printed_by_the_command_exactly_like_this() {
    let catalog = common::load_catalog();
    let input = serde_json::to_string(&example_basket()).unwrap();

    let result = run_checkout_command(&input, &catalog).unwrap();

    assert_snapshot!(result.output);
}

#[test]
fn the_biggest_shared_case_full_everything_matches_the_saved_snapshot() {
    let catalog = common::load_catalog();

    let receipt = checkout(&full_everything_basket(), &catalog).unwrap();

    assert_json_snapshot!(receipt);
}

#[test]
fn one_pen_written_out_as_receipt_lines_inline_snapshot() {
    let catalog = common::load_catalog();

    let receipt = checkout(&basket(&[("PEN", 1)]), &catalog).unwrap();

    assert_json_snapshot!(receipt.lines, @r#"
    [
      {
        "sku": "PEN",
        "name": "Pen",
        "quantity": 1,
        "unit_price": 120,
        "line_total": 120,
        "line_discount": 0
      }
    ]
    "#);
}

#[test]
fn an_error_as_the_command_prints_it_inline_snapshot() {
    let catalog = common::load_catalog();

    let result = run_checkout_command(r#"{"items":[]}"#, &catalog).unwrap();

    assert_eq!(result.exit_code, 1);
    assert_snapshot!(result.output, @r#"{"error":"EMPTY_BASKET"}"#);
}
