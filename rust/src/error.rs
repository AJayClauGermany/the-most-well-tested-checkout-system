//! Error codes (rules §9) and the error type returned by `checkout`.

use serde::{Deserialize, Serialize};
use std::fmt;

/// The six ways a checkout can fail (rules §9).
///
/// `rename_all = "SCREAMING_SNAKE_CASE"` lets serde read `"EMPTY_BASKET"` from JSON
/// as `ErrorCode::EmptyBasket` (and write it back the same way).
#[derive(Debug, Clone, Copy, PartialEq, Eq, Hash, Serialize, Deserialize)]
#[serde(rename_all = "SCREAMING_SNAKE_CASE")]
pub enum ErrorCode {
    EmptyBasket,
    UnknownProduct,
    InvalidQuantity,
    UnknownShippingMethod,
    UnknownCoupon,
    CouponMinimumNotMet,
}

impl ErrorCode {
    /// Every error code, in the order the rules list them.
    pub const ALL: [ErrorCode; 6] = [
        ErrorCode::EmptyBasket,
        ErrorCode::UnknownProduct,
        ErrorCode::InvalidQuantity,
        ErrorCode::UnknownShippingMethod,
        ErrorCode::UnknownCoupon,
        ErrorCode::CouponMinimumNotMet,
    ];

    /// The code exactly as the spec spells it, e.g. `"EMPTY_BASKET"`.
    pub fn as_str(self) -> &'static str {
        match self {
            ErrorCode::EmptyBasket => "EMPTY_BASKET",
            ErrorCode::UnknownProduct => "UNKNOWN_PRODUCT",
            ErrorCode::InvalidQuantity => "INVALID_QUANTITY",
            ErrorCode::UnknownShippingMethod => "UNKNOWN_SHIPPING_METHOD",
            ErrorCode::UnknownCoupon => "UNKNOWN_COUPON",
            ErrorCode::CouponMinimumNotMet => "COUPON_MINIMUM_NOT_MET",
        }
    }
}

/// `Display` decides what `println!("{}", code)` and `code.to_string()` show.
impl fmt::Display for ErrorCode {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        f.write_str(self.as_str())
    }
}

/// The error returned when a checkout fails. It carries exactly one [`ErrorCode`]
/// (what programs look at) and a human-readable `message` (what a person reads in a log).
///
/// In Rust, failures are normal return values (`Result<Receipt, CheckoutError>`),
/// not exceptions.
#[derive(Debug, Clone, PartialEq, Eq)]
pub struct CheckoutError {
    pub code: ErrorCode,
    /// Extra detail such as `"Unknown product: BANANA"`. Empty when there is none.
    pub message: String,
}

impl CheckoutError {
    /// An error with a code and no further detail.
    pub fn new(code: ErrorCode) -> Self {
        CheckoutError {
            code,
            message: String::new(),
        }
    }

    /// An error with a code and a message explaining it.
    pub fn with_message(code: ErrorCode, message: impl Into<String>) -> Self {
        CheckoutError {
            code,
            message: message.into(),
        }
    }
}

/// Lets you write `ErrorCode::EmptyBasket.into()` (or use `?`) to get a `CheckoutError`.
impl From<ErrorCode> for CheckoutError {
    fn from(code: ErrorCode) -> Self {
        CheckoutError::new(code)
    }
}

impl fmt::Display for CheckoutError {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        write!(f, "checkout failed: {}", self.code)?;
        if !self.message.is_empty() {
            write!(f, ": {}", self.message)?;
        }
        Ok(())
    }
}

/// Makes `CheckoutError` a "real" error that works with `Box<dyn Error>` and friends.
impl std::error::Error for CheckoutError {}

// ---------------------------------------------------------------------------
// Unit tests
//
// `#[cfg(test)]` means "only compile this module when running `cargo test`".
// Putting unit tests in the same file as the code is the normal Rust convention.
// ---------------------------------------------------------------------------
#[cfg(test)]
mod tests {
    // `super::*` imports everything from the file above (including private items).
    use super::*;

    // A test is just a function marked with `#[test]`.
    // Its name says what behaviour we expect, so a failure message reads like a sentence.
    #[test]
    fn error_codes_are_spelled_exactly_like_the_spec() {
        // Arrange: the spellings copied from spec/checkout-rules.md §9, in the same order.
        let spec_spellings = [
            "EMPTY_BASKET",
            "UNKNOWN_PRODUCT",
            "INVALID_QUANTITY",
            "UNKNOWN_SHIPPING_METHOD",
            "UNKNOWN_COUPON",
            "COUPON_MINIMUM_NOT_MET",
        ];

        // Act: turn every code into text.
        let our_spellings = ErrorCode::ALL.map(|code| code.to_string());

        // Assert: `assert_eq!(actual, expected)` fails the test and prints both sides
        // when they differ.
        assert_eq!(our_spellings, spec_spellings);
    }

    #[test]
    fn error_codes_can_be_read_from_json() {
        // The shared test cases store errors as JSON strings, e.g. "expected_error": "UNKNOWN_COUPON".
        let code: ErrorCode = serde_json::from_str("\"UNKNOWN_COUPON\"").unwrap();

        assert_eq!(code, ErrorCode::UnknownCoupon);
    }

    #[test]
    fn error_codes_are_written_to_json_exactly_like_the_spec() {
        // The CLI prints {"error":"EMPTY_BASKET"}, so the JSON spelling matters too.
        let json = serde_json::to_string(&ErrorCode::CouponMinimumNotMet).unwrap();

        assert_eq!(json, "\"COUPON_MINIMUM_NOT_MET\"");
    }

    #[test]
    fn checkout_error_message_mentions_the_code() {
        let error = CheckoutError::from(ErrorCode::EmptyBasket);

        assert_eq!(error.code, ErrorCode::EmptyBasket);
        assert_eq!(error.message, "");
        assert_eq!(error.to_string(), "checkout failed: EMPTY_BASKET");
    }

    #[test]
    fn checkout_error_with_a_message_shows_both_the_code_and_the_message() {
        let error =
            CheckoutError::with_message(ErrorCode::UnknownProduct, "Unknown product: BANANA");

        assert_eq!(error.code, ErrorCode::UnknownProduct);
        assert_eq!(error.message, "Unknown product: BANANA");
        assert_eq!(
            error.to_string(),
            "checkout failed: UNKNOWN_PRODUCT: Unknown product: BANANA"
        );
    }

    #[test]
    fn checkout_error_is_a_normal_error_so_it_works_with_box_dyn_error() {
        // `?` in a function returning `Result<_, Box<dyn Error>>` relies on this.
        let boxed: Box<dyn std::error::Error> =
            Box::new(CheckoutError::new(ErrorCode::UnknownCoupon));

        assert_eq!(boxed.to_string(), "checkout failed: UNKNOWN_COUPON");
    }
}
