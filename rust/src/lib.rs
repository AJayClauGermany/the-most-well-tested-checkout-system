//! A small checkout system, used to show off many kinds of software testing.
//!
//! The rules live in `spec/checkout-rules.md` at the root of the repository.
//! This crate only contains the *types* and an empty [`checkout`] function so far:
//! we write the tests first (TDD) and the implementation later.
//!
//! Where do the tests live? (Rust has its own conventions!)
//!
//! * **Unit tests** sit at the bottom of each `src/*.rs` file, inside a
//!   `#[cfg(test)] mod tests { ... }` block. They can even test private functions.
//! * **Integration tests** live in the `tests/` folder. Each file there is compiled as a
//!   separate program that can only use the *public* API, exactly like a real user would.
//! * **Doc tests** are the code examples inside `///` comments, like the one on [`checkout`].

// Each `mod` line pulls in the file with the same name from `src/`.
mod basket;
mod catalog;
mod error;
mod receipt;

// Re-export the types so users can write `checkout::Basket` instead of `checkout::basket::Basket`.
pub use basket::{Basket, BasketItem};
pub use catalog::{Catalog, Product, Promotion};
pub use error::{CheckoutError, ErrorCode};
pub use receipt::{Receipt, ReceiptLine};

/// An amount of money in whole cents: `199` means $1.99.
///
/// We never use floating point numbers for money (see rules §1).
/// `i64` (a signed number) keeps subtraction simple and can never overflow for our baskets.
pub type Cents = i64;

/// Calculates the receipt for a basket, using the products in `catalog`.
///
/// Returns `Ok(receipt)` on success, or `Err(error)` carrying one of the six
/// [`ErrorCode`]s when the basket is invalid (rules §4.1).
///
/// # Example
///
/// `no_run` means: this example must compile, but is not run (yet), because
/// `checkout` is not implemented. Change it to a normal example once it works!
///
/// ```no_run
/// use checkout::{Basket, BasketItem, Catalog, Product, checkout};
///
/// let catalog = Catalog::new(vec![Product {
///     sku: "APPLE".to_string(),
///     name: "Apple".to_string(),
///     unit_price: 50,
///     promotion: None,
/// }]);
/// let basket = Basket {
///     items: vec![BasketItem { sku: "APPLE".to_string(), quantity: 2 }],
///     coupon: None,
///     shipping: "standard".to_string(),
/// };
///
/// let receipt = checkout(&basket, &catalog).unwrap();
/// assert_eq!(receipt.subtotal, 100);
/// ```
pub fn checkout(basket: &Basket, catalog: &Catalog) -> Result<Receipt, CheckoutError> {
    // `let _ = ...` tells the compiler "yes, we know these are unused for now".
    let _ = (basket, catalog);

    // `todo!()` compiles fine, but panics (crashes) when it is reached.
    // Remove it when you start implementing – and un-ignore the test in tests/spec_cases.rs.
    todo!("checkout() is not implemented yet")
}
