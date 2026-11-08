import type { Basket, Catalog, Receipt } from "./types.ts";

/**
 * Calculates the receipt for a basket, following spec/checkout-rules.md.
 *
 * @throws {CheckoutError} with one of the six error codes when the basket is invalid.
 */
export function checkout(basket: Basket, catalog: Catalog): Receipt {
  // Not written yet - on purpose! This repository is test-first (TDD):
  // the tests in tests/ already describe what checkout() must do.
  // Start by removing the `.skip` in tests/integration/spec-cases.test.ts,
  // watch the tests fail ("red"), then make them pass one by one ("green").
  throw new Error("checkout() is not implemented yet");
}
