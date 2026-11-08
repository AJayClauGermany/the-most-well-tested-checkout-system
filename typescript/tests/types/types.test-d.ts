// Type-level tests: checks that run in the TYPE CHECKER, not at runtime.
//
// A normal test asks "does this value equal that?". A type test asks "does TypeScript
// agree that this type is that type?". They catch a public API quietly changing shape.
// Vitest runs files ending in .test-d.ts with `vitest --typecheck` (npm run test:types);
// expectTypeOf(...) assertions are checked by tsc and never execute.

import { describe, expectTypeOf, it } from "vitest";

import {
  checkout,
  ERROR_CODES,
  type InMemoryOrderRepository,
  placeOrder,
} from "../../src/index.ts";
import type {
  Basket,
  BasketItem,
  Catalog,
  CheckoutError,
  Clock,
  ErrorCode,
  OrderRepository,
  Promotion,
  Receipt,
  ReceiptLine,
} from "../../src/index.ts";

describe("the checkout types", () => {
  it("a promotion is one of exactly two kinds, told apart by `type`", () => {
    expectTypeOf<Promotion["type"]>().toEqualTypeOf<"buy_x_get_y_free" | "bulk_discount">();
  });

  it("the basket's shipping is any string (so invalid values can be tested), the coupon may be null", () => {
    expectTypeOf<Basket["shipping"]>().toEqualTypeOf<string | undefined>();
    expectTypeOf<Basket["coupon"]>().toEqualTypeOf<string | null | undefined>();
    expectTypeOf<Basket["items"]>().toEqualTypeOf<BasketItem[]>();
  });

  it("a receipt is made of numbers and lines, nothing optional", () => {
    expectTypeOf<Receipt["total"]>().toBeNumber();
    expectTypeOf<Receipt["lines"]>().toEqualTypeOf<ReceiptLine[]>();
    expectTypeOf<keyof Receipt>().toEqualTypeOf<
      "lines" | "subtotal" | "promotion_discount" | "coupon_discount" | "shipping" | "tax" | "total"
    >();
  });

  it("there are exactly six error codes, and the list is read-only", () => {
    expectTypeOf<ErrorCode>().toEqualTypeOf<
      | "EMPTY_BASKET"
      | "UNKNOWN_PRODUCT"
      | "INVALID_QUANTITY"
      | "UNKNOWN_SHIPPING_METHOD"
      | "UNKNOWN_COUPON"
      | "COUPON_MINIMUM_NOT_MET"
    >();
    expectTypeOf(ERROR_CODES).toEqualTypeOf<readonly ErrorCode[]>();
    expectTypeOf<CheckoutError["code"]>().toEqualTypeOf<ErrorCode>();
  });
});

describe("the functions", () => {
  it("checkout takes a basket and a catalog and gives a receipt (never undefined)", () => {
    expectTypeOf(checkout).parameters.toEqualTypeOf<[Basket, Catalog]>();
    expectTypeOf(checkout).returns.toEqualTypeOf<Receipt>();
  });

  it("placeOrder works with ANY repository, not just the in-memory one", () => {
    expectTypeOf(placeOrder).parameter(2).toEqualTypeOf<OrderRepository>();
    expectTypeOf<InMemoryOrderRepository>().toMatchTypeOf<OrderRepository>();
  });

  it("a clock is just a function that gives a Date", () => {
    expectTypeOf<Clock>().toEqualTypeOf<() => Date>();
    expectTypeOf<ReturnType<Clock>>().toEqualTypeOf<Date>();
  });

  it("a receipt has no field named after a typo", () => {
    // @ts-expect-error - 'totl' does not exist; if this line ever compiles, the test fails.
    expectTypeOf<Receipt["totl"]>().toBeNumber();
  });
});
