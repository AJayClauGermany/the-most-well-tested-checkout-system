#!/usr/bin/env bash
# Fails when code coverage drops below the thresholds below.
#
# Reads the merged report that ReportGenerator writes (reports/coverage/Cobertura.xml, see the
# README) and compares its line and branch rates with the minimums. ReportGenerator's own
# -minimumCoverageThresholds option is not in the version we pin, and this is short enough.
#
#   ./check-coverage.sh [path/to/Cobertura.xml]
set -euo pipefail

MIN_LINE_PERCENT=100
MIN_BRANCH_PERCENT=95

report="${1:-reports/coverage/Cobertura.xml}"
if [ ! -f "$report" ]; then
  echo "error: $report not found - run the coverage commands from the README first" >&2
  exit 2
fi

# The first <coverage ...> element holds the totals, e.g. line-rate="0.9898" branch-rate="0.96".
header="$(grep -m 1 -o '<coverage [^>]*>' "$report")"
rate_of() { sed -E "s/.* $1=\"([0-9.]+)\".*/\1/" <<<"$header"; }

# Rates are fractions (0.9898); compare in whole hundredths of a percent with awk, not floats in bash.
check() {
  local name="$1" rate="$2" minimum="$3"
  if awk -v rate="$rate" -v min="$minimum" 'BEGIN { exit !(int(rate * 10000 + 0.5) >= min * 100) }'; then
    printf 'ok    %s coverage %.2f%% (minimum %s%%)\n' "$name" "$(awk -v r="$rate" 'BEGIN { print r * 100 }')" "$minimum"
  else
    printf 'FAIL  %s coverage %.2f%% is below the minimum of %s%%\n' "$name" "$(awk -v r="$rate" 'BEGIN { print r * 100 }')" "$minimum"
    return 1
  fi
}

status=0
check line "$(rate_of line-rate)" "$MIN_LINE_PERCENT" || status=1
check branch "$(rate_of branch-rate)" "$MIN_BRANCH_PERCENT" || status=1
exit "$status"
