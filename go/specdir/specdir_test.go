package specdir

// UNIT TESTS for finding the shared spec/ folder. They use real folders (the repository,
// and empty temporary folders from t.TempDir(), which Go deletes when the test ends).

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// fakeRepository builds a tiny stand-in for the repository in a temporary folder:
// <root>/spec/catalog.json and an empty <root>/go/specdir. The tests below use it instead
// of the real repository, so they pass wherever the code is copied to – a mutation testing
// tool, for example, runs the tests in a copy of go/ far away from the real spec/ folder.
func fakeRepository(t *testing.T) (root, spec string) {
	t.Helper()
	root = t.TempDir()
	spec = filepath.Join(root, "spec")
	for _, dir := range []string{spec, filepath.Join(root, "go", "specdir")} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(spec, "catalog.json"), []byte(`{"products":[]}`), 0o644); err != nil {
		t.Fatal(err)
	}
	return root, spec
}

func TestFindFindsTheFolderThatHoldsCatalogJSON(t *testing.T) {
	t.Parallel()
	root, spec := fakeRepository(t)

	dir, err := Find(filepath.Join(root, "go", "specdir"))

	if err != nil || dir != spec {
		t.Errorf("Find() = %s, %v; want %s", dir, err, spec)
	}
}

func TestFindFindsTheSameFolderWhenStartingDeeperInside(t *testing.T) {
	t.Parallel()
	_, spec := fakeRepository(t)

	for _, start := range []string{
		spec,
		filepath.Join(spec, "test-cases"),
		filepath.Join(spec, "schema", "nested", "deeper"), // need not even exist
	} {
		got, err := Find(start)
		if err != nil || got != spec {
			t.Errorf("Find(%s) = %s, %v; want %s", start, got, err, spec)
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

// Not parallel: t.Setenv and t.Chdir change the whole test process.
func TestLocateFindsTheSpecFolderFromTheWorkingDirectory(t *testing.T) {
	root, spec := fakeRepository(t)
	t.Setenv(EnvVar, "")
	t.Chdir(filepath.Join(root, "go"))

	got, err := Locate()

	if err != nil || got != spec {
		t.Errorf("Locate() = %s, %v; want %s", got, err, spec)
	}
}

// Not parallel: t.Setenv and t.Chdir change the whole test process.
func TestLocateUsesTheFolderNamedByTheEnvironmentVariableFirst(t *testing.T) {
	_, spec := fakeRepository(t)
	t.Setenv(EnvVar, spec)
	t.Chdir(t.TempDir()) // no spec/ above here: the variable alone must be enough

	got, err := Locate()

	if err != nil || got != spec {
		t.Errorf("Locate() = %s, %v; want %s", got, err, spec)
	}
}

// Not parallel: t.Setenv changes the whole test process.
func TestLocateRejectsAnEnvironmentVariableThatPointsAtAFolderWithoutACatalog(t *testing.T) {
	empty := t.TempDir()
	t.Setenv(EnvVar, empty)

	_, err := Locate()

	want := EnvVar + "=" + empty + " has no catalog.json"
	if err == nil || err.Error() != want {
		t.Errorf("Locate() error = %v, want %q", err, want)
	}
}

func TestTheEnvironmentVariableIsCalledCheckoutSpecDir(t *testing.T) {
	t.Parallel()
	if EnvVar != "CHECKOUT_SPEC_DIR" {
		t.Errorf("EnvVar = %q, want CHECKOUT_SPEC_DIR (the READMEs and CI use that name)", EnvVar)
	}
}

func TestProgramFolderFollowsASymlinkToTheRealProgram(t *testing.T) {
	t.Parallel()
	// Someone may link the built program into a folder on their PATH. The spec folder must
	// then be looked for next to the real program, not next to the link.
	root, _ := fakeRepository(t)
	program := filepath.Join(root, "go", "checkout")
	if err := os.WriteFile(program, []byte("a program"), 0o755); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(t.TempDir(), "checkout")
	if err := os.Symlink(program, link); err != nil {
		t.Skipf("cannot create a symlink here: %v", err) // e.g. Windows without developer mode
	}
	want, _ := filepath.EvalSymlinks(filepath.Dir(program)) // t.TempDir may itself sit behind a link

	if got := programFolder(link); got != want {
		t.Errorf("programFolder(%s) = %s, want %s", link, got, want)
	}
}

func TestProgramFolderOfAPlainFileIsTheFolderItIsIn(t *testing.T) {
	t.Parallel()
	// If the path cannot be resolved (here: it does not exist), it is used as it is.
	missing := filepath.Join(t.TempDir(), "no-such-program")

	if got := programFolder(missing); got != filepath.Dir(missing) {
		t.Errorf("programFolder(%s) = %s, want %s", missing, got, filepath.Dir(missing))
	}
}

// Not parallel: it changes the working directory of the whole test process.
func TestLocateFallsBackToTheProgramsOwnFolder(t *testing.T) {
	// From a folder with no spec/ above it, Locate looks next to the executable instead.
	// The test binary is built in a temporary folder outside the repository, so here that
	// fails too – and the error must say where it looked last.
	t.Setenv(EnvVar, "")
	t.Chdir(t.TempDir())
	exe, _ := os.Executable()

	_, err := Locate()

	if err == nil || !strings.Contains(err.Error(), filepath.Dir(exe)) {
		t.Errorf("Locate() error = %v, want it to mention %s", err, filepath.Dir(exe))
	}
}
