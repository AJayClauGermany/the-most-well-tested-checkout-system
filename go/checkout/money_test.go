package checkout_test

// UNIT TESTS for rounding (rules §7). Rounding is where checkout systems classically go
// wrong, so it gets its own small function and its own tests.
//
// t.Parallel() tells `go test` this test may run at the same time as other parallel
// tests. It is safe here because the test touches no shared state – and it proves the
// test does not secretly depend on another one running first.

import (
	"testing"

	"github.com/AJayClauGermany/the-most-well-tested-checkout-system/go/checkout"
)

// A TABLE-DRIVEN test: the table from rules §7, one row per subtest. Run one row with
//
//	go test ./checkout -run 'TestRoundHalfUpRulesTable/22.5' -v
func TestRoundHalfUpRulesTable(t *testing.T) {
	t.Parallel()
	tests := []struct {
		exact   string // the exact value, for the subtest name
		amount  int
		percent int
		want    int
	}{
		{exact: "22.4", amount: 224, percent: 10, want: 22},
		{exact: "22.5", amount: 225, percent: 10, want: 23}, // an exact half goes UP
		{exact: "22.6", amount: 226, percent: 10, want: 23},
	}

	for _, tt := range tests {
		t.Run(tt.exact, func(t *testing.T) {
			t.Parallel()
			if got := checkout.RoundHalfUp(tt.amount, tt.percent); got != tt.want {
				t.Errorf("RoundHalfUp(%d, %d) = %d, want %d", tt.amount, tt.percent, got, tt.want)
			}
		})
	}
}

func TestRoundHalfUpRoundsAnExactHalfUpNotToTheEvenNumber(t *testing.T) {
	t.Parallel()
	// 10% of $52.35 is exactly 523.5 cents. "Banker's rounding" (half to even) would give
	// 524 here too, but 22.5 → 22 (see the table above); the two together pin the rule down.
	if got := checkout.RoundHalfUp(5235, 10); got != 524 {
		t.Errorf("RoundHalfUp(5235, 10) = %d, want 524", got)
	}
}

func TestRoundHalfUpMoreExamples(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		amount  int
		percent int
		want    int
	}{
		{name: "already whole: 10% of 500", amount: 500, percent: 10, want: 50},
		{name: "already whole: 8% of 50", amount: 50, percent: 8, want: 4},
		{name: "8% of $1.20 is 9.6, rounds up", amount: 120, percent: 8, want: 10},
		{name: "8% of $2.40 is 19.2, rounds down", amount: 240, percent: 8, want: 19},
		{name: "0 stays 0", amount: 0, percent: 10, want: 0},
		{name: "100% is the whole amount", amount: 899, percent: 100, want: 899},
		{name: "0% is nothing", amount: 899, percent: 0, want: 0},
		// 0.1 + 0.2 != 0.3 with floating point numbers. 10% of 3 cents is 0.3: it must
		// round to 0, not 1. And 8% of 1050 is exactly 84.0; a float might land on 83.99999.
		{name: "stays exact where floats drift: 10% of 3", amount: 3, percent: 10, want: 0},
		{name: "stays exact where floats drift: 8% of 1050", amount: 1050, percent: 8, want: 84},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := checkout.RoundHalfUp(tt.amount, tt.percent); got != tt.want {
				t.Errorf("RoundHalfUp(%d, %d) = %d, want %d", tt.amount, tt.percent, got, tt.want)
			}
		})
	}
}
