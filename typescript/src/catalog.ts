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
