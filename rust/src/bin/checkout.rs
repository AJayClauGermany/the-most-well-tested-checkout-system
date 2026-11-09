//! The command line program (rules §8): basket JSON in on stdin, receipt JSON out on stdout.
//!
//! ```text
//! echo '{"items":[{"sku":"APPLE","quantity":1}]}' | cargo run --quiet --bin checkout
//! ```
//!
//! All the logic is in the library's `command` module (unit tested in-process); this file
//! only connects it to the terminal. It is tested from the outside in `tests/cli.rs`.

use std::process::ExitCode;

fn main() -> ExitCode {
    todo!("the checkout CLI is not implemented yet")
}
