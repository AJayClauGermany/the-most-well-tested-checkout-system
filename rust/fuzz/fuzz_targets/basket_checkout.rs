//! Fuzz target for property P10 (spec/properties.md): for ANY input, the checkout either
//! gives a receipt, an error code, or says "not a basket" – it never crashes.
//!
//! libFuzzer feeds this function random bytes, millions of times, mutating the ones that
//! reach new code. Any panic (an overflow, an index out of range, an `unwrap` on `None`)
//! is reported with the exact input that caused it, saved in `fuzz/artifacts/`.
//!
//! Run with a nightly toolchain, from the `rust/` folder:
//!
//! ```text
//! cargo +nightly fuzz run basket_checkout -- -max_total_time=60
//! ```

#![no_main]

use checkout::{Catalog, load_standard_catalog, run_checkout_command};
use libfuzzer_sys::fuzz_target;
use std::sync::LazyLock;

/// The standard catalog, loaded once (the fuzzer calls the target millions of times).
static CATALOG: LazyLock<Catalog> =
    LazyLock::new(|| load_standard_catalog().expect("spec/catalog.json should be found"));

fuzz_target!(|data: &[u8]| {
    // Only valid UTF-8 can be a JSON text; the CLI reads text too.
    if let Ok(text) = std::str::from_utf8(data) {
        // Ok(receipt or error code) and Err(not a basket) are both fine. A panic is a bug.
        let _ = run_checkout_command(text, &CATALOG);
    }
});
