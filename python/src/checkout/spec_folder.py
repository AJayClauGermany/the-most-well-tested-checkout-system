"""Finding the shared spec/ folder at the root of the repository."""

from pathlib import Path


def find_spec_folder(start: Path | None = None) -> Path:
    """The repository's ``spec/`` folder.

    Walks up from ``start`` (by default, this very file) until it reaches a folder that
    contains ``spec/catalog.json``. Looking upwards, instead of hard-coding ``../../spec``,
    means the code still finds the spec when it runs from a copy somewhere deeper inside
    the repository - which is exactly what a mutation testing tool does with it.

    Raises ``FileNotFoundError`` when no folder above ``start`` contains ``spec/catalog.json``.
    """
    here = (start if start is not None else Path(__file__)).resolve()
    # `here` itself first, then each parent up to the root of the disk.
    for folder in (here, *here.parents):
        if (folder / "spec" / "catalog.json").is_file():
            return folder / "spec"
    raise FileNotFoundError(f"No spec/catalog.json found in any folder above {here}")
