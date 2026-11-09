//! Test helper: JSON Schema validation with the `jsonschema` crate.
//!
//! The schemas in `spec/schema/` are the *contract* every language's output must honour.
//! A "contract test" checks a real receipt (or saved order, or CLI error) against them.
//! The schemas use JSON Schema draft 2020-12.

use jsonschema::{Draft, Validator};
use serde_json::Value;
use std::fs;

use super::spec_dir;

/// Validates `value` against one of the shared schemas, e.g. `"receipt.schema.json"`.
///
/// Returns an empty list when the value is valid, otherwise one message per problem
/// (an empty list is easy to assert with `assert_eq!(errors, [])` and prints nicely when
/// the assertion fails). Messages look like `"/tax -1 is less than the minimum of 0"`.
pub fn schema_errors(schema_name: &str, value: &Value) -> Vec<String> {
    let validator = validator_for(schema_name);
    validator
        .iter_errors(value)
        .map(|error| {
            let path = error.instance_path().to_string();
            let path = if path.is_empty() {
                "/".to_string()
            } else {
                path
            };
            format!("{path} {error}")
        })
        .collect()
}

/// Compiles one schema from `spec/schema/`.
///
/// `order.schema.json` refers to `receipt.schema.json` with a *relative* `$ref`. A relative
/// reference is resolved against a base URI, so we tell the validator that the base is the
/// `spec/schema/` folder (as a `file://` URL); the crate's file resolver then reads the
/// referenced schema from there.
fn validator_for(schema_name: &str) -> Validator {
    let folder = spec_dir()
        .join("schema")
        .canonicalize()
        .expect("spec/schema should exist");
    let path = folder.join(schema_name);
    let text = fs::read_to_string(&path)
        .unwrap_or_else(|error| panic!("could not read {}: {error}", path.display()));
    let schema: Value = serde_json::from_str(&text)
        .unwrap_or_else(|error| panic!("could not parse {}: {error}", path.display()));

    jsonschema::options()
        .with_draft(Draft::Draft202012)
        .with_base_uri(format!("file://{}/", folder.display()))
        .build(&schema)
        .unwrap_or_else(|error| panic!("could not compile {schema_name}: {error}"))
}
