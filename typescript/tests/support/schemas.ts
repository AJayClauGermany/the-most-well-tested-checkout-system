// Test helper: JSON Schema validation with ajv (https://ajv.js.org/).
//
// The schemas in spec/schema/ are the *contract* every language's output must honour.
// A "contract test" checks a real receipt (or saved order, or CLI error) against them.
// The schemas use JSON Schema draft 2020-12, so we use ajv's 2020 class.

import { readFileSync } from "node:fs";
import { join } from "node:path";

import { Ajv2020 } from "ajv/dist/2020.js";
import type { ValidateFunction } from "ajv";

import { SPEC_DIR } from "./spec.ts";

function readSchema(name: string): object {
  return JSON.parse(readFileSync(join(SPEC_DIR, "schema", name), "utf8")) as object;
}

// One ajv instance knows every schema, so `order.schema.json` can `$ref` `receipt.schema.json`
// by that exact name.
const ajv = new Ajv2020({ allErrors: true, strict: true });
ajv.addSchema(readSchema("receipt.schema.json"), "receipt.schema.json");
ajv.addSchema(readSchema("order.schema.json"), "order.schema.json");
ajv.addSchema(readSchema("cli-error.schema.json"), "cli-error.schema.json");
ajv.addSchema(readSchema("basket.schema.json"), "basket.schema.json");

function compiled(name: string): ValidateFunction {
  const validate = ajv.getSchema(name);
  if (validate === undefined) {
    throw new Error(`schema ${name} was not added`);
  }
  return validate;
}

/**
 * Validates `value` against one of the shared schemas.
 *
 * @returns an empty list when the value is valid, otherwise one message per problem
 * (an empty list is easy to assert with `toEqual([])` and prints nicely when it fails).
 */
export function schemaErrors(
  schema:
    | "receipt.schema.json"
    | "order.schema.json"
    | "cli-error.schema.json"
    | "basket.schema.json",
  value: unknown,
): string[] {
  const validate = compiled(schema);
  validate(value);
  return (validate.errors ?? []).map(
    (error) => `${error.instancePath || "/"} ${error.message ?? ""}`,
  );
}
