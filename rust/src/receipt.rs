//! The receipt: the output of a successful checkout (rules §5).

use crate::Cents;
use serde::{Deserialize, Serialize};

/// One line of the receipt, after duplicate basket lines were merged.
#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
pub struct ReceiptLine {
    pub sku: String,
    pub name: String,
    pub quantity: i64,
    pub unit_price: Cents,
    /// `unit_price × quantity`, before any promotion.
    pub line_total: Cents,
    /// How much the product's promotion takes off this line (`0` if none).
    pub line_discount: Cents,
}

/// The result of a successful checkout. All amounts are whole cents.
#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
pub struct Receipt {
    /// In order of first appearance in the basket.
    pub lines: Vec<ReceiptLine>,
    pub subtotal: Cents,
    pub promotion_discount: Cents,
    pub coupon_discount: Cents,
    pub shipping: Cents,
    pub tax: Cents,
    pub total: Cents,
}
