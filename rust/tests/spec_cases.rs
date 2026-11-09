//! The shared, data-driven test cases from `spec/test-cases/*.json`.
//!
//! Each case says: "for this basket, the checkout gives these six totals" (or "fails with
//! this error code"). The same JSON files are used by every language in the repository,
//! so these 60 cases are the acceptance tests of the whole project.
//!
//! Run only this file with: `cargo test --test spec_cases`
//!
//! ## Why one test with a loop, and not one test per case?
//!
//! Rust's built-in test tool has no "parameterized tests" (one test generated per row of
//! data). So this single test loops over all cases, *collects* every mismatch together
//! with the case id, and only fails at the very end. That way one run shows you **all**
//! the broken cases, not just the first one.
//!
//! `tests/spec_cases_named.rs` shows the other way: the `libtest-mimic` crate turns each
//! case into its own named test, so `cargo test --test spec_cases_named promo-` runs only
//! the promotion cases.

mod common;

use checkout::{Catalog, checkout};
use common::{ExpectedTotals, SpecCase};

#[test]
fn all_shared_spec_cases_pass() {
    let catalog = common::load_catalog();
    let cases = common::load_spec_cases();
    let mut failures: Vec<String> = Vec::new();

    for case in &cases {
        // `if let Some(x) = ...` runs the block only when there *is* a problem.
        if let Some(problem) = check_case(case, &catalog) {
            failures.push(format!(
                "[{}] ({}) {}\n    {}",
                case.id, case.file, case.description, problem
            ));
        }
    }

    // One assertion at the end lists every failing case at once.
    assert!(
        failures.is_empty(),
        "{} of {} spec cases failed:\n\n{}\n",
        failures.len(),
        cases.len(),
        failures.join("\n\n")
    );
}

/// Runs the checkout for one case. Returns `None` when the case passes,
/// or `Some(description of the problem)` when it does not.
fn check_case(case: &SpecCase, catalog: &Catalog) -> Option<String> {
    let result = checkout(&case.basket, catalog);

    // Compare what we got with what the case expects. There are four possibilities:
    match (&case.expected, case.expected_error, result) {
        // 1. We expected totals and got a receipt: compare the six totals.
        (Some(expected), None, Ok(receipt)) => {
            let actual = ExpectedTotals::from_receipt(&receipt);
            if &actual == expected {
                None
            } else {
                Some(format!("expected {expected:?}\n    but got  {actual:?}"))
            }
        }
        // 2. We expected an error and got an error: compare the codes.
        (None, Some(expected_code), Err(error)) => {
            if error.code == expected_code {
                None
            } else {
                Some(format!(
                    "expected error {expected_code} but got error {}",
                    error.code
                ))
            }
        }
        // 3. We expected totals but the checkout failed.
        (Some(_), None, Err(error)) => {
            Some(format!("expected a receipt but got error {}", error.code))
        }
        // 4. We expected an error but the checkout succeeded.
        (None, Some(expected_code), Ok(receipt)) => Some(format!(
            "expected error {expected_code} but got a receipt with total {}",
            receipt.total
        )),
        // A broken case with both or neither (tests/spec_files.rs reports these too).
        _ => Some("the case must have exactly one of `expected` and `expected_error`".to_string()),
    }
}
