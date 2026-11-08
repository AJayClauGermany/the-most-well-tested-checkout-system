#!/usr/bin/env bash
# A FAKE checkout CLI used only to test e2e/run.sh itself.
#
# It does not calculate anything: it reads a basket from standard input, looks for a
# shared test case with exactly that basket, and prints that case's expected answer.
# If run.sh works, every case passes against this oracle. (Do not copy this idea into a
# real implementation – it is the testing equivalent of writing the answers on your arm.)
set -euo pipefail
here="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
basket="$(jq -c . )"   # reads standard input
case="$(jq -c --argjson b "$basket" 'first(.cases[] | select(.basket == $b))' "$here"/../../spec/test-cases/*.json | head -n 1)"
if [ -z "$case" ]; then
  echo '{"error": "ORACLE_HAS_NO_SUCH_BASKET"}'
  exit 1
fi
if jq -e 'has("expected_error")' <<<"$case" >/dev/null; then
  jq -c '{error: .expected_error}' <<<"$case"
  exit 1
fi
# A real receipt also has "lines"; the runner only compares the totals.
jq -c '.expected' <<<"$case"
