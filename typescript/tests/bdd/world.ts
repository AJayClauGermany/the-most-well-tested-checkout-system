// The Cucumber "World": one object that lives for one scenario and carries its state
// from step to step (the basket being built, the receipt or error that came out, ...).

import { setWorldConstructor, World } from "@cucumber/cucumber";
import type { IWorldOptions } from "@cucumber/cucumber";

import type { InMemoryOrderRepository } from "../../src/index.ts";
import type { Basket, Catalog, CheckoutError, Receipt } from "../../src/index.ts";

export class CheckoutWorld extends World {
  catalog: Catalog = { products: [] };
  basket: Basket = { items: [] };
  receipt: Receipt | undefined;
  error: CheckoutError | undefined;

  /** What the (replaceable) clock answers. "Given the clock is fixed at ..." changes it. */
  now = new Date("2026-01-01T00:00:00Z");
  repository: InMemoryOrderRepository | undefined;
  lastOrderId: string | undefined;

  constructor(options: IWorldOptions) {
    super(options);
  }
}

setWorldConstructor(CheckoutWorld);
