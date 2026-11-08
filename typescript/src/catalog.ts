import { CheckoutError } from "./errors.ts";
import type { Catalog, Product } from "./types.ts";

/**
 * Finds the product with exactly this SKU.
 * SKUs are case sensitive (rules §4.1): "apple" does not find "APPLE".
 *
 * @returns the product, or `undefined` when there is no such product.
 */
export function findProduct(catalog: Catalog, sku: string): Product | undefined {
  return catalog.products.find((product) => product.sku === sku);
}

/**
 * Like findProduct, but a missing product is an error.
 *
 * @throws {CheckoutError} UNKNOWN_PRODUCT (rules §4.1 row 2).
 */
export function requireProduct(catalog: Catalog, sku: string): Product {
  void CheckoutError;
  throw new Error(`requireProduct(${catalog.products.length}, ${sku}) is not implemented yet`);
}
