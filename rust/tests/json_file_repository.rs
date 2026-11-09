//! Integration tests for the JSON-file order repository: real files in a real temporary folder.
//!
//! Unlike a unit test, this touches the file system, so each test gets its own fresh folder
//! from `tempfile::tempdir()`, which is deleted when the `TempDir` value is dropped at the
//! end of the test. The stored files must follow `spec/schema/order.schema.json` – a
//! contract test inside an integration test.
//!
//! Run only this file with: `cargo test --test json_file_repository`

mod common;

use checkout::{JsonFileOrderRepository, Order, OrderRepository, checkout, parse_timestamp};
use common::fixtures::{SharedClock, one_apple_receipt, tiny_catalog};
use common::schemas::schema_errors;
use common::{FIXED_TIME, fixed_clock};
use pretty_assertions::assert_eq;
use rstest::rstest;
use serde_json::{Value, json};
use std::fs;
use std::path::Path;
use tempfile::TempDir;

/// A fresh temporary folder for one test.
fn temporary_folder() -> TempDir {
    tempfile::tempdir().expect("a temporary folder should be available")
}

fn open(folder: &Path) -> JsonFileOrderRepository {
    JsonFileOrderRepository::open(folder, fixed_clock()).expect("the folder should open")
}

/// The JSON a stored order should contain.
fn expected_order_json(id: &str) -> Value {
    json!({
        "id": id,
        "created_at": FIXED_TIME,
        "receipt": one_apple_receipt(),
    })
}

#[test]
fn writes_order_1_json_into_the_folder_on_the_first_save() {
    let folder = temporary_folder();
    let mut repository = open(folder.path());

    let id = repository.save(&one_apple_receipt()).unwrap();

    assert_eq!(id, "ORDER-1");
    assert!(folder.path().join("ORDER-1.json").is_file());
}

#[test]
fn writes_a_file_that_follows_order_schema_json_exactly() {
    let folder = temporary_folder();
    let mut repository = open(folder.path());
    repository.save(&one_apple_receipt()).unwrap();

    let text = fs::read_to_string(folder.path().join("ORDER-1.json")).unwrap();
    let stored: Value = serde_json::from_str(&text).unwrap();

    assert_eq!(
        schema_errors("order.schema.json", &stored),
        Vec::<String>::new()
    );
    assert_eq!(stored, expected_order_json("ORDER-1"));
}

#[test]
fn writes_human_readable_json_indented_and_ending_in_a_newline() {
    let folder = temporary_folder();
    let mut repository = open(folder.path());
    repository.save(&one_apple_receipt()).unwrap();

    let text = fs::read_to_string(folder.path().join("ORDER-1.json")).unwrap();

    assert!(text.contains("\n  "), "should be indented:\n{text}");
    assert!(text.ends_with('\n'), "should end with a newline");
}

#[test]
fn reads_a_saved_order_back_exactly_as_saved() {
    let folder = temporary_folder();
    let mut repository = open(folder.path());
    let basket = common::fixtures::basket_with(&[("COFFEE", 3)], Some("SAVE10"), "standard");
    let receipt = checkout(&basket, &tiny_catalog()).unwrap();

    let id = repository.save(&receipt).unwrap();

    assert_eq!(
        repository.find(&id),
        Ok(Some(Order {
            id: "ORDER-1".to_string(),
            created_at: FIXED_TIME.to_string(),
            receipt,
        }))
    );
}

#[test]
fn numbers_orders_sequentially() {
    let folder = temporary_folder();
    let mut repository = open(folder.path());

    let ids = [
        repository.save(&one_apple_receipt()).unwrap(),
        repository.save(&one_apple_receipt()).unwrap(),
        repository.save(&one_apple_receipt()).unwrap(),
    ];

    assert_eq!(ids, ["ORDER-1", "ORDER-2", "ORDER-3"]);
    let mut files: Vec<String> = fs::read_dir(folder.path())
        .unwrap()
        .map(|entry| entry.unwrap().file_name().to_string_lossy().to_string())
        .collect();
    files.sort();
    assert_eq!(files, ["ORDER-1.json", "ORDER-2.json", "ORDER-3.json"]);
}

#[test]
fn continues_numbering_when_the_same_folder_is_opened_again_later() {
    let folder = temporary_folder();
    let mut first = open(folder.path());
    first.save(&one_apple_receipt()).unwrap();
    first.save(&one_apple_receipt()).unwrap();

    let mut second = open(folder.path());

    assert_eq!(second.save(&one_apple_receipt()).unwrap(), "ORDER-3");
    assert_eq!(second.find("ORDER-1").unwrap().unwrap().id, "ORDER-1");
}

#[test]
fn numbers_after_the_highest_existing_order_even_with_gaps_and_two_digit_numbers() {
    let folder = temporary_folder();
    fs::write(
        folder.path().join("ORDER-3.json"),
        expected_order_json("ORDER-3").to_string(),
    )
    .unwrap();
    fs::write(
        folder.path().join("ORDER-12.json"),
        expected_order_json("ORDER-12").to_string(),
    )
    .unwrap();
    let mut repository = open(folder.path());

    // Not ORDER-4: alphabetically "3" > "12", but numerically 12 is the highest.
    assert_eq!(repository.save(&one_apple_receipt()).unwrap(), "ORDER-13");
}

#[test]
fn ignores_files_in_the_folder_that_are_not_orders() {
    let folder = temporary_folder();
    fs::write(folder.path().join("README.txt"), "not an order").unwrap();
    fs::write(
        folder.path().join("ORDER-9.json.bak"),
        "not an order either",
    )
    .unwrap();
    fs::write(
        folder.path().join("ORDER-9"),
        "no .json extension, so not an order",
    )
    .unwrap();
    fs::write(
        folder.path().join("ORDER-09.json"),
        "a leading zero is not an order id",
    )
    .unwrap();
    fs::create_dir(folder.path().join("ORDER-7.json")).unwrap(); // a folder, not a file
    let mut repository = open(folder.path());

    assert_eq!(repository.save(&one_apple_receipt()).unwrap(), "ORDER-1");
}

#[test]
fn creates_the_folder_if_it_does_not_exist_yet() {
    let folder = temporary_folder();
    let nested = folder.path().join("orders").join("2026");
    let mut repository = open(&nested);

    repository.save(&one_apple_receipt()).unwrap();

    assert!(nested.join("ORDER-1.json").is_file());
    assert_eq!(repository.folder(), nested.as_path());
}

#[rstest]
#[case("ORDER-2")]
#[case("ORDER-0")]
#[case("nonsense")]
#[case("")]
fn answers_not_found_for_ids_that_have_no_file_without_failing(#[case] id: &str) {
    let folder = temporary_folder();
    let mut repository = open(folder.path());
    repository.save(&one_apple_receipt()).unwrap();

    assert_eq!(repository.find(id), Ok(None));
}

#[test]
fn does_not_let_an_id_escape_the_folder_no_path_tricks() {
    let folder = temporary_folder();
    fs::write(
        folder.path().join("secret.json"),
        r#"{"id":"x","created_at":"x","receipt":{}}"#,
    )
    .unwrap();
    let repository = open(&folder.path().join("orders"));

    assert_eq!(repository.find("../secret"), Ok(None));
    assert_eq!(repository.find("../secret.json"), Ok(None));
}

#[test]
fn stamps_each_order_with_the_clocks_time_at_the_moment_of_saving() {
    let folder = temporary_folder();
    let clock = SharedClock::starting_at(parse_timestamp("2026-01-15T10:30:00Z").unwrap());
    let mut repository =
        JsonFileOrderRepository::open(folder.path(), Box::new(clock.clone())).unwrap();

    let first = repository.save(&one_apple_receipt()).unwrap();
    clock.set(parse_timestamp("2026-01-16T09:00:00Z").unwrap());
    let second = repository.save(&one_apple_receipt()).unwrap();

    assert_eq!(
        repository.find(&first).unwrap().unwrap().created_at,
        "2026-01-15T10:30:00Z"
    );
    assert_eq!(
        repository.find(&second).unwrap().unwrap().created_at,
        "2026-01-16T09:00:00Z"
    );
}

// --- When the file system misbehaves, the repository reports an error (it never panics) --

#[test]
fn opening_fails_when_the_path_is_a_file_not_a_folder() {
    let folder = temporary_folder();
    let file = folder.path().join("a-file");
    fs::write(&file, "I am a file").unwrap();

    let result = JsonFileOrderRepository::open(&file, fixed_clock());

    assert!(result.is_err(), "opening a file as a folder should fail");
}

#[test]
fn saving_fails_with_a_repository_error_when_the_folder_has_disappeared() {
    let folder = temporary_folder();
    let orders = folder.path().join("orders");
    let mut repository = open(&orders);
    fs::remove_dir_all(&orders).unwrap();

    let error = repository.save(&one_apple_receipt()).unwrap_err();

    assert!(!error.message.is_empty());
}

#[test]
fn finding_fails_with_a_repository_error_when_a_stored_file_is_corrupt() {
    let folder = temporary_folder();
    let repository = open(folder.path());
    fs::write(folder.path().join("ORDER-1.json"), "this is not json").unwrap();

    let error = repository.find("ORDER-1").unwrap_err();

    assert!(error.message.contains("expected"), "{}", error.message);
}
