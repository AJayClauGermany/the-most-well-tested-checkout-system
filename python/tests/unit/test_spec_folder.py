"""Unit tests for finding the shared spec/ folder.

They use pytest's `tmp_path` fixture: a fresh, empty temporary folder for each test,
deleted again afterwards.
"""

from pathlib import Path

import pytest

from checkout import find_spec_folder


def test_finds_the_folder_that_holds_catalog_json() -> None:
    folder = find_spec_folder()

    assert folder.name == "spec"
    assert (folder / "catalog.json").is_file()
    assert (folder / "checkout-rules.md").is_file()


def test_finds_the_same_folder_when_starting_deeper_inside_the_repository() -> None:
    folder = find_spec_folder()

    assert find_spec_folder(folder / "test-cases") == folder
    assert find_spec_folder(folder / "schema" / "nested" / "deeper") == folder


def test_gives_up_with_a_clear_error_when_there_is_no_spec_folder_above_the_start(
    tmp_path: Path,
) -> None:
    with pytest.raises(FileNotFoundError, match="No spec/catalog.json found"):
        find_spec_folder(tmp_path)
    with pytest.raises(FileNotFoundError, match="No spec/catalog.json found"):
        find_spec_folder(Path("/"))


def test_is_not_fooled_by_a_spec_folder_without_a_catalog_json_in_it(tmp_path: Path) -> None:
    (tmp_path / "spec").mkdir()

    with pytest.raises(FileNotFoundError, match="No spec/catalog.json found"):
        find_spec_folder(tmp_path / "deeper")


def test_a_fake_spec_folder_with_a_catalog_is_found(tmp_path: Path) -> None:
    (tmp_path / "spec").mkdir()
    (tmp_path / "spec" / "catalog.json").write_text("{}")

    assert find_spec_folder(tmp_path / "a" / "b" / "c") == tmp_path / "spec"
