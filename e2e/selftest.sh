#!/usr/bin/env bash
# Tests the runner itself (yes, test code deserves tests too).
#   1. Against an "oracle" that answers from the spec, every case must pass.
#   2. Against a program that is wrong on purpose (it always says EMPTY_BASKET), only the
#      cases that really expect EMPTY_BASKET may pass, and the runner must exit non-zero.
set -euo pipefail
here="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cases=("$here"/../spec/test-cases/*.json)
total="$(jq -s 'map(.cases | length) | add' "${cases[@]}")"
empty_basket="$(jq -s 'map(.cases[] | select(.expected_error == "EMPTY_BASKET")) | length' "${cases[@]}")"
log="$(mktemp)"
trap 'rm -f "$log"' EXIT

echo "--- selftest 1: the oracle passes all $total cases"
if ! "$here/run.sh" "$here/selftest/oracle.sh" >"$log" 2>&1; then
  cat "$log"; echo "FAILED: runner reported failures against the oracle"; exit 1
fi
tail -n 1 "$log"
grep -qx "$total passed, 0 failed" "$log"

echo "--- selftest 2: a wrong program is caught ($empty_basket of $total cases expect EMPTY_BASKET)"
if "$here/run.sh" "$here/selftest/always-empty-basket.sh" >"$log" 2>&1; then
  cat "$log"; echo "FAILED: runner should have exited non-zero"; exit 1
fi
tail -n 1 "$log"
grep -qx "$empty_basket passed, $((total - empty_basket)) failed" "$log"

echo "selftest OK"
