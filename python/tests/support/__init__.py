"""Helpers for tests, not tests themselves.

* ``fixtures``  - a tiny made-up catalog and a sample receipt
* ``dollars``   - "$12.39" -> 1239, without floating point
* ``schemas``   - schema_errors(schema, value) with jsonschema (draft 2020-12)

pytest finds this package because pyproject.toml has ``pythonpath = ["tests"]``.
"""
