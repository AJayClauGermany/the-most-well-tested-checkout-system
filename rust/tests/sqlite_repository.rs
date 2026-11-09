//! Integration tests for the SQLite order repository (the `rusqlite` crate, SQLite bundled).
//!
//! A real database, but no server to install: SQLite is one file. Each test gets its own
//! file in a temporary folder from `tempfile::tempdir()`.
//!
//! Run only this file with: `cargo test --test sqlite_repository`

mod common;

use checkout::{Order, OrderRepository, SqliteOrderRepository, checkout, parse_timestamp};
use common::fixtures::{SharedClock, basket_with, one_apple_receipt, tiny_catalog};
use common::schemas::schema_errors;
use common::{FIXED_TIME, fixed_clock};
use pretty_assertions::assert_eq;
use rstest::rstest;
use rusqlite::Connection;
use std::path::{Path, PathBuf};
use tempfile::TempDir;

/// A fresh temporary folder and the database path inside it.
fn temporary_database() -> (TempDir, PathBuf) {
    let folder = tempfile::tempdir().expect("a temporary folder should be available");
    let path = folder.path().join("orders.db");
    (folder, path)
}

fn open(path: &Path) -> SqliteOrderRepository {
    SqliteOrderRepository::open(path, fixed_clock()).expect("the database should open")
}

#[test]
fn saves_and_finds_an_order() {
    let (_folder, path) = temporary_database();
    let mut repository = open(&path);
    let basket = basket_with(&[("COFFEE", 3)], Some("SAVE10"), "standard");
    let receipt = checkout(&basket, &tiny_catalog()).unwrap();

    let id = repository.save(&receipt).unwrap();

    assert_eq!(id, "ORDER-1");
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
fn gives_back_an_order_that_follows_order_schema_json() {
    let (_folder, path) = temporary_database();
    let mut repository = open(&path);
    let id = repository.save(&one_apple_receipt()).unwrap();

    let order = repository.find(&id).unwrap().unwrap();

    let as_json = serde_json::to_value(&order).unwrap();
    assert_eq!(
        schema_errors("order.schema.json", &as_json),
        Vec::<String>::new()
    );
}

#[test]
fn numbers_orders_sequentially() {
    let (_folder, path) = temporary_database();
    let mut repository = open(&path);

    let ids = [
        repository.save(&one_apple_receipt()).unwrap(),
        repository.save(&one_apple_receipt()).unwrap(),
        repository.save(&one_apple_receipt()).unwrap(),
    ];

    assert_eq!(ids, ["ORDER-1", "ORDER-2", "ORDER-3"]);
}

#[test]
fn continues_numbering_after_the_database_is_closed_and_opened_again() {
    let (_folder, path) = temporary_database();
    let mut first = open(&path);
    first.save(&one_apple_receipt()).unwrap();
    first.save(&one_apple_receipt()).unwrap();
    first.close().unwrap();

    let mut second = open(&path);

    assert_eq!(second.save(&one_apple_receipt()).unwrap(), "ORDER-3");
    assert_eq!(second.find("ORDER-2").unwrap().unwrap().id, "ORDER-2");
}

#[test]
fn really_writes_to_the_database_file_checked_with_a_separate_connection() {
    let (_folder, path) = temporary_database();
    let mut repository = open(&path);
    repository.save(&one_apple_receipt()).unwrap();
    repository.close().unwrap();

    let connection = Connection::open(&path).unwrap();
    let count: i64 = connection
        .query_row("SELECT COUNT(*) FROM orders", [], |row| row.get(0))
        .unwrap();

    assert_eq!(count, 1);
}

#[test]
fn works_in_memory_too_a_throw_away_database_without_a_file() {
    let mut repository = SqliteOrderRepository::in_memory(fixed_clock()).unwrap();

    assert_eq!(repository.save(&one_apple_receipt()).unwrap(), "ORDER-1");
    assert_eq!(
        repository.find("ORDER-1").unwrap().unwrap().receipt,
        one_apple_receipt()
    );
}

// "?" placeholders keep values apart from the SQL, so even an id that looks like an SQL
// command is just text that matches no order.
#[rstest]
#[case("ORDER-2")]
#[case("ORDER-0")]
#[case("ORDER-01")]
#[case("ORDER-x")]
#[case("1")]
#[case("")]
#[case("ORDER-1; DROP TABLE orders")]
#[case("1 OR 1=1")]
fn answers_not_found_for_ids_that_do_not_exist_without_failing(#[case] id: &str) {
    let (_folder, path) = temporary_database();
    let mut repository = open(&path);
    repository.save(&one_apple_receipt()).unwrap();

    assert_eq!(repository.find(id), Ok(None));
    assert!(repository.find("ORDER-1").unwrap().is_some()); // and the table is still there
}

#[test]
fn stamps_each_order_with_the_clocks_time_at_the_moment_of_saving() {
    let (_folder, path) = temporary_database();
    let clock = SharedClock::starting_at(parse_timestamp("2026-01-15T10:30:00Z").unwrap());
    let mut repository = SqliteOrderRepository::open(&path, Box::new(clock.clone())).unwrap();

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

// --- When the database misbehaves, the repository reports an error (it never panics) ---

#[test]
fn opening_fails_when_the_folder_for_the_file_does_not_exist() {
    let (folder, _) = temporary_database();
    let path = folder.path().join("no-such-folder").join("orders.db");

    let result = SqliteOrderRepository::open(&path, fixed_clock());

    assert!(
        result.is_err(),
        "SQLite cannot create a file in a missing folder"
    );
}

#[test]
fn finding_fails_with_a_repository_error_when_a_stored_row_is_corrupt() {
    let (_folder, path) = temporary_database();
    let repository = open(&path);
    let connection = Connection::open(&path).unwrap();
    connection
        .execute(
            "INSERT INTO orders (created_at, receipt) VALUES (?1, ?2)",
            ("2026-01-15T10:30:00Z", "this is not a receipt"),
        )
        .unwrap();

    let error = repository.find("ORDER-1").unwrap_err();

    assert!(error.message.contains("expected"), "{}", error.message);
}
