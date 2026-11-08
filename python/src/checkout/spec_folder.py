"""Finding the shared spec/ folder at the root of the repository."""

from pathlib import Path


def find_spec_folder(start: Path | None = None) -> Path:
    """The repository's ``spec/`` folder, found by walking up from ``start`` (default: this file).

    Not implemented yet - see tests/unit/test_spec_folder.py.
    """
    raise NotImplementedError(f"find_spec_folder({start}) is not implemented yet")
