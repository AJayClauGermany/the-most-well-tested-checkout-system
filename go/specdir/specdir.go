// Package specdir finds the repository's shared spec/ folder by walking up the folder
// tree until a folder contains spec/catalog.json.
package specdir

import "errors"

var errNotImplemented = errors.New("specdir: not implemented yet")

// Find walks up from start until it reaches a folder that contains spec/catalog.json,
// and returns that spec folder. TODO: not implemented yet.
func Find(start string) (string, error) {
	return "", errNotImplemented
}

// Locate finds the spec folder for a running program, from the working directory or
// the executable's folder. TODO: not implemented yet.
func Locate() (string, error) {
	return "", errNotImplemented
}
