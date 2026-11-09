//! An order repository that writes one JSON file per order into a folder.
//!
//! Its tests are *integration tests* (they touch the real file system), so they live in
//! `tests/json_file_repository.rs` and use a temporary folder.

use crate::{Clock, Order, OrderRepository, Receipt, RepositoryError};
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
        let _ = (folder.as_ref(), clock);
        todo!("JsonFileOrderRepository::open() is not implemented yet")
    }

    /// The folder the orders are kept in.
    pub fn folder(&self) -> &Path {
        &self.folder
    }
}

impl OrderRepository for JsonFileOrderRepository {
    fn save(&mut self, receipt: &Receipt) -> Result<String, RepositoryError> {
        let _ = (receipt, &self.clock);
        todo!("JsonFileOrderRepository::save() is not implemented yet")
    }

    fn find(&self, id: &str) -> Result<Option<Order>, RepositoryError> {
        let _ = id;
        todo!("JsonFileOrderRepository::find() is not implemented yet")
    }
}
