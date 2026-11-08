//! The basket: the input to the checkout (rules §3).

use serde::{Deserialize, Serialize};

/// One line in the basket: which product, and how many.
#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
pub struct BasketItem {
    pub sku: String,
    /// Signed on purpose: invalid quantities like `0` or `-1` must be representable,
    /// so that the checkout can reject them with `INVALID_QUANTITY`.
    pub quantity: i64,
}

/// What the customer wants to buy.
#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
pub struct Basket {
    pub items: Vec<BasketItem>,
    /// `None` when there is no coupon. A missing field and `null` both become `None`.
    #[serde(default)]
    pub coupon: Option<String>,
    /// `"standard"` or `"express"`. It is a `String` (not an enum) because invalid
    /// values like `"overnight"` must reach the checkout so it can report
    /// `UNKNOWN_SHIPPING_METHOD`.
    #[serde(default = "default_shipping")]
    pub shipping: String,
}

/// Used by serde when the `"shipping"` field is missing from the JSON.
fn default_shipping() -> String {
    "standard".to_string()
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn missing_coupon_and_shipping_get_their_defaults() {
        let json = r#"{ "items": [{ "sku": "APPLE", "quantity": 2 }] }"#;

        let basket: Basket = serde_json::from_str(json).unwrap();

        assert_eq!(basket.coupon, None);
        assert_eq!(basket.shipping, "standard");
    }
}
