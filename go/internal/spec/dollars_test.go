package spec

// Test code deserves tests too: the BDD steps rely on DollarsToCents, so a bug there
// would hide bugs in the checkout.

import (
	"strings"
	"testing"
)

func TestDollarsToCents(t *testing.T) {
	t.Parallel()
	tests := []struct {
		text string
		want int
	}{
		{text: "$0.00", want: 0},
		{text: "$0.50", want: 50},
		{text: "$12.39", want: 1239},
		{text: "$5", want: 500},
		{text: "$5.5", want: 550},
		{text: "  $8.99 ", want: 899},
		{text: "0.04", want: 4},
		{text: "$17.98", want: 1798},
	}

	for _, tt := range tests {
		t.Run(tt.text, func(t *testing.T) {
			t.Parallel()
			got, err := DollarsToCents(tt.text)

			if err != nil {
				t.Fatalf("DollarsToCents(%q) returned an error: %v", tt.text, err)
			}
			if got != tt.want {
				t.Errorf("DollarsToCents(%q) = %d, want %d", tt.text, got, tt.want)
			}
		})
	}
}

func TestDollarsToCentsRefusesTextThatIsNotAnAmount(t *testing.T) {
	t.Parallel()
	for _, text := range []string{"ten dollars", "$1.234", "$-1.00", "", "$", "$1.", "99999999999999999999999"} {
		_, err := DollarsToCents(text)

		if err == nil || !strings.Contains(err.Error(), "not a dollar amount") {
			t.Errorf("DollarsToCents(%q) error = %v, want 'not a dollar amount'", text, err)
		}
	}
}
