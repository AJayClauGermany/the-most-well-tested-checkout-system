//! Products, promotions and the catalog (rules §2).

use crate::{Cents, CheckoutError};
use serde::{Deserialize, Serialize};

/// A special offer attached to a product (rules §4.2).
///
/// In JSON the kind of promotion is stored in a `"type"` field:
///
/// ```json
/// { "type": "buy_x_get_y_free", "buy": 2, "free": 1 }
/// { "type": "bulk_discount", "min_quantity": 10, "percent": 10 }
/// ```
///
/// `#[serde(tag = "type")]` tells serde to use that field to pick the enum variant.
/// A Rust `enum` is perfect here: a promotion is *either* one kind *or* the other,
/// and each kind has its own fields.
///
/// The numbers are `i64`, the same type as a basket quantity, so they can be
/// combined in calculations without converting between number types.
#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[serde(tag = "type", rename_all = "snake_case")]
pub enum Promotion {
    /// For every complete group of `buy + free` items, `free` of them cost nothing.
    BuyXGetYFree { buy: i64, free: i64 },
    /// `percent`% off the line total when at least `min_quantity` items are bought.
    BulkDiscount { min_quantity: i64, percent: i64 },
}

/// One product that can be bought.
#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
pub struct Product {
    pub sku: String,
    pub name: String,
    pub unit_price: Cents,
    /// `None` when the product has no promotion (`"promotion": null` in JSON).
    pub promotion: Option<Promotion>,
}

/// The list of products the shop sells. The checkout gets it as an input, so tests
/// can use a tiny made-up catalog or the standard one from `spec/catalog.json`.
#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
pub struct Catalog {
    pub products: Vec<Product>,
}

impl Catalog {
    pub fn new(products: Vec<Product>) -> Self {
        Catalog { products }
    }

    /// Finds a product by its SKU. SKUs are case sensitive, so `"apple"` does not find `"APPLE"`.
    ///
    /// Returns `None` when there is no such product (Rust has no `null`).
    pub fn find(&self, sku: &str) -> Option<&Product> {
        self.products.iter().find(|product| product.sku == sku)
    }

    /// Like [`Catalog::find`], but a missing product is an error: `UNKNOWN_PRODUCT`
    /// (rules §4.1 row 2).
    pub fn require(&self, sku: &str) -> Result<&Product, CheckoutError> {
        let _ = sku;
        todo!("Catalog::require() is not implemented yet")
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::ErrorCode;

    /// A small helper that builds a product without a promotion, to keep the tests short.
    fn product(sku: &str, unit_price: Cents) -> Product {
        Product {
            sku: sku.to_string(),
            name: format!("Test {sku}"),
            unit_price,
            promotion: None,
        }
    }

    #[test]
    fn find_returns_the_product_with_that_sku() {
        // Arrange
        let catalog = Catalog::new(vec![product("APPLE", 50), product("PEN", 120)]);

        // Act
        let found = catalog.find("PEN");

        // Assert: `Some(...)` means "found something". We compare the whole value.
        assert_eq!(found, Some(&product("PEN", 120)));
    }

    #[test]
    fn find_returns_none_for_an_unknown_sku() {
        let catalog = Catalog::new(vec![product("APPLE", 50)]);

        // `assert!(condition)` fails the test when the condition is false.
        assert!(catalog.find("BANANA").is_none());
    }

    #[test]
    fn find_is_case_sensitive() {
        let catalog = Catalog::new(vec![product("APPLE", 50)]);

        assert_eq!(catalog.find("apple"), None);
    }

    #[test]
    fn require_returns_the_product_when_it_exists() {
        let catalog = Catalog::new(vec![product("APPLE", 50), product("PEN", 120)]);

        assert_eq!(catalog.require("APPLE"), Ok(&product("APPLE", 50)));
    }

    #[test]
    fn require_fails_with_unknown_product_when_it_does_not() {
        let catalog = Catalog::new(vec![product("APPLE", 50)]);

        let error = catalog.require("BANANA").unwrap_err();

        assert_eq!(error.code, ErrorCode::UnknownProduct);
    }

    #[test]
    fn require_mentions_the_sku_in_the_error_message_so_it_is_useful_to_a_human() {
        let catalog = Catalog::new(vec![product("APPLE", 50)]);

        let error = catalog.require("BANANA").unwrap_err();

        assert_eq!(error.message, "Unknown product: BANANA");
    }

    #[test]
    fn promotions_are_read_from_json_using_their_type_field() {
        let json = r#"{ "type": "bulk_discount", "min_quantity": 10, "percent": 10 }"#;

        let promotion: Promotion = serde_json::from_str(json).unwrap();

        assert_eq!(
            promotion,
            Promotion::BulkDiscount {
                min_quantity: 10,
                percent: 10
            }
        );
    }
}
