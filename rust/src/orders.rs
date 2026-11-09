//! Saving orders (rules §6): the order type, the clock, the repository trait and the
//! simplest repository there is – a list in memory.

use crate::{Basket, Catalog, CheckoutError, ErrorCode, Receipt, checkout};
use serde::{Deserialize, Serialize};
use std::fmt;
use std::time::{Duration, SystemTime, UNIX_EPOCH};

/// A saved order, exactly as `spec/schema/order.schema.json` describes it.
#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
pub struct Order {
    pub id: String,
    /// UTC, ISO 8601, whole seconds, trailing Z: `"2026-01-15T10:30:00Z"`.
    pub created_at: String,
    pub receipt: Receipt,
}

// ---------------------------------------------------------------------------
// Time
// ---------------------------------------------------------------------------

/// Anything that answers "what time is it?".
///
/// The rules (§6) demand that the time source can be replaced, because `created_at` is
/// the one value that changes from run to run. Production uses [`SystemClock`]; tests use
/// [`FixedClock`] (a *stub*) or the `MockClock` that `mockall` generates from this trait.
#[cfg_attr(test, mockall::automock)]
pub trait Clock {
    fn now(&self) -> SystemTime;
}

/// The real clock.
#[derive(Debug, Clone, Copy, Default)]
pub struct SystemClock;

impl Clock for SystemClock {
    fn now(&self) -> SystemTime {
        todo!("SystemClock::now() is not implemented yet")
    }
}

/// A clock that always answers the same time – the textbook *stub*.
#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub struct FixedClock {
    time: SystemTime,
}

impl FixedClock {
    pub fn new(time: SystemTime) -> Self {
        FixedClock { time }
    }

    /// A fixed clock from a timestamp written like `created_at`: `"2026-01-15T10:30:00Z"`.
    /// `None` when the text is not such a timestamp.
    pub fn at(timestamp: &str) -> Option<Self> {
        let _ = timestamp;
        todo!("FixedClock::at() is not implemented yet")
    }
}

impl Clock for FixedClock {
    fn now(&self) -> SystemTime {
        todo!("FixedClock::now() is not implemented yet")
    }
}

/// Formats a time the way `created_at` must look: UTC, whole seconds, trailing Z.
///
/// Done with the standard library only (no date crate): the number of seconds since
/// 1970 is turned into a calendar date with whole-number maths.
///
/// ```
/// use checkout::format_timestamp;
/// use std::time::{Duration, UNIX_EPOCH};
///
/// assert_eq!(format_timestamp(UNIX_EPOCH), "1970-01-01T00:00:00Z");
/// assert_eq!(format_timestamp(UNIX_EPOCH + Duration::from_secs(1768473000)), "2026-01-15T10:30:00Z");
/// ```
pub fn format_timestamp(time: SystemTime) -> String {
    let _ = time;
    todo!("format_timestamp() is not implemented yet")
}

/// The inverse of [`format_timestamp`]: reads `"2026-01-15T10:30:00Z"`.
/// `None` for anything that is not exactly that shape (or not a real date).
pub fn parse_timestamp(text: &str) -> Option<SystemTime> {
    let _ = text;
    todo!("parse_timestamp() is not implemented yet")
}

// ---------------------------------------------------------------------------
// Order ids
// ---------------------------------------------------------------------------

/// The id of the n-th order saved into a repository: `"ORDER-1"`, `"ORDER-2"`, ...
pub fn order_id(n: u64) -> String {
    let _ = n;
    todo!("order_id() is not implemented yet")
}

/// The inverse of [`order_id`]: `"ORDER-7"` gives `Some(7)`. Anything that is not an
/// order id (as `order.schema.json` spells it: no leading zero, no `ORDER-0`) gives `None`.
pub fn order_number(id: &str) -> Option<u64> {
    let _ = id;
    todo!("order_number() is not implemented yet")
}

// ---------------------------------------------------------------------------
// The repository
// ---------------------------------------------------------------------------

/// Something went wrong *storing* or *reading* an order (a file could not be written,
/// the database refused the row). Not a checkout error: the basket was fine.
#[derive(Debug, Clone, PartialEq, Eq)]
pub struct RepositoryError {
    pub message: String,
}

impl RepositoryError {
    pub fn new(message: impl Into<String>) -> Self {
        RepositoryError {
            message: message.into(),
        }
    }
}

impl fmt::Display for RepositoryError {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        f.write_str(&self.message)
    }
}

impl std::error::Error for RepositoryError {}

impl From<std::io::Error> for RepositoryError {
    fn from(error: std::io::Error) -> Self {
        RepositoryError::new(error.to_string())
    }
}

impl From<serde_json::Error> for RepositoryError {
    fn from(error: serde_json::Error) -> Self {
        RepositoryError::new(error.to_string())
    }
}

/// Where orders are kept (rules §6). Implemented in memory, as JSON files and in SQLite.
///
/// `#[cfg_attr(test, mockall::automock)]` generates a `MockOrderRepository` for the unit
/// tests, which can check *how* it was called (for example: `save` was never called).
#[cfg_attr(test, mockall::automock)]
pub trait OrderRepository {
    /// Saves a receipt and returns the new order's id.
    fn save(&mut self, receipt: &Receipt) -> Result<String, RepositoryError>;

    /// Looks an order up. `Ok(None)` means "not found" – a normal answer, not an error.
    fn find(&self, id: &str) -> Result<Option<Order>, RepositoryError>;
}

/// Keeps orders in a plain `Vec`. Doubles as the hand-written *fake* for tests.
pub struct InMemoryOrderRepository {
    orders: Vec<Order>,
    clock: Box<dyn Clock>,
}

impl InMemoryOrderRepository {
    /// An empty repository that stamps orders with this clock's time.
    pub fn new(clock: Box<dyn Clock>) -> Self {
        InMemoryOrderRepository {
            orders: Vec::new(),
            clock,
        }
    }

    /// How many orders are saved. Handy for tests ("the order repository is empty").
    pub fn count(&self) -> usize {
        todo!("InMemoryOrderRepository::count() is not implemented yet")
    }
}

/// Shows the number of orders (a `Box<dyn Clock>` cannot be printed, so this is hand-written).
impl fmt::Debug for InMemoryOrderRepository {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        f.debug_struct("InMemoryOrderRepository")
            .field("orders", &self.orders.len())
            .finish_non_exhaustive()
    }
}

/// `InMemoryOrderRepository::default()` uses the real clock.
impl Default for InMemoryOrderRepository {
    fn default() -> Self {
        InMemoryOrderRepository::new(Box::new(SystemClock))
    }
}

impl OrderRepository for InMemoryOrderRepository {
    fn save(&mut self, receipt: &Receipt) -> Result<String, RepositoryError> {
        let _ = (receipt, &self.clock);
        todo!("InMemoryOrderRepository::save() is not implemented yet")
    }

    fn find(&self, id: &str) -> Result<Option<Order>, RepositoryError> {
        let _ = id;
        todo!("InMemoryOrderRepository::find() is not implemented yet")
    }
}

// ---------------------------------------------------------------------------
// Placing an order = checkout + save
// ---------------------------------------------------------------------------

/// What [`place_order`] gives back: the receipt and the id the repository assigned.
#[derive(Debug, Clone, PartialEq, Eq)]
pub struct PlacedOrder {
    pub order_id: String,
    pub receipt: Receipt,
}

/// Why [`place_order`] failed: either the basket was invalid, or saving went wrong.
#[derive(Debug, Clone, PartialEq, Eq)]
pub enum PlaceOrderError {
    Checkout(CheckoutError),
    Repository(RepositoryError),
}

impl PlaceOrderError {
    /// The checkout error code, when the basket was the problem.
    pub fn checkout_code(&self) -> Option<ErrorCode> {
        match self {
            PlaceOrderError::Checkout(error) => Some(error.code),
            PlaceOrderError::Repository(_) => None,
        }
    }
}

impl fmt::Display for PlaceOrderError {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        match self {
            PlaceOrderError::Checkout(error) => error.fmt(f),
            PlaceOrderError::Repository(error) => write!(f, "could not save the order: {error}"),
        }
    }
}

impl std::error::Error for PlaceOrderError {}

impl From<CheckoutError> for PlaceOrderError {
    fn from(error: CheckoutError) -> Self {
        PlaceOrderError::Checkout(error)
    }
}

impl From<RepositoryError> for PlaceOrderError {
    fn from(error: RepositoryError) -> Self {
        PlaceOrderError::Repository(error)
    }
}

/// Checks out, then saves the receipt as an order. A failed checkout saves nothing:
/// the error is returned before the repository is touched.
///
/// `R: OrderRepository + ?Sized` means any repository works, including `&mut dyn OrderRepository`.
pub fn place_order<R: OrderRepository + ?Sized>(
    basket: &Basket,
    catalog: &Catalog,
    repository: &mut R,
) -> Result<PlacedOrder, PlaceOrderError> {
    let _ = (basket, catalog, repository, checkout);
    todo!("place_order() is not implemented yet")
}

// Keeps the unused-import warnings away while the functions above are still stubs.
#[allow(dead_code)]
fn unused_until_implemented(_: Duration, _: SystemTime) -> SystemTime {
    UNIX_EPOCH
}

// ---------------------------------------------------------------------------
// Unit tests – and the project's showcase of TEST DOUBLES.
//
// A *test double* stands in for a real dependency:
//   * a FAKE is a simple working version – InMemoryOrderRepository is one,
//   * a STUB gives canned answers – FixedClock always says the same time,
//   * a MOCK records how it was called – mockall's MockOrderRepository lets us check
//     that save() was never called.
// ---------------------------------------------------------------------------
#[cfg(test)]
mod tests {
    use super::*;
    use crate::test_fixtures::{basket, one_apple_receipt, tiny_catalog};
    use pretty_assertions::assert_eq;
    use rstest::rstest;

    const FIXED_TIME: &str = "2026-01-15T10:30:00Z";

    /// A stub clock: always 10:30 on 15 January 2026, UTC.
    fn fixed_clock() -> Box<dyn Clock> {
        Box::new(FixedClock::at(FIXED_TIME).unwrap())
    }

    fn seconds(unix_seconds: u64) -> SystemTime {
        UNIX_EPOCH + Duration::from_secs(unix_seconds)
    }

    // --- format_timestamp / parse_timestamp -------------------------------

    #[rstest]
    #[case::the_epoch(0, "1970-01-01T00:00:00Z")]
    #[case::the_year_2000(946_684_800, "2000-01-01T00:00:00Z")]
    #[case::the_fixed_test_time(1_768_473_000, "2026-01-15T10:30:00Z")]
    #[case::a_leap_day(1_709_208_000, "2024-02-29T12:00:00Z")]
    #[case::the_last_second_of_2099(4_102_444_799, "2099-12-31T23:59:59Z")]
    fn format_timestamp_writes_utc_with_whole_seconds_and_a_trailing_z(
        #[case] unix_seconds: u64,
        #[case] expected: &str,
    ) {
        assert_eq!(format_timestamp(seconds(unix_seconds)), expected);
    }

    #[test]
    fn format_timestamp_drops_fractions_of_a_second_instead_of_rounding_them() {
        let time = seconds(1_768_473_000) + Duration::from_millis(999);

        assert_eq!(format_timestamp(time), "2026-01-15T10:30:00Z");
    }

    #[rstest]
    #[case("1970-01-01T00:00:00Z", 0)]
    #[case("2000-01-01T00:00:00Z", 946_684_800)]
    #[case("2026-01-15T10:30:00Z", 1_768_473_000)]
    #[case("2024-02-29T12:00:00Z", 1_709_208_000)]
    #[case("2099-12-31T23:59:59Z", 4_102_444_799)]
    fn parse_timestamp_reads_what_format_timestamp_writes(
        #[case] text: &str,
        #[case] unix_seconds: u64,
    ) {
        assert_eq!(parse_timestamp(text), Some(seconds(unix_seconds)));
        assert_eq!(format_timestamp(parse_timestamp(text).unwrap()), text);
    }

    #[rstest]
    #[case::empty("")]
    #[case::space_instead_of_t("2026-01-15 10:30:00Z")]
    #[case::no_z("2026-01-15T10:30:00")]
    #[case::milliseconds("2026-01-15T10:30:00.000Z")]
    #[case::two_digit_year("26-01-15T10:30:00Z")]
    #[case::month_13("2026-13-01T00:00:00Z")]
    #[case::month_0("2026-00-01T00:00:00Z")]
    #[case::day_0("2026-01-00T00:00:00Z")]
    #[case::february_30th("2026-02-30T00:00:00Z")]
    #[case::february_29th_in_a_non_leap_year("2026-02-29T00:00:00Z")]
    #[case::hour_24("2026-01-15T24:00:00Z")]
    #[case::minute_60("2026-01-15T10:60:00Z")]
    #[case::second_60("2026-01-15T10:30:60Z")]
    #[case::letters("2026-01-15T10:30:0aZ")]
    #[case::trailing_text("2026-01-15T10:30:00Zzz")]
    fn parse_timestamp_rejects_anything_that_is_not_a_timestamp(#[case] text: &str) {
        assert_eq!(parse_timestamp(text), None);
    }

    #[test]
    fn parse_timestamp_accepts_february_29th_in_a_leap_year() {
        assert_eq!(
            parse_timestamp("2024-02-29T00:00:00Z"),
            Some(seconds(1_709_164_800))
        );
        // 2000 was a leap year (divisible by 400), 1900 was not (divisible by 100).
        assert!(parse_timestamp("2000-02-29T00:00:00Z").is_some());
        assert!(parse_timestamp("1900-02-29T00:00:00Z").is_none());
    }

    // --- order ids ---------------------------------------------------------

    #[rstest]
    #[case(1, "ORDER-1")]
    #[case(2, "ORDER-2")]
    #[case(10, "ORDER-10")]
    #[case(1234, "ORDER-1234")]
    fn order_id_and_order_number_are_inverses(#[case] n: u64, #[case] id: &str) {
        assert_eq!(order_id(n), id);
        assert_eq!(order_number(id), Some(n));
    }

    #[rstest]
    #[case("")]
    #[case("ORDER-")]
    #[case("ORDER-0")]
    #[case("ORDER-01")]
    #[case("ORDER--1")]
    #[case("ORDER-1.5")]
    #[case("ORDER-1 ")]
    #[case(" ORDER-1")]
    #[case("MYORDER-1")]
    #[case("ORDER-1x")]
    #[case("order-1")]
    #[case("ORDER-x")]
    #[case("1")]
    #[case("ORDER_1")]
    #[case("ORDER-١")] // an Arabic-Indic digit is not an ASCII digit
    #[case("ORDER-99999999999999999999")] // too big for a u64
    fn order_number_does_not_recognise_text_that_is_not_an_order_id(#[case] id: &str) {
        assert_eq!(order_number(id), None);
    }

    // --- clocks ------------------------------------------------------------

    #[test]
    fn the_system_clock_reads_the_real_time() {
        let before = SystemTime::now();

        let now = SystemClock.now();

        let after = SystemTime::now();
        assert!(before <= now && now <= after);
    }

    #[test]
    fn a_fixed_clock_always_answers_the_same_time() {
        let clock = FixedClock::at(FIXED_TIME).unwrap();

        assert_eq!(clock.now(), seconds(1_768_473_000));
        assert_eq!(clock.now(), seconds(1_768_473_000));
        assert_eq!(FixedClock::new(seconds(5)).now(), seconds(5));
    }

    #[test]
    fn a_fixed_clock_cannot_be_made_from_text_that_is_not_a_timestamp() {
        assert_eq!(FixedClock::at("tomorrow"), None);
    }

    // --- InMemoryOrderRepository (a hand-written fake) ---------------------

    #[test]
    fn starts_empty() {
        let repository = InMemoryOrderRepository::new(fixed_clock());

        assert_eq!(repository.count(), 0);
        assert_eq!(repository.find("ORDER-1"), Ok(None));
    }

    #[test]
    fn numbers_orders_order_1_2_3_in_the_order_they_are_saved() {
        let mut repository = InMemoryOrderRepository::new(fixed_clock());

        let ids = [
            repository.save(&one_apple_receipt()).unwrap(),
            repository.save(&one_apple_receipt()).unwrap(),
            repository.save(&one_apple_receipt()).unwrap(),
        ];

        assert_eq!(ids, ["ORDER-1", "ORDER-2", "ORDER-3"]);
        assert_eq!(repository.count(), 3);
    }

    #[test]
    fn finds_a_saved_order_with_its_receipt_and_the_time_the_clock_said() {
        let mut repository = InMemoryOrderRepository::new(fixed_clock());

        let id = repository.save(&one_apple_receipt()).unwrap();

        assert_eq!(
            repository.find(&id),
            Ok(Some(Order {
                id: "ORDER-1".to_string(),
                created_at: FIXED_TIME.to_string(),
                receipt: one_apple_receipt(),
            }))
        );
    }

    #[test]
    fn asks_the_clock_at_the_moment_of_saving_once_per_save() {
        // A MOCK clock: it answers like the stub, and mockall checks it was called exactly twice.
        let mut clock = MockClock::new();
        clock
            .expect_now()
            .times(2)
            .returning(|| seconds(1_768_473_000));
        let mut repository = InMemoryOrderRepository::new(Box::new(clock));

        repository.save(&one_apple_receipt()).unwrap();
        repository.save(&one_apple_receipt()).unwrap();
        // The expectation is checked when `clock` is dropped at the end of the test.
    }

    #[test]
    fn stamps_each_order_with_the_clocks_time_at_the_moment_of_saving() {
        let mut clock = MockClock::new();
        let mut times = [1_768_473_000, 1_768_554_000].into_iter(); // Jan 15 10:30, Jan 16 09:00
        clock
            .expect_now()
            .returning(move || seconds(times.next().unwrap()));
        let mut repository = InMemoryOrderRepository::new(Box::new(clock));

        let first = repository.save(&one_apple_receipt()).unwrap();
        let second = repository.save(&one_apple_receipt()).unwrap();

        assert_eq!(
            repository.find(&first).unwrap().unwrap().created_at,
            "2026-01-15T10:30:00Z"
        );
        assert_eq!(
            repository.find(&second).unwrap().unwrap().created_at,
            "2026-01-16T09:00:00Z"
        );
    }

    #[test]
    fn keeps_a_copy_of_the_receipt_so_changing_the_original_later_does_not_change_the_order() {
        let mut repository = InMemoryOrderRepository::new(fixed_clock());
        let mut receipt = one_apple_receipt();

        let id = repository.save(&receipt).unwrap();
        receipt.total = 1;
        receipt.lines[0].quantity = 42;

        assert_eq!(
            repository.find(&id).unwrap().unwrap().receipt,
            one_apple_receipt()
        );
    }

    #[rstest]
    #[case("ORDER-2")]
    #[case("ORDER-0")]
    #[case("ORDER-99")]
    #[case("order-1")]
    #[case("")]
    fn answers_not_found_for_ids_that_were_never_handed_out(#[case] id: &str) {
        let mut repository = InMemoryOrderRepository::new(fixed_clock());
        repository.save(&one_apple_receipt()).unwrap();

        assert_eq!(repository.find(id), Ok(None));
    }

    #[test]
    fn uses_the_real_clock_by_default() {
        let mut repository = InMemoryOrderRepository::default();
        let before = SystemTime::now();

        let id = repository.save(&one_apple_receipt()).unwrap();

        let created_at = repository.find(&id).unwrap().unwrap().created_at;
        let stamped = parse_timestamp(&created_at).expect("created_at should be a timestamp");
        // Whole seconds, so allow one second of slack on the early side.
        assert!(stamped + Duration::from_secs(1) >= before);
        assert!(stamped <= SystemTime::now());
    }

    // --- place_order -------------------------------------------------------

    #[test]
    fn checks_out_and_saves_the_receipt_returning_both_the_receipt_and_the_new_id() {
        let catalog = tiny_catalog();
        let mut repository = InMemoryOrderRepository::new(fixed_clock());

        let placed = place_order(&basket(&[("APPLE", 1)]), &catalog, &mut repository).unwrap();

        assert_eq!(placed.order_id, "ORDER-1");
        assert_eq!(
            placed.receipt,
            checkout(&basket(&[("APPLE", 1)]), &catalog).unwrap()
        );
        assert_eq!(
            repository.find("ORDER-1").unwrap().unwrap().receipt,
            placed.receipt
        );
    }

    #[test]
    fn saves_exactly_what_it_returns_checked_with_a_mock_repository() {
        // mockall's MockOrderRepository records every call, so the test can see exactly
        // how place_order used it: save() once, with the receipt it returns; find() never.
        let catalog = tiny_catalog();
        let expected_receipt = checkout(&basket(&[("TV", 1)]), &catalog).unwrap();
        let mut repository = MockOrderRepository::new();
        repository
            .expect_save()
            .withf(move |receipt| *receipt == expected_receipt)
            .times(1)
            .returning(|_| Ok("ORDER-7".to_string()));
        repository.expect_find().never();

        let placed = place_order(&basket(&[("TV", 1)]), &catalog, &mut repository).unwrap();

        assert_eq!(placed.order_id, "ORDER-7"); // the id comes from the repository, not from place_order
    }

    #[test]
    fn saves_nothing_when_checkout_fails_p14() {
        let catalog = tiny_catalog();
        let mut repository = MockOrderRepository::new();
        repository.expect_save().never(); // the same as `.times(0)`
        repository.expect_find().never();

        let empty = place_order(&basket(&[]), &catalog, &mut repository).unwrap_err();
        let unknown =
            place_order(&basket(&[("BANANA", 1)]), &catalog, &mut repository).unwrap_err();

        assert_eq!(empty.checkout_code(), Some(ErrorCode::EmptyBasket));
        assert_eq!(unknown.checkout_code(), Some(ErrorCode::UnknownProduct));
    }

    #[test]
    fn leaves_the_fake_repository_empty_after_a_failed_checkout() {
        let catalog = tiny_catalog();
        let mut repository = InMemoryOrderRepository::new(fixed_clock());

        let error = place_order(&basket(&[("APPLE", 0)]), &catalog, &mut repository).unwrap_err();

        assert_eq!(error.checkout_code(), Some(ErrorCode::InvalidQuantity));
        assert_eq!(repository.count(), 0);
        assert_eq!(repository.find("ORDER-1"), Ok(None));
    }

    #[test]
    fn passes_a_repository_failure_on_instead_of_hiding_it() {
        let catalog = tiny_catalog();
        let mut repository = MockOrderRepository::new();
        repository
            .expect_save()
            .times(1)
            .returning(|_| Err(RepositoryError::new("disk full")));

        let error = place_order(&basket(&[("APPLE", 1)]), &catalog, &mut repository).unwrap_err();

        assert_eq!(
            error,
            PlaceOrderError::Repository(RepositoryError::new("disk full"))
        );
        assert_eq!(error.checkout_code(), None);
        assert_eq!(error.to_string(), "could not save the order: disk full");
    }

    #[test]
    fn a_checkout_failure_reads_like_the_checkout_error() {
        let error = PlaceOrderError::from(CheckoutError::with_message(
            ErrorCode::UnknownCoupon,
            "Unknown coupon: NOPE",
        ));

        assert_eq!(
            error.to_string(),
            "checkout failed: UNKNOWN_COUPON: Unknown coupon: NOPE"
        );
    }

    #[test]
    fn repository_errors_can_be_made_from_io_and_json_errors() {
        let io = std::io::Error::new(std::io::ErrorKind::NotFound, "no such folder");
        let json = serde_json::from_str::<Order>("not json").unwrap_err();

        assert_eq!(RepositoryError::from(io).message, "no such folder");
        assert!(RepositoryError::from(json).message.contains("expected"));
    }
}
