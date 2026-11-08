// Unit tests for saving orders (rules §6), and the project's showcase of TEST DOUBLES.
//
// A *test double* stands in for a real dependency:
//   * a FAKE is a simple working version - InMemoryOrderRepository is one,
//   * a STUB gives canned answers - a clock that always says the same time,
//   * a MOCK records how it was called - vi.fn() lets us check that save() was never called.
// Vitest can also fake the system clock itself (vi.useFakeTimers / vi.setSystemTime).

import { afterEach, describe, expect, it, vi } from "vitest";

import {
  CheckoutError,
  checkout,
  formatTimestamp,
  InMemoryOrderRepository,
  orderId,
  orderNumber,
  placeOrder,
  systemClock,
} from "../../src/index.ts";
import type { OrderRepository } from "../../src/index.ts";
import { oneAppleReceipt, tinyCatalog } from "../support/fixtures.ts";

/** A stub clock: always 10:30 on 15 January 2026, UTC. */
const fixedClock = () => new Date("2026-01-15T10:30:00Z");

describe("formatTimestamp", () => {
  it("writes UTC with whole seconds and a trailing Z", () => {
    expect(formatTimestamp(new Date("2026-01-15T10:30:00Z"))).toBe("2026-01-15T10:30:00Z");
  });

  it("drops the milliseconds instead of rounding them", () => {
    expect(formatTimestamp(new Date("2026-01-15T10:30:00.999Z"))).toBe("2026-01-15T10:30:00Z");
  });

  it("converts other time zones to UTC", () => {
    expect(formatTimestamp(new Date("2026-01-15T12:30:00+02:00"))).toBe("2026-01-15T10:30:00Z");
  });
});

describe("orderId and orderNumber", () => {
  it.each([
    [1, "ORDER-1"],
    [2, "ORDER-2"],
    [10, "ORDER-10"],
    [1234, "ORDER-1234"],
  ])("number %i ↔ id %s", (n, id) => {
    expect(orderId(n)).toBe(id);
    expect(orderNumber(id)).toBe(n);
  });

  it.each([
    "",
    "ORDER-",
    "ORDER-0",
    "ORDER-01",
    "ORDER--1",
    "ORDER-1.5",
    "ORDER-1 ",
    "order-1",
    "ORDER-x",
    "1",
    "ORDER_1",
  ])("does not recognise %j as an order id", (id) => {
    expect(orderNumber(id)).toBeUndefined();
  });
});

describe("systemClock", () => {
  afterEach(() => {
    vi.useRealTimers();
  });

  it("reads the real clock - which Vitest can fake", () => {
    vi.useFakeTimers();
    vi.setSystemTime(new Date("2030-06-01T08:00:00Z"));

    expect(systemClock().toISOString()).toBe("2030-06-01T08:00:00.000Z");
  });
});

describe("InMemoryOrderRepository (a hand-written fake)", () => {
  it("starts empty", () => {
    const repository = new InMemoryOrderRepository(fixedClock);

    expect(repository.count()).toBe(0);
    expect(repository.find("ORDER-1")).toBeUndefined();
  });

  it("numbers orders ORDER-1, ORDER-2, ORDER-3 in the order they are saved", () => {
    const repository = new InMemoryOrderRepository(fixedClock);

    const ids = [
      repository.save(oneAppleReceipt),
      repository.save(oneAppleReceipt),
      repository.save(oneAppleReceipt),
    ];

    expect(ids).toEqual(["ORDER-1", "ORDER-2", "ORDER-3"]);
    expect(repository.count()).toBe(3);
  });

  it("finds a saved order with its receipt and the time the clock said", () => {
    const repository = new InMemoryOrderRepository(fixedClock);

    const id = repository.save(oneAppleReceipt);

    expect(repository.find(id)).toEqual({
      id: "ORDER-1",
      created_at: "2026-01-15T10:30:00Z",
      receipt: oneAppleReceipt,
    });
  });

  it("asks the clock at the moment of saving, once per save", () => {
    // A mock clock: it answers like the stub, and remembers it was called.
    const clock = vi.fn(fixedClock);
    const repository = new InMemoryOrderRepository(clock);

    repository.save(oneAppleReceipt);
    repository.save(oneAppleReceipt);

    expect(clock).toHaveBeenCalledTimes(2);
  });

  it("keeps a copy of the receipt, so changing the original later does not change the order", () => {
    const repository = new InMemoryOrderRepository(fixedClock);
    const receipt = structuredClone(oneAppleReceipt);

    const id = repository.save(receipt);
    receipt.total = 1;
    receipt.lines[0]!.quantity = 42;

    expect(repository.find(id)?.receipt).toEqual(oneAppleReceipt);
  });

  it("answers 'not found' (undefined) for ids that were never handed out", () => {
    const repository = new InMemoryOrderRepository(fixedClock);
    repository.save(oneAppleReceipt);

    expect(repository.find("ORDER-2")).toBeUndefined();
    expect(repository.find("ORDER-0")).toBeUndefined();
    expect(repository.find("ORDER-99")).toBeUndefined();
    expect(repository.find("order-1")).toBeUndefined();
    expect(repository.find("")).toBeUndefined();
  });

  it("uses the real clock when none is given", () => {
    vi.useFakeTimers();
    vi.setSystemTime(new Date("2031-02-03T04:05:06.789Z"));
    const repository = new InMemoryOrderRepository();

    const id = repository.save(oneAppleReceipt);

    expect(repository.find(id)?.created_at).toBe("2031-02-03T04:05:06Z");
    vi.useRealTimers();
  });
});

describe("placeOrder", () => {
  it("checks out and saves the receipt, returning both the receipt and the new id", () => {
    const repository = new InMemoryOrderRepository(fixedClock);

    const result = placeOrder({ items: [{ sku: "APPLE", quantity: 1 }] }, tinyCatalog, repository);

    expect(result.orderId).toBe("ORDER-1");
    expect(result.receipt).toEqual(
      checkout({ items: [{ sku: "APPLE", quantity: 1 }] }, tinyCatalog),
    );
    expect(repository.find("ORDER-1")?.receipt).toEqual(result.receipt);
  });

  it("saves exactly what it returns - checked with a mock repository", () => {
    // vi.fn() creates a MOCK: a function that records every call. We build a repository
    // out of two of them, so the test can see exactly how placeOrder used it.
    const repository: OrderRepository = {
      save: vi.fn(() => "ORDER-7"),
      find: vi.fn(() => undefined),
    };

    const result = placeOrder({ items: [{ sku: "TV", quantity: 1 }] }, tinyCatalog, repository);

    expect(repository.save).toHaveBeenCalledTimes(1);
    expect(repository.save).toHaveBeenCalledWith(result.receipt);
    expect(result.orderId).toBe("ORDER-7"); // the id comes from the repository, not from placeOrder
    expect(repository.find).not.toHaveBeenCalled();
  });

  it("saves NOTHING when checkout fails (P14)", () => {
    const repository: OrderRepository = {
      save: vi.fn(() => "ORDER-1"),
      find: vi.fn(() => undefined),
    };

    expect(() => placeOrder({ items: [] }, tinyCatalog, repository)).toThrow(CheckoutError);
    expect(() =>
      placeOrder({ items: [{ sku: "BANANA", quantity: 1 }] }, tinyCatalog, repository),
    ).toThrow(expect.objectContaining({ code: "UNKNOWN_PRODUCT" }));

    expect(repository.save).not.toHaveBeenCalled();
  });

  it("leaves the fake repository empty after a failed checkout", () => {
    const repository = new InMemoryOrderRepository(fixedClock);

    expect(() =>
      placeOrder({ items: [{ sku: "APPLE", quantity: 0 }] }, tinyCatalog, repository),
    ).toThrow(CheckoutError);

    expect(repository.count()).toBe(0);
    expect(repository.find("ORDER-1")).toBeUndefined();
  });
});
