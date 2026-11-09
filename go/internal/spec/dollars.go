package spec

// Test helper: turns "$12.39" into 1239 cents – without ever touching a floating point
// number (spec/features/README.md insists, and rules §1 explains why).

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// dollarPattern matches "$12.39", "12.39", "$5" or "$5.5": optional $, digits, and up to
// two decimals. Group 1 is the dollars, group 2 the cents (possibly one digit, or empty).
var dollarPattern = regexp.MustCompile(`^\$?(\d+)(?:\.(\d{1,2}))?$`)

// DollarsToCents converts "$12.39" (or "12.39") to 1239 and "$5" to 500, using whole numbers only.
func DollarsToCents(text string) (int, error) {
	match := dollarPattern.FindStringSubmatch(strings.TrimSpace(text))
	if match == nil {
		return 0, fmt.Errorf("not a dollar amount: %q", text)
	}
	dollars, err := strconv.Atoi(match[1])
	if err != nil {
		return 0, fmt.Errorf("not a dollar amount: %q: %w", text, err)
	}
	cents := 0
	if match[2] != "" {
		digits := match[2]
		if len(digits) == 1 { // "$5.5" means 50 cents, not 5
			digits += "0"
		}
		cents, _ = strconv.Atoi(digits) // two digits always parse
	}
	return dollars*100 + cents, nil
}
