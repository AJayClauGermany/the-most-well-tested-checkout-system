//! Finding the shared `spec/` folder at the root of the repository.

use std::path::{Path, PathBuf};

/// The repository's `spec/` folder, found by walking up from `start` until a folder
/// contains `spec/catalog.json`. `None` when no folder above `start` has one.
///
/// Looking upwards (instead of hard-coding `../spec`) means the program finds the spec
/// no matter where the compiled binary ends up (`target/debug/`, a copy made by a tool, ...).
pub fn find_spec_folder(start: &Path) -> Option<PathBuf> {
    let _ = start;
    todo!("find_spec_folder() is not implemented yet")
}

/// The `spec/` folder for the running program: looks above the program's own location
/// first, then above the current directory, and finally falls back to the location this
/// crate was compiled from (`CARGO_MANIFEST_DIR`, which is what `cargo run` uses).
pub fn locate_spec_folder() -> Option<PathBuf> {
    todo!("locate_spec_folder() is not implemented yet")
}

/// The first of `candidates` that has a `spec/` folder somewhere above it.
fn locate_spec_folder_among(candidates: impl IntoIterator<Item = PathBuf>) -> Option<PathBuf> {
    let _ = candidates;
    todo!("locate_spec_folder_among() is not implemented yet")
}

#[allow(dead_code)]
fn unused_until_implemented() -> Option<PathBuf> {
    locate_spec_folder_among(Vec::new())
}

// ---------------------------------------------------------------------------
// Unit tests. tempfile::tempdir() gives a fresh folder that is deleted when the test ends.
// ---------------------------------------------------------------------------
#[cfg(test)]
mod tests {
    use super::*;

    /// The `rust/` folder: where Cargo.toml is. Cargo fills this in at compile time.
    fn manifest_dir() -> PathBuf {
        PathBuf::from(env!("CARGO_MANIFEST_DIR"))
    }

    #[test]
    fn finds_the_folder_that_holds_catalog_json() {
        let folder = find_spec_folder(&manifest_dir()).expect("the spec folder should be found");

        assert!(folder.ends_with("spec"));
        assert!(folder.join("catalog.json").is_file());
        assert!(folder.join("checkout-rules.md").is_file());
    }

    #[test]
    fn finds_the_same_folder_when_starting_deeper_inside_the_repository() {
        let folder = find_spec_folder(&manifest_dir()).unwrap();

        assert_eq!(
            find_spec_folder(&folder.join("test-cases")),
            Some(folder.clone())
        );
        assert_eq!(
            find_spec_folder(&folder.join("schema").join("nested").join("deeper")),
            Some(folder)
        );
    }

    #[test]
    fn gives_none_when_there_is_no_spec_folder_above_the_start() {
        let temporary = tempfile::tempdir().unwrap();

        assert_eq!(find_spec_folder(temporary.path()), None);
        assert_eq!(find_spec_folder(Path::new("/")), None);
    }

    #[test]
    fn is_not_fooled_by_a_spec_folder_without_a_catalog_json_in_it() {
        let temporary = tempfile::tempdir().unwrap();
        std::fs::create_dir(temporary.path().join("spec")).unwrap();

        assert_eq!(find_spec_folder(&temporary.path().join("deeper")), None);
    }

    #[test]
    fn locate_tries_the_candidates_in_order_and_takes_the_first_hit() {
        let temporary = tempfile::tempdir().unwrap();
        let expected = find_spec_folder(&manifest_dir()).unwrap();

        let found = locate_spec_folder_among([temporary.path().to_path_buf(), manifest_dir()]);

        assert_eq!(found, Some(expected));
    }

    #[test]
    fn locate_gives_none_when_no_candidate_has_a_spec_folder_above_it() {
        let temporary = tempfile::tempdir().unwrap();

        assert_eq!(
            locate_spec_folder_among([temporary.path().to_path_buf()]),
            None
        );
        assert_eq!(locate_spec_folder_among([]), None);
    }

    #[test]
    fn locate_spec_folder_finds_the_spec_from_inside_the_repository() {
        // The test binary lives in rust/target/..., and `cargo test` runs in rust/: both are
        // inside the repository, and so is the compile-time fallback.
        let folder = locate_spec_folder().expect("the spec folder should be found");

        assert_eq!(folder, find_spec_folder(&manifest_dir()).unwrap());
    }
}
