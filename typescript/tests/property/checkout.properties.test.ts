// Property-based tests with fast-check (https://fast-check.dev/).
//
// An example test says "THIS basket gives THAT total". A property says something that must
// be true for EVERY basket, and fast-check generates hundreds of random baskets trying to
// break it - then shrinks any failure to the smallest example. The properties are listed
// once for all six languages in spec/properties.md; each test is named after its id.

import fc from "fast-check";
import { describe, expect, it } from "vitest";

import {
  CheckoutError,
  checkout,
  ERROR_CODES,
  InMemoryOrderRepository,
  placeOrder,
  roundHalfUp,
} from "../../src/index.ts";
import type { Basket, BasketItem, Receipt } from "../../src/index.ts";
import { schemaErrors } from "../support/schemas.ts";
import { loadCatalog } from "../support/spec.ts";

const catalog = loadCatalog();
const SKUS = catalog.products.map((product) => product.sku);

// ---------------------------------------------------------------------------
// Generators ("arbitraries"), following the "Generating baskets" section of properties.md
// ---------------------------------------------------------------------------

/** A sensible item: a real SKU, a quantity from 1 to 99. */
const item = fc.record({
  sku: fc.constantFrom(...SKUS),
  quantity: fc.integer({ min: 1, max: 99 }),
});

/** SAVE10 or WELCOME5 written any old way: mixed case, spaces around it. */
const sloppyCoupon = fc
  .tuple(
    fc.constantFrom("SAVE10", "WELCOME5"),
    fc.array(fc.boolean(), { minLength: 8, maxLength: 8 }),
    fc.nat(3),
    fc.nat(3),
  )
  .map(([code, upper, before, after]) => {
    const letters = [...code].map((letter, i) =>
      upper[i] ? letter.toUpperCase() : letter.toLowerCase(),
    );
    return `${" ".repeat(before)}${letters.join("")}${" ".repeat(after)}`;
  });

/** A whole basket: 1 to 8 items (the same SKU may repeat), maybe a coupon, maybe a shipping method. */
const basket: fc.Arbitrary<Basket> = fc
  .record(
    {
      items: fc.array(item, { minLength: 1, maxLength: 8 }),
      coupon: fc.option(
        fc.oneof(fc.constant("SAVE10"), fc.constant("WELCOME5"), sloppyCoupon, fc.string()),
        {
          nil: undefined,
        },
      ),
      shipping: fc.option(fc.constantFrom("standard", "express", "overnight", "Express", ""), {
        nil: undefined,
      }),
    },
    { requiredKeys: ["items"] },
  )
  .map((b) => ({
    ...b,
    ...(b.coupon === undefined && Math.random() < 0.2 ? { coupon: null } : {}),
  }));

/** Anything at all that still fits the Basket type: nonsense SKUs, silly quantities, random strings. */
const anyBasket: fc.Arbitrary<Basket> = fc.record(
  {
    items: fc.array(
      fc.record({
        sku: fc.oneof(fc.constantFrom(...SKUS), fc.string()),
        quantity: fc.oneof(
          fc.integer({ min: -5, max: 150 }),
          fc.double({ noNaN: false }),
          fc.constant(Number.NaN),
        ),
      }),
      { maxLength: 6 },
    ),
    coupon: fc.option(
      fc.oneof(fc.constantFrom("SAVE10", "WELCOME5", "save10", " welcome5 "), fc.string()),
      {
        nil: undefined,
      },
    ),
    shipping: fc.option(fc.oneof(fc.constantFrom("standard", "express"), fc.string()), {
      nil: undefined,
    }),
  },
  { requiredKeys: ["items"] },
);

/** Runs checkout; `undefined` means the basket was invalid (a CheckoutError was thrown). */
function tryCheckout(b: Basket): Receipt | undefined {
  try {
    return checkout(b, catalog);
  } catch (error) {
    if (error instanceof CheckoutError) {
      return undefined;
    }
    throw error;
  }
}

/** Only baskets that check out successfully - what P1 to P9 are about. */
const validBasket = basket.filter((b) => tryCheckout(b) !== undefined);

const totalsOf = (receipt: Receipt) => ({
  subtotal: receipt.subtotal,
  promotion_discount: receipt.promotion_discount,
  coupon_discount: receipt.coupon_discount,
  shipping: receipt.shipping,
  tax: receipt.tax,
  total: receipt.total,
});

// ---------------------------------------------------------------------------
// Properties of a successful checkout
// ---------------------------------------------------------------------------

describe("properties of a successful checkout", () => {
  it("P1: total == subtotal − promotion_discount − coupon_discount + shipping + tax", () => {
    fc.assert(
      fc.property(validBasket, (b) => {
        const r = checkout(b, catalog);
        expect(r.total).toBe(
          r.subtotal - r.promotion_discount - r.coupon_discount + r.shipping + r.tax,
        );
      }),
    );
  });

  it("P2: nothing is negative and a discount never exceeds what it discounts", () => {
    fc.assert(
      fc.property(validBasket, (b) => {
        const r = checkout(b, catalog);
        for (const amount of [
          r.subtotal,
          r.promotion_discount,
          r.coupon_discount,
          r.shipping,
          r.tax,
          r.total,
        ]) {
          expect(amount).toBeGreaterThanOrEqual(0);
        }
        for (const line of r.lines) {
          expect(line.unit_price).toBeGreaterThanOrEqual(0);
          expect(line.line_total).toBeGreaterThanOrEqual(0);
          expect(line.line_discount).toBeGreaterThanOrEqual(0);
          expect(line.line_discount).toBeLessThanOrEqual(line.line_total);
        }
        expect(r.promotion_discount).toBeLessThanOrEqual(r.subtotal);
        expect(r.coupon_discount).toBeLessThanOrEqual(r.subtotal - r.promotion_discount);
      }),
    );
  });

  it("P3: the receipt lines explain the totals", () => {
    fc.assert(
      fc.property(validBasket, (b) => {
        const r = checkout(b, catalog);
        const sum = (numbers: number[]) => numbers.reduce((a, n) => a + n, 0);
        expect(r.subtotal).toBe(sum(r.lines.map((line) => line.line_total)));
        expect(r.promotion_discount).toBe(sum(r.lines.map((line) => line.line_discount)));
        for (const line of r.lines) {
          expect(line.line_total).toBe(line.unit_price * line.quantity);
        }
      }),
    );
  });

  it("P4: shuffling the items gives the same totals", () => {
    fc.assert(
      fc.property(validBasket, fc.infiniteStream(fc.nat()), (b, randoms) => {
        // A permutation of the items, driven by fast-check's random numbers (so it shrinks).
        const shuffled = [...b.items];
        for (let i = shuffled.length - 1; i > 0; i--) {
          const j = (randoms.next().value as number) % (i + 1);
          [shuffled[i], shuffled[j]] = [shuffled[j] as BasketItem, shuffled[i] as BasketItem];
        }
        expect(totalsOf(checkout({ ...b, items: shuffled }, catalog))).toEqual(
          totalsOf(checkout(b, catalog)),
        );
      }),
    );
  });

  it("P5: splitting one line into two lines of the same SKU gives the same receipt", () => {
    const splittable = validBasket.filter((b) => b.items.some((i) => i.quantity >= 2));
    fc.assert(
      fc.property(splittable, fc.nat(), fc.nat(), (b, pick, cut) => {
        const candidates = [...b.items.keys()].filter(
          (index) => (b.items[index] as BasketItem).quantity >= 2,
        );
        const index = candidates[pick % candidates.length] as number;
        const line = b.items[index] as BasketItem;
        const first = 1 + (cut % (line.quantity - 1)); // 1 .. quantity-1
        const split: BasketItem[] = [
          ...b.items.slice(0, index),
          { sku: line.sku, quantity: first },
          ...b.items.slice(index + 1),
          { sku: line.sku, quantity: line.quantity - first },
        ];
        expect(checkout({ ...b, items: split }, catalog)).toEqual(checkout(b, catalog));
      }),
    );
  });

  it("P6: the shipping table is applied exactly", () => {
    fc.assert(
      fc.property(validBasket, (b) => {
        const r = checkout(b, catalog);
        expect([0, 499, 999]).toContain(r.shipping);
        const discounted = r.subtotal - r.promotion_discount - r.coupon_discount;
        if (b.shipping === "express") {
          expect(r.shipping).toBe(999);
        } else {
          expect(r.shipping).toBe(discounted >= 5000 ? 0 : 499);
        }
      }),
    );
  });

  it("P7: tax == (discounted_subtotal × 8 + 50) ÷ 100 with integer division", () => {
    fc.assert(
      fc.property(validBasket, (b) => {
        const r = checkout(b, catalog);
        const discounted = r.subtotal - r.promotion_discount - r.coupon_discount;
        expect(r.tax).toBe(Math.floor((discounted * 8 + 50) / 100));
      }),
    );
  });

  it("P8: writing the coupon in another case or with surrounding spaces gives the same receipt", () => {
    const withRealCoupon = validBasket.filter(
      (b) => b.coupon === "SAVE10" || b.coupon === "WELCOME5",
    );
    fc.assert(
      fc.property(withRealCoupon, sloppyCoupon, (b, sloppy) => {
        // Make the sloppy version spell the SAME code as the basket's.
        const sameCode =
          sloppy.trim().toUpperCase() === b.coupon
            ? sloppy
            : sloppy.replace(/\S+/, b.coupon as string).toLowerCase();
        expect(checkout({ ...b, coupon: sameCode }, catalog)).toEqual(checkout(b, catalog));
      }),
    );
  });

  it("P9: one line per distinct SKU, in order of first appearance, each with quantity 1–99", () => {
    fc.assert(
      fc.property(validBasket, (b) => {
        const r = checkout(b, catalog);
        const distinctInOrder = [...new Set(b.items.map((i) => i.sku))];
        expect(r.lines.map((line) => line.sku)).toEqual(distinctInOrder);
        for (const line of r.lines) {
          expect(Number.isInteger(line.quantity)).toBe(true);
          expect(line.quantity).toBeGreaterThanOrEqual(1);
          expect(line.quantity).toBeLessThanOrEqual(99);
          const expectedQuantity = b.items
            .filter((i) => i.sku === line.sku)
            .reduce((sum, i) => sum + i.quantity, 0);
          expect(line.quantity).toBe(expectedQuantity);
        }
      }),
    );
  });
});

// ---------------------------------------------------------------------------
// Properties of any input at all
// ---------------------------------------------------------------------------

describe("properties of any input at all", () => {
  it("P10: for ANY basket, checkout gives a schema-valid receipt or exactly one of the six error codes - never a crash", () => {
    fc.assert(
      fc.property(anyBasket, (b) => {
        try {
          const receipt = checkout(b, catalog);
          expect(schemaErrors("receipt.schema.json", receipt)).toEqual([]);
        } catch (error) {
          expect(error).toBeInstanceOf(CheckoutError);
          expect(ERROR_CODES).toContain((error as CheckoutError).code);
        }
      }),
      { numRuns: 500 },
    );
  });

  it("P11: round_half_up is within half a cent of the exact value, and exact when the value is whole", () => {
    fc.assert(
      fc.property(fc.nat(10_000_000), fc.integer({ min: 0, max: 100 }), (amount, percent) => {
        const rounded = roundHalfUp(amount, percent);
        const exactTimes100 = amount * percent; // the exact value is this ÷ 100
        expect(Number.isInteger(rounded)).toBe(true);
        expect(Math.abs(rounded * 100 - exactTimes100)).toBeLessThanOrEqual(50);
        if (exactTimes100 % 100 === 0) {
          expect(rounded * 100).toBe(exactTimes100);
        }
        if (exactTimes100 % 100 === 50) {
          expect(rounded * 100).toBe(exactTimes100 + 50); // halves go UP
        }
      }),
    );
  });
});

// ---------------------------------------------------------------------------
// Properties of saving orders
// ---------------------------------------------------------------------------

describe("properties of saving orders", () => {
  const fixedClock = () => new Date("2026-01-15T10:30:00Z");

  it("P12: after saving n orders the ids are exactly ORDER-1 … ORDER-n, in that order", () => {
    fc.assert(
      fc.property(fc.array(validBasket, { minLength: 0, maxLength: 20 }), (baskets) => {
        const repository = new InMemoryOrderRepository(fixedClock);
        const ids = baskets.map((b) => placeOrder(b, catalog, repository).orderId);
        expect(ids).toEqual(baskets.map((_, i) => `ORDER-${i + 1}`));
      }),
    );
  });

  it("P13: find(save(receipt)) gives back the same receipt and the clock's time", () => {
    fc.assert(
      fc.property(
        validBasket,
        fc.date({ min: new Date("2000-01-01T00:00:00Z"), max: new Date("2099-12-31T23:59:59Z") }),
        (b, when) => {
          const repository = new InMemoryOrderRepository(() => when);
          const receipt = checkout(b, catalog);
          const id = repository.save(receipt);
          const order = repository.find(id);
          expect(order?.receipt).toEqual(receipt);
          expect(order?.created_at).toBe(when.toISOString().replace(/\.\d{3}Z$/, "Z"));
        },
      ),
    );
  });

  it("P14: a failed checkout saves nothing, and find of an unknown id is 'not found'", () => {
    const invalidBasket = anyBasket.filter((b) => tryCheckout(b) === undefined);
    fc.assert(
      fc.property(invalidBasket, fc.string(), (b, randomId) => {
        const repository = new InMemoryOrderRepository(fixedClock);
        expect(() => placeOrder(b, catalog, repository)).toThrow(CheckoutError);
        expect(repository.count()).toBe(0);
        expect(repository.find(randomId)).toBeUndefined();
        expect(repository.find("ORDER-1")).toBeUndefined();
      }),
    );
  });
});
