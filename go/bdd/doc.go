// Package bdd runs the Given/When/Then scenarios in spec/features/*.feature with godog,
// the Cucumber implementation for Go.
//
// There is no production code here, only tests: features_test.go points godog at the
// shared feature files (read in place, never copied) and steps_test.go holds one step
// definition for every sentence listed in spec/features/README.md.
package bdd
