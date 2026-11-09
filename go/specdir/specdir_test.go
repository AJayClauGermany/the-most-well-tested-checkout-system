package specdir

// UNIT TESTS for finding the shared spec/ folder. They use real folders (the repository,
// and empty temporary folders from t.TempDir(), which Go deletes when the test ends).

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// fileExists reports whether path exists.
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func TestFindFindsTheFolderThatHoldsCatalogJSON(t *testing.T) {
	t.Parallel()
	wd, _ := os.Getwd() // go test runs in this package's folder: go/specdir

	dir, err := Find(wd)

	if err != nil {
		t.Fatalf("Find() returned an error: %v", err)
	}
	if filepath.Base(dir) != "spec" {
		t.Errorf("Find() = %s, want a folder named spec", dir)
	}
	for _, name := range []string{"catalog.json", "checkout-rules.md"} {
		if !fileExists(filepath.Join(dir, name)) {
			t.Errorf("%s has no %s", dir, name)
		}
	}
}

func TestFindFindsTheSameFolderWhenStartingDeeperInside(t *testing.T) {
	t.Parallel()
	wd, _ := os.Getwd()
	dir, _ := Find(wd)

	for _, start := range []string{
		filepath.Join(dir, "test-cases"),
		filepath.Join(dir, "schema", "nested", "deeper"), // need not even exist
	} {
		got, err := Find(start)
		if err != nil || got != dir {
			t.Errorf("Find(%s) = %s, %v; want %s", start, got, err, dir)
		}
	}
}

func TestFindGivesUpWithAClearErrorWhenThereIsNoSpecFolderAbove(t *testing.T) {
	t.Parallel()
	start := t.TempDir()

	_, err := Find(start)

	if err == nil || !strings.Contains(err.Error(), "no spec/catalog.json found in any folder above "+start) {
		t.Errorf("Find(%s) error = %v, want a clear 'not found'", start, err)
	}
}

func TestFindIsNotFooledByASpecFolderWithoutACatalog(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, "spec"), 0o755); err != nil {
		t.Fatal(err)
	}

	_, err := Find(filepath.Join(root, "deeper"))

	if err == nil {
		t.Error("Find() found a spec folder without a catalog.json")
	}
}

func TestLocateFindsTheSpecFolderFromTheWorkingDirectory(t *testing.T) {
	t.Parallel()
	wd, _ := os.Getwd()
	want, _ := Find(wd)

	got, err := Locate()

	if err != nil || got != want {
		t.Errorf("Locate() = %s, %v; want %s", got, err, want)
	}
}

// Not parallel: it changes the working directory of the whole test process.
func TestLocateFallsBackToTheProgramsOwnFolder(t *testing.T) {
	// From a folder with no spec/ above it, Locate looks next to the executable instead.
	// The test binary is built in a temporary folder outside the repository, so here that
	// fails too – and the error must say where it looked last.
	t.Chdir(t.TempDir())
	exe, _ := os.Executable()

	_, err := Locate()

	if err == nil || !strings.Contains(err.Error(), filepath.Dir(exe)) {
		t.Errorf("Locate() error = %v, want it to mention %s", err, filepath.Dir(exe))
	}
}
