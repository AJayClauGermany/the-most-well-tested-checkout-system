//! Shared helpers for the integration tests: loading the language-agnostic spec, small
//! fixtures, schema validation and the `$12.39` → `1239` helper.
//!
//! Why `tests/common/mod.rs` and not `tests/common.rs`? Cargo treats every `.rs` file
//! directly inside `tests/` as its own test program. Files in a sub-folder are not,
//! so this is the standard place for code that several test files share.
//! A test file uses it with `mod common;`.
//!
//! The spec files are read straight from the repository's `spec/` folder (never copied),
//! so every language is checked against exactly the same data.

// Each test file only uses some of these helpers. Without this line the compiler
// would warn about the unused ones (and our CI treats warnings as errors).
#![allow(dead_code)]

pub mod dollars;
pub mod fixtures;
pub mod schemas;

use checkout::{Basket, Catalog, Cents, Clock, ErrorCode, FixedClock, Receipt};
use serde::{Deserialize, Serialize};
use serde_json::Value;
use std::fs;
use std::path::{Path, PathBuf};

/// The time every fixed clock in the tests answers (the one from `spec/examples/order.json`).
pub const FIXED_TIME: &str = "2026-01-15T10:30:00Z";

/// A stub clock that always says [`FIXED_TIME`].
pub fn fixed_clock() -> Box<dyn Clock> {
    Box::new(FixedClock::at(FIXED_TIME).expect("FIXED_TIME should be a valid timestamp"))
}

/// The six totals of a receipt, as written in a test case's `"expected"` field.
#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
pub struct ExpectedTotals {
    pub subtotal: Cents,
    pub promotion_discount: Cents,
    pub coupon_discount: Cents,
    pub shipping: Cents,
    pub tax: Cents,
    pub total: Cents,
}

impl ExpectedTotals {
    /// Picks the six totals out of a receipt, so they can be compared with `assert_eq!`.
    pub fn from_receipt(receipt: &Receipt) -> Self {
        ExpectedTotals {
            subtotal: receipt.subtotal,
            promotion_discount: receipt.promotion_discount,
            coupon_discount: receipt.coupon_discount,
            shipping: receipt.shipping,
            tax: receipt.tax,
            total: receipt.total,
        }
    }
}

/// One shared test case from `spec/test-cases/*.json`.
///
/// A valid case has **exactly one** of `expected` and `expected_error`
/// (`tests/spec_files.rs` checks that).
#[derive(Debug, Clone, Deserialize)]
pub struct SpecCase {
    /// The file the case came from, e.g. `"02-promotions.json"`. Not in the JSON itself:
    /// the loader fills it in, so failure messages can say where to look.
    #[serde(skip)]
    pub file: String,
    pub id: String,
    pub description: String,
    pub basket: Basket,
    pub expected: Option<ExpectedTotals>,
    pub expected_error: Option<ErrorCode>,
}

/// The shape of a whole test-case file: a description plus a list of cases.
#[derive(Deserialize)]
struct SpecFile {
    cases: Vec<SpecCase>,
}

/// The repository's `spec/` folder.
///
/// `env!("CARGO_MANIFEST_DIR")` is the folder containing `Cargo.toml` (that is, `rust/`),
/// filled in by Cargo at compile time. Using it means the tests find the spec no matter
/// which folder you run `cargo test` from.
pub fn spec_dir() -> PathBuf {
    Path::new(env!("CARGO_MANIFEST_DIR"))
        .join("..")
        .join("spec")
}

/// Loads the standard catalog from `spec/catalog.json`.
pub fn load_catalog() -> Catalog {
    let path = spec_dir().join("catalog.json");
    let text = read_file(&path);
    serde_json::from_str(&text)
        .unwrap_or_else(|error| panic!("could not parse {}: {error}", path.display()))
}

/// Loads every case from every `spec/test-cases/*.json` file, sorted by file name.
pub fn load_spec_cases() -> Vec<SpecCase> {
    let mut all_cases = Vec::new();

    for path in test_case_files() {
        let file_name = path.file_name().unwrap().to_string_lossy().to_string();
        let text = read_file(&path);
        let spec_file: SpecFile = serde_json::from_str(&text)
            .unwrap_or_else(|error| panic!("could not parse {}: {error}", path.display()));

        for mut case in spec_file.cases {
            case.file = file_name.clone();
            all_cases.push(case);
        }
    }

    all_cases
}

/// Loads one of the hand-checked outputs in `spec/examples/`, e.g. `"receipt.json"`.
pub fn load_example(name: &str) -> Value {
    let path = spec_dir().join("examples").join(name);
    let text = read_file(&path);
    serde_json::from_str(&text)
        .unwrap_or_else(|error| panic!("could not parse {}: {error}", path.display()))
}

/// All `.json` files in `spec/test-cases/`, sorted so the order is always the same.
fn test_case_files() -> Vec<PathBuf> {
    let dir = spec_dir().join("test-cases");
    let entries = fs::read_dir(&dir)
        .unwrap_or_else(|error| panic!("could not list {}: {error}", dir.display()));

    let mut files: Vec<PathBuf> = entries
        .map(|entry| entry.unwrap().path())
        .filter(|path| {
            path.extension()
                .is_some_and(|extension| extension == "json")
        })
        .collect();
    files.sort();
    files
}

fn read_file(path: &Path) -> String {
    fs::read_to_string(path)
        .unwrap_or_else(|error| panic!("could not read {}: {error}", path.display()))
}
