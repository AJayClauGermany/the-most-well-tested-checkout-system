//! The command line program's logic (rules §8), kept free of stdin/stdout so it can be
//! unit tested in-process. `src/bin/checkout.rs` is the thin wrapper that talks to the
//! terminal; it is tested from the outside in `tests/cli.rs`.

use crate::{Basket, Catalog, checkout};
use std::fmt;
use std::path::PathBuf;

/// What the program prints and the exit code it ends with: `0` with a receipt,
/// `1` with `{"error":"<CODE>"}`.
#[derive(Debug, Clone, PartialEq, Eq)]
pub struct CommandResult {
    pub output: String,
    pub exit_code: u8,
}

/// The program could not run at all: the input was not a basket, or the standard
/// catalog could not be found or read. The CLI reports these on stderr with exit code 2.
#[derive(Debug, Clone, PartialEq, Eq)]
pub struct CommandError {
    pub message: String,
}

impl CommandError {
    pub fn new(message: impl Into<String>) -> Self {
        CommandError {
            message: message.into(),
        }
    }
}

impl fmt::Display for CommandError {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        f.write_str(&self.message)
    }
}

impl std::error::Error for CommandError {}

/// Where the standard catalog lives: `spec/catalog.json` at the root of the repository.
pub fn standard_catalog_path() -> Result<PathBuf, CommandError> {
    todo!("standard_catalog_path() is not implemented yet")
}

/// Reads `spec/catalog.json`.
pub fn load_standard_catalog() -> Result<Catalog, CommandError> {
    todo!("load_standard_catalog() is not implemented yet")
}

/// Runs one checkout from the JSON text of a basket.
///
/// * a valid basket gives the receipt as one line of JSON and exit code 0,
/// * an invalid basket gives `{"error":"<CODE>"}` and exit code 1,
/// * text that is not a basket at all is an `Err` (the CLI turns that into exit code 2).
///
/// ```
/// use checkout::{Catalog, Product, run_checkout_command};
///
/// let catalog = Catalog::new(vec![Product {
///     sku: "APPLE".to_string(), name: "Apple".to_string(), unit_price: 50, promotion: None,
/// }]);
///
/// let result = run_checkout_command(r#"{"items":[]}"#, &catalog).unwrap();
/// assert_eq!(result.exit_code, 1);
/// assert_eq!(result.output, r#"{"error":"EMPTY_BASKET"}"#);
/// ```
pub fn run_checkout_command(input: &str, catalog: &Catalog) -> Result<CommandResult, CommandError> {
    let _ = (input, catalog, checkout);
    todo!("run_checkout_command() is not implemented yet")
}

#[allow(dead_code)]
fn unused_until_implemented(basket: Basket) -> Basket {
    basket
}

// ---------------------------------------------------------------------------
// Unit tests, run in-process. The real program is tested from the outside in tests/cli.rs.
// ---------------------------------------------------------------------------
#[cfg(test)]
mod tests {
    use super::*;
    use crate::test_fixtures::{one_apple_receipt, tiny_catalog};
    use pretty_assertions::assert_eq;
    use rstest::rstest;
    use serde_json::{Value, json};

    const ONE_APPLE: &str = r#"{"items":[{"sku":"APPLE","quantity":1}]}"#;

    #[test]
    fn the_standard_catalog_is_spec_catalog_json_at_the_root_of_the_repository() {
        let path = standard_catalog_path().unwrap();

        assert!(path.ends_with(PathBuf::from("spec").join("catalog.json")));
        // The repository root is the folder that also holds the rules.
        assert!(path.with_file_name("checkout-rules.md").is_file());
    }

    #[test]
    fn loads_the_six_standard_products() {
        let catalog = load_standard_catalog().unwrap();

        let skus: Vec<&str> = catalog.products.iter().map(|p| p.sku.as_str()).collect();
        assert_eq!(
            skus,
            ["APPLE", "BREAD", "COFFEE", "NOTEBOOK", "PEN", "HEADPHONES"]
        );
    }

    #[test]
    fn prints_the_receipt_as_json_and_exits_with_0_for_a_valid_basket() {
        let result = run_checkout_command(ONE_APPLE, &tiny_catalog()).unwrap();

        assert_eq!(result.exit_code, 0);
        let printed: Value = serde_json::from_str(&result.output).unwrap();
        assert_eq!(printed, serde_json::to_value(one_apple_receipt()).unwrap());
    }

    #[test]
    fn prints_the_receipt_fields_in_the_order_of_rules_5() {
        // Field order is not something JSON cares about, but humans reading the output do.
        let result = run_checkout_command(ONE_APPLE, &tiny_catalog()).unwrap();

        assert_eq!(
            result.output,
            r#"{"lines":[{"sku":"APPLE","name":"Apple","quantity":1,"unit_price":50,"line_total":50,"line_discount":0}],"subtotal":50,"promotion_discount":0,"coupon_discount":0,"shipping":499,"tax":4,"total":553}"#
        );
    }

    #[test]
    fn prints_an_error_object_and_exits_with_1_for_an_invalid_basket() {
        let result = run_checkout_command(r#"{"items":[]}"#, &tiny_catalog()).unwrap();

        assert_eq!(
            result,
            CommandResult {
                output: r#"{"error":"EMPTY_BASKET"}"#.to_string(),
                exit_code: 1,
            }
        );
    }

    #[test]
    fn reports_every_kind_of_error_the_same_way() {
        let input = r#"{"items":[{"sku":"APPLE","quantity":1}],"coupon":"NOPE"}"#;

        let result = run_checkout_command(input, &tiny_catalog()).unwrap();

        assert_eq!(result.exit_code, 1);
        assert_eq!(
            serde_json::from_str::<Value>(&result.output).unwrap(),
            json!({ "error": "UNKNOWN_COUPON" })
        );
    }

    #[test]
    fn prints_compact_json_on_a_single_line() {
        let result = run_checkout_command(ONE_APPLE, &tiny_catalog()).unwrap();

        assert!(!result.output.contains('\n'));
        assert!(!result.output.contains(": "));
    }

    #[test]
    fn accepts_surrounding_whitespace_as_a_shell_pipeline_adds_a_newline() {
        let input = format!("  {ONE_APPLE}\n");

        let result = run_checkout_command(&input, &tiny_catalog()).unwrap();

        assert_eq!(result.exit_code, 0);
    }

    #[test]
    fn missing_coupon_and_shipping_mean_no_coupon_and_standard_shipping() {
        let result = run_checkout_command(ONE_APPLE, &tiny_catalog()).unwrap();
        let printed: Value = serde_json::from_str(&result.output).unwrap();

        assert_eq!(printed["coupon_discount"], 0);
        assert_eq!(printed["shipping"], 499);
    }

    #[test]
    fn fails_rather_than_inventing_an_error_code_when_the_input_is_not_json() {
        let error = run_checkout_command("not json", &tiny_catalog()).unwrap_err();
        assert!(
            error.message.starts_with("The input is not a basket: "),
            "unexpected message: {}",
            error.message
        );

        let error = run_checkout_command("", &tiny_catalog()).unwrap_err();
        assert!(error.message.starts_with("The input is not a basket: "));
    }

    // The JSON is fine, but it is not a basket. The reason comes from serde, so it is
    // precise: it names the field or the wrong type.
    #[rstest]
    #[case("null", "null, expected struct Basket")]
    #[case("42", "integer `42`, expected struct Basket")]
    #[case(r#""text""#, "string \"text\", expected struct Basket")]
    #[case("[]", "sequence, expected struct Basket")]
    #[case("{}", "missing field `items`")]
    #[case(r#"{"items": "APPLE"}"#, "string \"APPLE\", expected a sequence")]
    #[case(r#"{"items": [1, 2]}"#, "integer `1`, expected struct BasketItem")]
    #[case(r#"{"items": [null]}"#, "null, expected struct BasketItem")]
    #[case(
        r#"{"items": ["APPLE"]}"#,
        "string \"APPLE\", expected struct BasketItem"
    )]
    #[case(r#"{"items": [{"sku": 5}]}"#, "integer `5`, expected a string")]
    #[case(r#"{"items": [{"sku": "APPLE"}]}"#, "missing field `quantity`")]
    #[case(
        r#"{"items": [{"sku": "APPLE", "quantity": 1.5}]}"#,
        "floating point `1.5`, expected i64"
    )]
    #[case(
        r#"{"items": [{"sku": "APPLE", "quantity": 1}], "shipping": null}"#,
        "null, expected a string"
    )]
    fn fails_when_the_json_is_not_a_basket(#[case] input: &str, #[case] reason: &str) {
        let error = run_checkout_command(input, &tiny_catalog()).unwrap_err();

        assert!(
            error.message.starts_with("The input is not a basket: "),
            "{}",
            error.message
        );
        assert!(
            error.message.contains(reason),
            "{} should mention {reason:?}",
            error.message
        );
    }

    #[test]
    fn command_errors_display_their_message() {
        let error = CommandError::new("something went wrong");

        assert_eq!(error.to_string(), "something went wrong");
    }
}
