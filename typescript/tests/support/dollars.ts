// Test helper: turns "$12.39" into 1239 cents - without ever touching a floating point
// number (spec/features/README.md insists, and rules §1 explains why).

/** "$12.39" or "12.39" → 1239. "$5" → 500. */
export function dollarsToCents(text: string): number {
  const trimmed = text.trim().replace(/^\$/, "");
  if (!/^\d+(\.\d{1,2})?$/.test(trimmed)) {
    throw new Error(`not a dollar amount: ${JSON.stringify(text)}`);
  }
  const [dollars = "0", cents = "0"] = trimmed.split(".");
  return Number(dollars) * 100 + Number(cents.padEnd(2, "0"));
}
