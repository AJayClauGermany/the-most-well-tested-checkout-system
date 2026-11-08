// Step definitions: the glue between the Given/When/Then sentences in spec/features/*.feature
// and the checkout code. Every sentence used by the features is listed in
// spec/features/README.md; each one appears exactly once here.

import assert from "node:assert/strict";

import { defineParameterType, Given, Then, When } from "@cucumber/cucumber";
import type { DataTable } from "@cucumber/cucumber";

import { CheckoutError, checkout, InMemoryOrderRepository, placeOrder } from "../../src/index.ts";
import { dollarsToCents } from "../support/dollars.ts";
import { loadCatalog } from "../support/spec.ts";
import type { CheckoutWorld } from "./world.ts";

// A custom parameter type: "{money}" in a step matches "$12.39" and hands the step 1239 cents.
defineParameterType({
  name: "money",
  regexp: /\$\d+(?:\.\d{1,2})?/,
  transformer: (text: string) => dollarsToCents(text),
});

// --- Given ------------------------------------------------------------------

Given("the standard catalog", function (this: CheckoutWorld) {
  this.catalog = loadCatalog();
});

Given("an empty basket", function (this: CheckoutWorld) {
  this.basket = { items: [] };
});

Given("a basket containing:", function (this: CheckoutWorld, table: DataTable) {
  for (const row of table.hashes()) {
    this.basket.items.push({ sku: row.sku as string, quantity: Number(row.quantity) });
  }
});

Given("the customer uses the coupon {string}", function (this: CheckoutWorld, code: string) {
  this.basket.coupon = code;
});

Given("the customer chooses {string} shipping", function (this: CheckoutWorld, method: string) {
  this.basket.shipping = method;
});

Given("an empty order repository", function (this: CheckoutWorld) {
  // The repository's clock reads `this.now`, so a later "the clock is fixed at" step still applies.
  this.repository = new InMemoryOrderRepository(() => this.now);
});

Given("the clock is fixed at {string}", function (this: CheckoutWorld, timestamp: string) {
  this.now = new Date(timestamp);
});

// --- When -------------------------------------------------------------------

When("the customer checks out", function (this: CheckoutWorld) {
  try {
    this.receipt = checkout(this.basket, this.catalog);
  } catch (error) {
    if (!(error instanceof CheckoutError)) {
      throw error;
    }
    this.error = error;
  }
});

When("the customer checks out and saves the order", function (this: CheckoutWorld) {
  assert.ok(this.repository, "no order repository - add 'Given an empty order repository'");
  try {
    const result = placeOrder(this.basket, this.catalog, this.repository);
    this.receipt = result.receipt;
    this.lastOrderId = result.orderId;
  } catch (error) {
    if (!(error instanceof CheckoutError)) {
      throw error;
    }
    this.error = error;
  }
});

// --- Then: the receipt ------------------------------------------------------

function receiptOf(world: CheckoutWorld) {
  assert.equal(world.error, undefined, `checkout failed with ${world.error?.code}`);
  assert.ok(world.receipt, "the customer has not checked out yet");
  return world.receipt;
}

Then("the subtotal is {money}", function (this: CheckoutWorld, cents: number) {
  assert.equal(receiptOf(this).subtotal, cents);
});

Then("the promotion discount is {money}", function (this: CheckoutWorld, cents: number) {
  assert.equal(receiptOf(this).promotion_discount, cents);
});

Then("the coupon discount is {money}", function (this: CheckoutWorld, cents: number) {
  assert.equal(receiptOf(this).coupon_discount, cents);
});

Then("the shipping cost is {money}", function (this: CheckoutWorld, cents: number) {
  assert.equal(receiptOf(this).shipping, cents);
});

Then("the tax is {money}", function (this: CheckoutWorld, cents: number) {
  assert.equal(receiptOf(this).tax, cents);
});

Then("the total is {money}", function (this: CheckoutWorld, cents: number) {
  assert.equal(receiptOf(this).total, cents);
});

Then("the receipt has {int} lines", function (this: CheckoutWorld, count: number) {
  assert.equal(receiptOf(this).lines.length, count);
});

Then(
  "the receipt line for {string} has a quantity of {int}",
  function (this: CheckoutWorld, sku: string, quantity: number) {
    const line = receiptOf(this).lines.find((l) => l.sku === sku);
    assert.ok(line, `no receipt line for ${sku}`);
    assert.equal(line.quantity, quantity);
  },
);

Then("checkout fails with the error {string}", function (this: CheckoutWorld, code: string) {
  assert.equal(this.receipt, undefined, "checkout succeeded but was expected to fail");
  assert.equal(this.error?.code, code);
});

// --- Then: saved orders -----------------------------------------------------

function repositoryOf(world: CheckoutWorld): InMemoryOrderRepository {
  assert.ok(world.repository, "no order repository - add 'Given an empty order repository'");
  return world.repository;
}

Then("the order is saved with the id {string}", function (this: CheckoutWorld, id: string) {
  assert.equal(this.lastOrderId, id);
});

Then(
  "the saved order {string} has a total of {money}",
  function (this: CheckoutWorld, id: string, cents: number) {
    const order = repositoryOf(this).find(id);
    assert.ok(order, `order ${id} not found`);
    assert.equal(order.receipt.total, cents);
  },
);

Then("the order repository is empty", function (this: CheckoutWorld) {
  assert.equal(repositoryOf(this).count(), 0);
});

Then("the order {string} cannot be found", function (this: CheckoutWorld, id: string) {
  assert.equal(repositoryOf(this).find(id), undefined);
});

Then(
  "the saved order {string} was created at {string}",
  function (this: CheckoutWorld, id: string, timestamp: string) {
    const order = repositoryOf(this).find(id);
    assert.ok(order, `order ${id} not found`);
    assert.equal(order.created_at, timestamp);
  },
);
