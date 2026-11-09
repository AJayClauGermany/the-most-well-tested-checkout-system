//! An order repository that writes one JSON file per order into a folder.
//!
//! Its tests are *integration tests* (they touch the real file system), so they live in
//! `tests/json_file_repository.rs` and use a temporary folder.

use crate::{
    Clock, Order, OrderRepository, Receipt, RepositoryError, format_timestamp, order_id,
    order_number,
};
use std::fs;
use std::path::{Path, PathBuf};

/// Keeps each order as `<folder>/ORDER-<n>.json`, in the shape of
/// `spec/schema/order.schema.json`. Ids continue where the files in the folder left off,
/// so reopening the same folder later keeps numbering correctly.
pub struct JsonFileOrderRepository {
    folder: PathBuf,
    clock: Box<dyn Clock>,
}

impl JsonFileOrderRepository {
    /// Opens (and creates, if needed) the folder the orders are kept in.
    pub fn open(folder: impl AsRef<Path>, clock: Box<dyn Clock>) -> Result<Self, RepositoryError> {
        let folder = folder.as_ref().to_path_buf();
        fs::create_dir_all(&folder)?;
        Ok(JsonFileOrderRepository { folder, clock })
    }

    /// The folder the orders are kept in.
    pub fn folder(&self) -> &Path {
        &self.folder
    }

    fn path_of(&self, id: &str) -> PathBuf {
        self.folder.join(format!("{id}.json"))
    }

    /// The highest number among the `ORDER-<n>.json` files in the folder, or 0 when there
    /// are none. Numbers, not names, are compared: `ORDER-12` comes after `ORDER-3`.
    fn highest_order_number(&self) -> Result<u64, RepositoryError> {
        let mut highest = 0;
        for entry in fs::read_dir(&self.folder)? {
            let entry = entry?;
            let name = entry.file_name();
            let number = name
                .to_string_lossy()
                .strip_suffix(".json")
                .and_then(order_number);
            if let Some(number) = number.filter(|_| entry.path().is_file()) {
                highest = highest.max(number);
            }
        }
        Ok(highest)
    }
}

impl OrderRepository for JsonFileOrderRepository {
    fn save(&mut self, receipt: &Receipt) -> Result<String, RepositoryError> {
        let order = Order {
            id: order_id(self.highest_order_number()? + 1),
            created_at: format_timestamp(self.clock.now()),
            receipt: receipt.clone(),
        };
        // Indented, with a newline at the end: friendly to humans and to `git diff`.
        let mut text = serde_json::to_string_pretty(&order)?;
        text.push('\n');
        fs::write(self.path_of(&order.id), text)?;
        Ok(order.id)
    }

    fn find(&self, id: &str) -> Result<Option<Order>, RepositoryError> {
        // Only accept real order ids: "../secret" must not read files outside the folder.
        if order_number(id).is_none() {
            return Ok(None);
        }
        let path = self.path_of(id);
        if !path.is_file() {
            return Ok(None);
        }
        let text = fs::read_to_string(path)?;
        Ok(Some(serde_json::from_str(&text)?))
    }
}
