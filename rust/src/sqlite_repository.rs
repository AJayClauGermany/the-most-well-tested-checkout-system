//! An order repository backed by SQLite, through the `rusqlite` crate.
//!
//! Its tests are *integration tests* (a real database file), so they live in
//! `tests/sqlite_repository.rs` and use a temporary folder.

use crate::{Clock, Order, OrderRepository, Receipt, RepositoryError};
use rusqlite::Connection;
use std::path::Path;

/// Keeps orders in one table of an SQLite database file (or in memory).
/// The database row number becomes the order number, so ids are sequential and survive
/// closing and reopening the file.
pub struct SqliteOrderRepository {
    connection: Connection,
    clock: Box<dyn Clock>,
}

impl SqliteOrderRepository {
    /// Opens (and creates, if needed) the database file and its `orders` table.
    pub fn open(path: impl AsRef<Path>, clock: Box<dyn Clock>) -> Result<Self, RepositoryError> {
        let _ = (path.as_ref(), clock);
        todo!("SqliteOrderRepository::open() is not implemented yet")
    }

    /// A throw-away database that lives in memory only (SQLite's `:memory:`).
    pub fn in_memory(clock: Box<dyn Clock>) -> Result<Self, RepositoryError> {
        let _ = clock;
        todo!("SqliteOrderRepository::in_memory() is not implemented yet")
    }

    /// Closes the database. It takes `self` by value, so a closed repository cannot be
    /// used by mistake – the compiler refuses.
    pub fn close(self) -> Result<(), RepositoryError> {
        let _ = (self.connection, self.clock);
        todo!("SqliteOrderRepository::close() is not implemented yet")
    }
}

impl OrderRepository for SqliteOrderRepository {
    fn save(&mut self, receipt: &Receipt) -> Result<String, RepositoryError> {
        let _ = receipt;
        todo!("SqliteOrderRepository::save() is not implemented yet")
    }

    fn find(&self, id: &str) -> Result<Option<Order>, RepositoryError> {
        let _ = id;
        todo!("SqliteOrderRepository::find() is not implemented yet")
    }
}

impl From<rusqlite::Error> for RepositoryError {
    fn from(error: rusqlite::Error) -> Self {
        RepositoryError::new(error.to_string())
    }
}
