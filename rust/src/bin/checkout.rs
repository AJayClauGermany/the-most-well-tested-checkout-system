//! The command line program (rules §8): basket JSON in on stdin, receipt JSON out on stdout.
//!
//! ```text
//! echo '{"items":[{"sku":"APPLE","quantity":1}]}' | cargo run --quiet --bin checkout
//! ```
//!
//! All the logic is in the library's `command` module (unit tested in-process); this file
//! only connects it to the terminal. It is tested from the outside in `tests/cli.rs`, and it
//! is left out of the coverage and mutation reports, which only see code running inside
//! the test process.

use checkout::{CommandResult, load_standard_catalog, run_checkout_command};
use std::io::Read;
use std::process::ExitCode;

fn main() -> ExitCode {
    match run() {
        // A receipt (exit code 0) or {"error":"<CODE>"} (exit code 1), on one line.
        Ok(result) => {
            println!("{}", result.output);
            ExitCode::from(result.exit_code)
        }
        // Not a basket at all, or no catalog: nothing on stdout, a message on stderr.
        Err(message) => {
            eprintln!("checkout: {message}");
            ExitCode::from(2)
        }
    }
}

fn run() -> Result<CommandResult, String> {
    let mut input = String::new();
    std::io::stdin()
        .read_to_string(&mut input)
        .map_err(|error| format!("could not read standard input: {error}"))?;
    let catalog = load_standard_catalog().map_err(|error| error.to_string())?;
    run_checkout_command(&input, &catalog).map_err(|error| error.to_string())
}
