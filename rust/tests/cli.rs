//! End-to-end tests: run the real command line program (rules §8) as a separate process,
//! exactly as a user (or the shared runner in `e2e/run.sh`) would, and look only at what
//! comes out: stdout, stderr and the exit code.
//!
//! `assert_cmd` finds the binary Cargo built for this test run (`cargo_bin!("checkout")`),
//! feeds it stdin and gives us readable assertions; `predicates` describes what the output
//! should look like. Because the test only sees the outside of the program, it would pass
//! unchanged for a checkout written in any language.
//!
//! Run only this file with: `cargo test --test cli`

mod common;

use assert_cmd::Command;
use assert_cmd::cargo::cargo_bin;
use common::schemas::schema_errors;
use predicates::prelude::*;
use pretty_assertions::assert_eq;
use serde_json::Value;

/// The real program, ready to receive `input` on standard input.
fn checkout_cli(input: &str) -> Command {
    let mut command = Command::new(cargo_bin!("checkout"));
    command.write_stdin(input);
    command
}

#[test]
fn every_shared_spec_case_passes_through_the_real_program() {
    let mut failures = Vec::new();

    for case in common::load_spec_cases() {
        let input = format!("{}\n", serde_json::to_string(&case.basket).unwrap());
        let output = checkout_cli(&input).output().unwrap();
        let stdout = String::from_utf8(output.stdout).unwrap();
        let printed: Result<Value, _> = serde_json::from_str(&stdout);

        let problem = match (case.expected_error, &case.expected, printed) {
            (Some(code), None, Ok(printed)) => {
                let expected = serde_json::json!({ "error": code });
                if output.status.code() == Some(1) && printed == expected {
                    None
                } else {
                    Some(format!(
                        "expected {expected} with exit 1, got {printed} with {}",
                        output.status
                    ))
                }
            }
            (None, Some(totals), Ok(printed)) => {
                let expected = serde_json::to_value(totals).unwrap();
                let got: Value = expected
                    .as_object()
                    .unwrap()
                    .keys()
                    .map(|key| (key.clone(), printed[key].clone()))
                    .collect::<serde_json::Map<_, _>>()
                    .into();
                let schema_problems = schema_errors("receipt.schema.json", &printed);
                if output.status.code() == Some(0) && got == expected && schema_problems.is_empty()
                {
                    None
                } else {
                    Some(format!(
                        "expected {expected} with exit 0, got {got} with {}; schema: {schema_problems:?}",
                        output.status
                    ))
                }
            }
            (_, _, Err(error)) => Some(format!("no JSON on stdout ({error}): {stdout:?}")),
            _ => Some("broken case".to_string()),
        };

        if let Some(problem) = problem {
            failures.push(format!("[{}] {problem}", case.id));
        }
    }

    assert_eq!(failures, Vec::<String>::new());
}

#[test]
fn prints_exactly_one_line_of_json_and_exits_with_0() {
    checkout_cli(r#"{"items":[{"sku":"APPLE","quantity":1}]}"#)
        .assert()
        .success()
        .code(0)
        .stdout(predicate::str::ends_with("\n"))
        .stdout(predicate::function(|out: &str| {
            out.trim_end().lines().count() == 1
        }))
        .stdout(predicate::str::starts_with(r#"{"lines":[{"sku":"APPLE""#))
        .stderr(predicate::str::is_empty());
}

#[test]
fn prints_the_error_object_and_exits_with_1_for_an_invalid_basket() {
    checkout_cli(r#"{"items":[]}"#)
        .assert()
        .failure()
        .code(1)
        .stdout("{\"error\":\"EMPTY_BASKET\"}\n")
        .stderr(predicate::str::is_empty());
}

#[test]
fn exits_with_2_and_explains_on_stderr_when_the_input_is_not_json() {
    checkout_cli("this is not json")
        .assert()
        .code(2)
        .stdout(predicate::str::is_empty())
        .stderr(predicate::str::starts_with("checkout: "));
}

#[test]
fn exits_with_2_when_the_json_is_not_a_basket() {
    checkout_cli(r#"{"hello": "world"}"#)
        .assert()
        .code(2)
        .stdout(predicate::str::is_empty())
        .stderr(predicate::str::contains("not a basket"));
}

#[test]
fn exits_with_2_on_empty_input() {
    checkout_cli("")
        .assert()
        .code(2)
        .stdout(predicate::str::is_empty());
}

#[test]
fn finds_the_standard_catalog_from_its_own_location_whatever_the_current_directory() {
    // The shared runner starts the program from the repository root, a user might start it
    // from anywhere: the catalog is found by walking up from the executable, not from here.
    let elsewhere = tempfile::tempdir().unwrap();

    checkout_cli(r#"{"items":[{"sku":"COFFEE","quantity":3}]}"#)
        .current_dir(elsewhere.path())
        .assert()
        .success()
        .stdout(predicate::str::contains(r#""promotion_discount":899"#));
}
