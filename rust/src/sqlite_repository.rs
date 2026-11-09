//! An order repository backed by SQLite, through the `rusqlite` crate.
//!
//! Its tests are *integration tests* (a real database file), so they live in
//! `tests/sqlite_repository.rs` and use a temporary folder.

use crate::{
    Clock, Order, OrderRepository, Receipt, RepositoryError, format_timestamp, order_id,
    order_number,
};
use rusqlite::{Connection, OptionalExtension};
use std::path::Path;

/// Keeps orders in one table of an SQLite database file (or in memory).
/// The database row number becomes the order number, so ids are sequential and survive
/// closing and reopening the file.
///
/// There is no `close()` method: like every resource in Rust, the database is closed
/// automatically when the repository is dropped (goes out of scope, or `drop(repository)`).
pub struct SqliteOrderRepository {
    connection: Connection,
    clock: Box<dyn Clock>,
}

impl SqliteOrderRepository {
    /// Opens (and creates, if needed) the database file and its `orders` table.
    pub fn open(path: impl AsRef<Path>, clock: Box<dyn Clock>) -> Result<Self, RepositoryError> {
        Self::with_connection(Connection::open(path)?, clock)
    }

    /// A throw-away database that lives in memory only (SQLite's `:memory:`).
    pub fn in_memory(clock: Box<dyn Clock>) -> Result<Self, RepositoryError> {
        Self::with_connection(Connection::open_in_memory()?, clock)
    }

    /// Creates the `orders` table if this is a new database. The row number (`number`)
    /// becomes the order number; AUTOINCREMENT makes sure a number is never reused.
    fn with_connection(
        connection: Connection,
        clock: Box<dyn Clock>,
    ) -> Result<Self, RepositoryError> {
        connection.execute_batch(
            "CREATE TABLE IF NOT EXISTS orders (
                number     INTEGER PRIMARY KEY AUTOINCREMENT,
                created_at TEXT NOT NULL,
                receipt    TEXT NOT NULL
            )",
        )?;
        Ok(SqliteOrderRepository { connection, clock })
    }
}

impl OrderRepository for SqliteOrderRepository {
    fn save(&mut self, receipt: &Receipt) -> Result<String, RepositoryError> {
        // "?1" and "?2" are placeholders: the values travel separately from the SQL, so a
        // receipt can never be mistaken for a command (no SQL injection).
        self.connection.execute(
            "INSERT INTO orders (created_at, receipt) VALUES (?1, ?2)",
            (
                format_timestamp(self.clock.now()),
                serde_json::to_string(receipt)?,
            ),
        )?;
        // Row numbers start at 1 and only go up, so they are never negative.
        Ok(order_id(self.connection.last_insert_rowid().unsigned_abs()))
    }

    fn find(&self, id: &str) -> Result<Option<Order>, RepositoryError> {
        // Anything that is not an order id (or too big for SQLite's numbers) is simply not found.
        let Some(number) = order_number(id).and_then(|n| i64::try_from(n).ok()) else {
            return Ok(None);
        };
        let row: Option<(String, String)> = self
            .connection
            .query_row(
                "SELECT created_at, receipt FROM orders WHERE number = ?1",
                [number],
                |row| Ok((row.get(0)?, row.get(1)?)),
            )
            .optional()?; // no row is "not found", not an error
        let Some((created_at, receipt)) = row else {
            return Ok(None);
        };
        Ok(Some(Order {
            id: id.to_string(),
            created_at,
            receipt: serde_json::from_str(&receipt)?,
        }))
    }
}

impl From<rusqlite::Error> for RepositoryError {
    fn from(error: rusqlite::Error) -> Self {
        RepositoryError::new(error.to_string())
    }
}
