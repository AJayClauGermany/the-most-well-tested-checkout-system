//! The shared spec cases again – this time **one named test per case**.
//!
//! `tests/spec_cases.rs` runs all 60 cases inside one `#[test]`. That is simple, but the
//! test runner only ever sees one test. Here the `libtest-mimic` crate builds the list of
//! tests at run time (one `Trial` per case, named after the case id) and runs them with
//! the same command line options and the same output as the built-in harness. So:
//!
//! ```text
//! cargo test --test spec_cases_named               # all 60, each reported by name
//! cargo test --test spec_cases_named promo-        # only the promotion cases
//! cargo test --test spec_cases_named -- --list     # just print the names
//! ```
//!
//! `harness = false` in Cargo.toml tells Cargo this file has its own `main`.

mod common;

use checkout::{Catalog, checkout};
use common::{ExpectedTotals, SpecCase};
use libtest_mimic::{Arguments, Failed, Trial};
use std::process::ExitCode;
use std::sync::Arc;

fn main() -> ExitCode {
    let arguments = Arguments::from_args();
    let catalog = Arc::new(common::load_catalog());

    let trials: Vec<Trial> = common::load_spec_cases()
        .into_iter()
        .map(|case| {
            let catalog = Arc::clone(&catalog);
            // The "kind" is shown in brackets before the name: `[01-subtotal] subtotal-single-item`.
            let kind = case.file.trim_end_matches(".json").to_string();
            Trial::test(case.id.clone(), move || check_case(&case, &catalog)).with_kind(kind)
        })
        .collect();

    libtest_mimic::run(&arguments, trials).exit_code()
}

/// Runs the checkout for one case; `Err` describes what went wrong.
fn check_case(case: &SpecCase, catalog: &Catalog) -> Result<(), Failed> {
    let result = checkout(&case.basket, catalog);

    match (&case.expected, case.expected_error, result) {
        (Some(expected), None, Ok(receipt)) => {
            let actual = ExpectedTotals::from_receipt(&receipt);
            if &actual == expected {
                Ok(())
            } else {
                Err(format!(
                    "{}\nexpected {expected:?}\n but got {actual:?}",
                    case.description
                )
                .into())
            }
        }
        (None, Some(expected_code), Err(error)) if error.code == expected_code => Ok(()),
        (None, Some(expected_code), Err(error)) => Err(format!(
            "{}\nexpected error {expected_code} but got error {}",
            case.description, error.code
        )
        .into()),
        (Some(_), None, Err(error)) => Err(format!(
            "{}\nexpected a receipt but got error {}",
            case.description, error.code
        )
        .into()),
        (None, Some(expected_code), Ok(receipt)) => Err(format!(
            "{}\nexpected error {expected_code} but got a receipt with total {}",
            case.description, receipt.total
        )
        .into()),
        _ => Err("the case must have exactly one of `expected` and `expected_error`".into()),
    }
}
