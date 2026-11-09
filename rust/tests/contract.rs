//! Contract tests: is the output even the right shape?
//!
//! `spec/schema/receipt.schema.json` is the contract every language's receipt must honour,
//! and `cli-error.schema.json` the contract for what the CLI prints on failure. Here the
//! `jsonschema` crate checks real output against them. The shared e2e runner only compares
//! totals; this is where the shape of the whole receipt (lines included) is held to the
//! contract.
//!
//! Run only this file with: `cargo test --test contract`

mod common;

use checkout::{InMemoryOrderRepository, OrderRepository, checkout, run_checkout_command};
use common::fixtures::{basket, one_apple_receipt, tiny_catalog};
use common::schemas::schema_errors;
use pretty_assertions::assert_eq;
use serde_json::{Value, json};

fn no_errors() -> Vec<String> {
    Vec::new()
}

fn to_json<T: serde::Serialize>(value: &T) -> Value {
    serde_json::to_value(value).unwrap()
}

#[test]
fn the_receipt_of_every_valid_shared_case_follows_receipt_schema_json() {
    let catalog = common::load_catalog();
    let mut failures = Vec::new();

    for case in common::load_spec_cases()
        .into_iter()
        .filter(|c| c.expected.is_some())
    {
        let receipt = checkout(&case.basket, &catalog).unwrap();
        let errors = schema_errors("receipt.schema.json", &to_json(&receipt));
        if !errors.is_empty() {
            failures.push(format!("[{}] {}", case.id, errors.join("; ")));
        }
    }

    assert_eq!(failures, Vec::<String>::new());
}

#[test]
fn a_receipt_from_a_made_up_catalog_follows_the_schema_too() {
    let receipt = checkout(&basket(&[("TV", 1), ("COFFEE", 6)]), &tiny_catalog()).unwrap();

    assert_eq!(
        schema_errors("receipt.schema.json", &to_json(&receipt)),
        no_errors()
    );
}

#[test]
fn the_hand_checked_example_receipt_follows_the_schema_the_spec_agrees_with_itself() {
    let example = common::load_example("receipt.json");

    assert_eq!(schema_errors("receipt.schema.json", &example), no_errors());
}

#[test]
fn the_schema_check_itself_works_a_broken_receipt_is_rejected_testing_the_test() {
    let mut broken = to_json(&one_apple_receipt());
    broken["shipping"] = json!(123);
    broken["tax"] = json!(-1);
    broken["extra"] = json!(true);

    let mut errors = schema_errors("receipt.schema.json", &broken);
    errors.sort();

    assert_eq!(
        errors,
        vec![
            "/ Additional properties are not allowed ('extra' was unexpected)".to_string(),
            "/shipping 123 is not one of 0, 499 or 999".to_string(),
            "/tax -1 is less than the minimum of 0".to_string(),
        ]
    );
}

#[test]
fn the_commands_output_for_every_failing_shared_case_follows_cli_error_schema_json() {
    let catalog = common::load_catalog();
    let mut failures = Vec::new();

    for case in common::load_spec_cases()
        .into_iter()
        .filter(|c| c.expected_error.is_some())
    {
        let input = serde_json::to_string(&case.basket).unwrap();
        let result = run_checkout_command(&input, &catalog).unwrap();
        let printed: Value = serde_json::from_str(&result.output).unwrap();
        let errors = schema_errors("cli-error.schema.json", &printed);
        if result.exit_code != 1 || !errors.is_empty() {
            failures.push(format!(
                "[{}] exit {} {}",
                case.id,
                result.exit_code,
                errors.join("; ")
            ));
        }
    }

    assert_eq!(failures, Vec::<String>::new());
}

#[test]
fn the_cli_error_schema_rejects_made_up_codes_testing_the_test() {
    let errors = schema_errors(
        "cli-error.schema.json",
        &json!({ "error": "OUT_OF_COFFEE" }),
    );

    assert_eq!(errors.len(), 1, "{errors:?}");
    assert!(errors[0].starts_with("/error "), "{errors:?}");
}

#[test]
fn every_shared_test_cases_basket_is_a_valid_basket_the_spec_is_consistent_with_itself() {
    for case in common::load_spec_cases() {
        let errors = schema_errors("basket.schema.json", &to_json(&case.basket));

        assert_eq!(errors, no_errors(), "case {}", case.id);
    }
}

#[test]
fn a_saved_order_follows_order_schema_json_which_refers_to_the_receipt_schema() {
    let mut repository = InMemoryOrderRepository::new(common::fixed_clock());
    let id = repository.save(&one_apple_receipt()).unwrap();
    let order = repository.find(&id).unwrap().unwrap();

    assert_eq!(
        schema_errors("order.schema.json", &to_json(&order)),
        no_errors()
    );
    // The $ref really is followed: a broken receipt inside the order is caught.
    let mut broken = to_json(&order);
    broken["receipt"]["total"] = json!(-5);
    assert_eq!(
        schema_errors("order.schema.json", &broken),
        vec!["/receipt/total -5 is less than the minimum of 0".to_string()]
    );
}
