//! BDD: runs the Given/When/Then scenarios in `spec/features/*.feature` with the
//! `cucumber` crate.
//!
//! The `.feature` files are read straight from the shared `spec/features/` folder (never
//! copied), so Rust is held to exactly the same scenarios as the other five languages.
//! Every sentence in the step vocabulary of `spec/features/README.md` appears exactly once
//! below as a step definition; `fail_on_skipped()` makes sure a sentence without a
//! definition fails the run instead of being quietly skipped.
//!
//! The cucumber crate is `async`, so `main` drives it with a small executor from the
//! `futures` crate. `harness = false` in Cargo.toml tells Cargo this file has its own `main`.
//!
//! Run only this file with: `cargo test --test features`

mod common;

use checkout::{
    Basket, BasketItem, Catalog, Clock, ErrorCode, InMemoryOrderRepository, OrderRepository,
    Receipt, checkout, parse_timestamp, place_order,
};
use common::dollars::dollars_to_cents;
use common::fixtures::SharedClock;
use cucumber::gherkin::Step;
use cucumber::{Parameter, World, given, then, when};
use std::str::FromStr;

// --- The World: one per scenario, carries the state from step to step -------

/// Everything a scenario builds up: the basket, and what came out of the checkout.
#[derive(Debug, World)]
#[world(init = Self::new)]
pub struct CheckoutWorld {
    catalog: Catalog,
    basket: Basket,
    receipt: Option<Receipt>,
    error: Option<ErrorCode>,
    /// What the (replaceable) clock answers. "Given the clock is fixed at ..." changes it,
    /// and the repository sees the change because it holds a clone of the same clock.
    clock: SharedClock,
    repository: Option<InMemoryOrderRepository>,
    last_order_id: Option<String>,
}

impl CheckoutWorld {
    fn new() -> Self {
        CheckoutWorld {
            catalog: Catalog::new(Vec::new()),
            basket: Basket {
                items: Vec::new(),
                coupon: None,
                shipping: "standard".to_string(),
            },
            receipt: None,
            error: None,
            clock: SharedClock::starting_at(parse_timestamp("2026-01-01T00:00:00Z").unwrap()),
            repository: None,
            last_order_id: None,
        }
    }

    /// The receipt of the last checkout; fails the step when the checkout failed.
    fn receipt(&self) -> &Receipt {
        if let Some(code) = self.error {
            panic!("checkout failed with {code}");
        }
        self.receipt
            .as_ref()
            .expect("the customer has not checked out yet")
    }

    fn repository(&self) -> &InMemoryOrderRepository {
        self.repository
            .as_ref()
            .expect("no order repository - add 'Given an empty order repository'")
    }
}

/// A custom parameter type: `{money}` in a step matches `$12.39` and hands the step 1239 cents.
#[derive(Debug, Parameter)]
#[param(regex = r"\$\d+(?:\.\d{1,2})?", name = "money")]
struct Money(i64);

impl FromStr for Money {
    type Err = String;

    fn from_str(text: &str) -> Result<Self, Self::Err> {
        dollars_to_cents(text).map(Money)
    }
}

// --- Given ------------------------------------------------------------------

#[given("the standard catalog")]
fn the_standard_catalog(world: &mut CheckoutWorld) {
    world.catalog = common::load_catalog();
}

#[given("an empty basket")]
fn an_empty_basket(world: &mut CheckoutWorld) {
    world.basket.items.clear();
}

#[given("a basket containing:")]
fn a_basket_containing(world: &mut CheckoutWorld, step: &Step) {
    let table = step
        .table
        .as_ref()
        .expect("this step needs a sku / quantity table");
    let header = &table.rows[0];
    let column = |name: &str| {
        header
            .iter()
            .position(|cell| cell == name)
            .unwrap_or_else(|| panic!("the table has no {name:?} column"))
    };
    let (sku, quantity) = (column("sku"), column("quantity"));

    for row in &table.rows[1..] {
        world.basket.items.push(BasketItem {
            sku: row[sku].clone(),
            quantity: row[quantity]
                .parse()
                .unwrap_or_else(|_| panic!("quantity {:?} is not a whole number", row[quantity])),
        });
    }
}

#[given(expr = "the customer uses the coupon {string}")]
fn the_customer_uses_the_coupon(world: &mut CheckoutWorld, code: String) {
    world.basket.coupon = Some(code);
}

#[given(expr = "the customer chooses {string} shipping")]
fn the_customer_chooses_shipping(world: &mut CheckoutWorld, method: String) {
    world.basket.shipping = method;
}

#[given("an empty order repository")]
fn an_empty_order_repository(world: &mut CheckoutWorld) {
    // The repository gets a clone of the world's clock, so a later "the clock is fixed at"
    // step still applies to it.
    let clock: Box<dyn Clock> = Box::new(world.clock.clone());
    world.repository = Some(InMemoryOrderRepository::new(clock));
}

#[given(expr = "the clock is fixed at {string}")]
fn the_clock_is_fixed_at(world: &mut CheckoutWorld, timestamp: String) {
    let time = parse_timestamp(&timestamp)
        .unwrap_or_else(|| panic!("{timestamp:?} is not a timestamp like 2026-01-15T10:30:00Z"));
    world.clock.set(time);
}

// --- When -------------------------------------------------------------------

#[when("the customer checks out")]
fn the_customer_checks_out(world: &mut CheckoutWorld) {
    match checkout(&world.basket, &world.catalog) {
        Ok(receipt) => world.receipt = Some(receipt),
        Err(error) => world.error = Some(error.code),
    }
}

#[when("the customer checks out and saves the order")]
fn the_customer_checks_out_and_saves_the_order(world: &mut CheckoutWorld) {
    let repository = world
        .repository
        .as_mut()
        .expect("no order repository - add 'Given an empty order repository'");
    match place_order(&world.basket, &world.catalog, repository) {
        Ok(placed) => {
            world.receipt = Some(placed.receipt);
            world.last_order_id = Some(placed.order_id);
        }
        Err(error) => {
            world.error = Some(
                error
                    .checkout_code()
                    .unwrap_or_else(|| panic!("saving failed unexpectedly: {error}")),
            );
        }
    }
}

// --- Then: the receipt ------------------------------------------------------

#[then(expr = "the subtotal is {money}")]
fn the_subtotal_is(world: &mut CheckoutWorld, amount: Money) {
    assert_eq!(world.receipt().subtotal, amount.0);
}

#[then(expr = "the promotion discount is {money}")]
fn the_promotion_discount_is(world: &mut CheckoutWorld, amount: Money) {
    assert_eq!(world.receipt().promotion_discount, amount.0);
}

#[then(expr = "the coupon discount is {money}")]
fn the_coupon_discount_is(world: &mut CheckoutWorld, amount: Money) {
    assert_eq!(world.receipt().coupon_discount, amount.0);
}

#[then(expr = "the shipping cost is {money}")]
fn the_shipping_cost_is(world: &mut CheckoutWorld, amount: Money) {
    assert_eq!(world.receipt().shipping, amount.0);
}

#[then(expr = "the tax is {money}")]
fn the_tax_is(world: &mut CheckoutWorld, amount: Money) {
    assert_eq!(world.receipt().tax, amount.0);
}

#[then(expr = "the total is {money}")]
fn the_total_is(world: &mut CheckoutWorld, amount: Money) {
    assert_eq!(world.receipt().total, amount.0);
}

#[then(expr = "the receipt has {int} lines")]
fn the_receipt_has_lines(world: &mut CheckoutWorld, count: usize) {
    assert_eq!(world.receipt().lines.len(), count);
}

#[then(expr = "the receipt line for {string} has a quantity of {int}")]
fn the_receipt_line_for_has_a_quantity_of(world: &mut CheckoutWorld, sku: String, quantity: i64) {
    let line = world
        .receipt()
        .lines
        .iter()
        .find(|line| line.sku == sku)
        .unwrap_or_else(|| panic!("no receipt line for {sku}"));
    assert_eq!(line.quantity, quantity);
}

#[then(expr = "checkout fails with the error {string}")]
fn checkout_fails_with_the_error(world: &mut CheckoutWorld, code: String) {
    assert!(
        world.receipt.is_none(),
        "checkout succeeded but was expected to fail"
    );
    assert_eq!(world.error.map(ErrorCode::as_str), Some(code.as_str()));
}

// --- Then: saved orders -----------------------------------------------------

#[then(expr = "the order is saved with the id {string}")]
fn the_order_is_saved_with_the_id(world: &mut CheckoutWorld, id: String) {
    assert_eq!(world.last_order_id.as_deref(), Some(id.as_str()));
}

#[then(expr = "the saved order {string} has a total of {money}")]
fn the_saved_order_has_a_total_of(world: &mut CheckoutWorld, id: String, amount: Money) {
    let order = world
        .repository()
        .find(&id)
        .expect("the repository should not fail")
        .unwrap_or_else(|| panic!("order {id} not found"));
    assert_eq!(order.receipt.total, amount.0);
}

#[then("the order repository is empty")]
fn the_order_repository_is_empty(world: &mut CheckoutWorld) {
    assert_eq!(world.repository().count(), 0);
}

#[then(expr = "the order {string} cannot be found")]
fn the_order_cannot_be_found(world: &mut CheckoutWorld, id: String) {
    assert_eq!(world.repository().find(&id), Ok(None));
}

#[then(expr = "the saved order {string} was created at {string}")]
fn the_saved_order_was_created_at(world: &mut CheckoutWorld, id: String, timestamp: String) {
    let order = world
        .repository()
        .find(&id)
        .expect("the repository should not fail")
        .unwrap_or_else(|| panic!("order {id} not found"));
    // Timestamps are compared as strings, exactly as spec/features/README.md says.
    assert_eq!(order.created_at, timestamp);
}

// --- The runner -------------------------------------------------------------

fn main() {
    let features = common::spec_dir().join("features");
    futures::executor::block_on(
        CheckoutWorld::cucumber()
            .fail_on_skipped()
            .run_and_exit(features),
    );
}
