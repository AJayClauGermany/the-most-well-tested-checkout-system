// Package specdir finds the repository's shared spec/ folder.
//
// The spec files (catalog, schemas, test cases) live once, at the root of the repository,
// and are never copied into a language folder. Code and tests that need them look upwards
// from where they are until they find a folder that contains spec/catalog.json. Looking
// upwards (instead of hard-coding "../spec") means the same code works from go/, from
// go/checkout, from the repository root – and from a copy of the folder somewhere deeper,
// which is what a mutation testing tool makes.
package specdir

import (
	"fmt"
	"os"
	"path/filepath"
)

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

// Locate finds the spec folder for a running program: first by looking up from the current
// working directory (this is what makes `go run ./cmd/checkout` work from go/ or from the
// repository root), then from the folder the executable itself is in (for a built binary
// started from anywhere).
func Locate() (string, error) {
	if wd, err := os.Getwd(); err == nil {
		if dir, err := Find(wd); err == nil {
			return dir, nil
		}
	}
	// If the program cannot tell where it is (very rare), exe is "" and Find below simply
	// reports that it found nothing.
	exe, _ := os.Executable()
	if resolved, err := filepath.EvalSymlinks(exe); err == nil {
		exe = resolved
	}
	return Find(filepath.Dir(exe))
}
