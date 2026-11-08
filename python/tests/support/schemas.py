"""Test helper: JSON Schema validation with jsonschema (https://python-jsonschema.readthedocs.io/).

The schemas in spec/schema/ are the *contract* every language's output must honour.
A "contract test" checks a real receipt (or saved order, or CLI error) against them.
The schemas use JSON Schema draft 2020-12, so we use ``Draft202012Validator``.
"""

import json
from typing import Any, Literal

from jsonschema import Draft202012Validator
from referencing import Registry, Resource
from spec_loader import SPEC_FOLDER

SchemaName = Literal[
    "receipt.schema.json",
    "order.schema.json",
    "cli-error.schema.json",
    "basket.schema.json",
]

SCHEMA_NAMES: tuple[SchemaName, ...] = (
    "receipt.schema.json",
    "order.schema.json",
    "cli-error.schema.json",
    "basket.schema.json",
)


def _read_schema(name: str) -> dict[str, Any]:
    data: dict[str, Any] = json.loads((SPEC_FOLDER / "schema" / name).read_text(encoding="utf-8"))
    return data


_SCHEMAS: dict[str, dict[str, Any]] = {name: _read_schema(name) for name in SCHEMA_NAMES}

# One registry knows every schema by its file name, so that order.schema.json can
# `"$ref": "receipt.schema.json"` and the validator finds it without touching the network.
_REGISTRY: Registry[Any] = Registry().with_resources(
    (name, Resource.from_contents(schema)) for name, schema in _SCHEMAS.items()
)


def schema_errors(schema: SchemaName, value: object) -> list[str]:
    """Validates ``value`` against one of the shared schemas.

    Returns an empty list when the value is valid, otherwise one message per problem,
    sorted, as ``"<path> <message>"`` (an empty list is easy to assert with ``== []``
    and the messages print nicely when it fails).
    """
    validator = Draft202012Validator(_SCHEMAS[schema], registry=_REGISTRY)
    errors = []
    for error in validator.iter_errors(value):
        path = "/" + "/".join(str(part) for part in error.absolute_path)
        errors.append(f"{path} {error.message}")
    return sorted(errors)
