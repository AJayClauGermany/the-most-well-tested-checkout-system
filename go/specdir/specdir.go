// Package specdir finds the repository's shared spec/ folder.
//
// The spec files (catalog, schemas, test cases) live once, at the root of the repository,
// and are never copied into a language folder. Code and tests that need them look upwards
// from where they are until they find a folder that contains spec/catalog.json. Looking
// upwards (instead of hard-coding "../spec") means the same code works from go/, from
// go/checkout and from the repository root.
//
// A copy of go/ somewhere else – which is what the mutation testing tool makes, in a
// temporary folder – has no spec/ above it. For that case the environment variable
// CHECKOUT_SPEC_DIR can name the spec folder directly.
package specdir

import (
	"fmt"
	"os"
	"path/filepath"
)

// EnvVar is the environment variable that names the spec folder directly, for programs
// and tests that run outside the repository. When it is unset or empty, Locate looks up the
// folder tree instead.
const EnvVar = "CHECKOUT_SPEC_DIR"

// Find walks up from start, one parent folder at a time, until it reaches a folder that
// contains spec/catalog.json, and returns that spec folder.
func Find(start string) (string, error) {
	dir := start
	for {
		candidate := filepath.Join(dir, "spec")
		if _, err := os.Stat(filepath.Join(candidate, "catalog.json")); err == nil {
			return candidate, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir { // reached the root of the file system
			return "", fmt.Errorf("no spec/catalog.json found in any folder above %s", start)
		}
		dir = parent
	}
}

// Locate finds the spec folder for a running program: the folder named by CHECKOUT_SPEC_DIR
// if it is set; otherwise by looking up from the current working directory (this is what
// makes `go run ./cmd/checkout` work from go/ or from the repository root), then from the
// folder the executable itself is in (for a built binary started from anywhere).
func Locate() (string, error) {
	if dir := os.Getenv(EnvVar); dir != "" {
		if _, err := os.Stat(filepath.Join(dir, "catalog.json")); err != nil {
			return "", fmt.Errorf("%s=%s has no catalog.json", EnvVar, dir)
		}
		return dir, nil
	}
	if wd, err := os.Getwd(); err == nil {
		if dir, err := Find(wd); err == nil {
			return dir, nil
		}
	}
	// If the program cannot tell where it is (very rare), exe is "" and Find below simply
	// reports that it found nothing.
	exe, _ := os.Executable()
	return Find(programFolder(exe))
}

// programFolder is the folder the program at path really lives in: when path is a symlink
// (say, a link to the built program in a folder on the PATH), the folder of the file it
// points to. A path that cannot be resolved is used as it is.
func programFolder(path string) string {
	if resolved, err := filepath.EvalSymlinks(path); err == nil {
		path = resolved
	}
	return filepath.Dir(path)
}
