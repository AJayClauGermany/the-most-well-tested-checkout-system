#!/usr/bin/env bash
# Shared end-to-end runner (rules §8).
#
# Feeds every shared test case in spec/test-cases/ to a checkout command line program,
# one basket at a time on standard input, and compares what comes back:
#
#   * a case with "expected"        → exit code 0 and the six totals must match
#   * a case with "expected_error"  → exit code 1 and {"error": "<CODE>"} must match
#
# The program can be written in any language. Examples:
#
#   e2e/run.sh python -m checkout               (run from the python/ folder)
#   e2e/run.sh node dist/cli.js
#   e2e/run.sh ./go/bin/checkout
#   e2e/run.sh java -jar java/target/checkout.jar
#   e2e/run.sh dotnet csharp/src/Checkout.Cli/bin/Debug/net10.0/Checkout.Cli.dll
#   e2e/run.sh rust/target/debug/checkout
#
# Needs: bash and jq (https://jqlang.github.io/jq/). Exit code: 0 if every case passed.
set -euo pipefail

if [ "$#" -eq 0 ]; then
  echo "usage: $0 <command that runs the checkout CLI> [args...]" >&2
  exit 2
fi
if ! command -v jq >/dev/null 2>&1; then
  echo "error: this script needs jq (https://jqlang.github.io/jq/)" >&2
  exit 2
fi

here="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cases_dir="$here/../spec/test-cases"
stderr_file="$(mktemp)"
trap 'rm -f "$stderr_file"' EXIT

passed=0
failed=0

report_fail() {
  failed=$((failed + 1))
  echo "FAIL  $1"
  echo "      $2"
  if [ -s "$stderr_file" ]; then
    sed 's/^/      stderr: /' "$stderr_file"
  fi
}

for file in "$cases_dir"/*.json; do
  echo "== $(basename "$file")"
  while IFS= read -r case; do
    id="$(jq -r '.id' <<<"$case")"
    basket="$(jq -c '.basket' <<<"$case")"

    # Run the program. Its exit code is part of what we check, so don't let it stop us.
    set +e
    output="$(printf '%s\n' "$basket" | "$@" 2>"$stderr_file")"
    exit_code=$?
    set -e

    if jq -e 'has("expected_error")' <<<"$case" >/dev/null; then
      want="$(jq -r '.expected_error' <<<"$case")"
      got="$(jq -r '.error // empty' <<<"$output" 2>/dev/null || true)"
      if [ "$exit_code" -eq 1 ] && [ "$got" = "$want" ]; then
        passed=$((passed + 1)); echo "ok    $id"
      else
        report_fail "$id" "expected error $want with exit code 1; got error '${got:-<none>}' with exit code $exit_code"
      fi
    else
      want="$(jq -cS '.expected' <<<"$case")"
      got="$(jq -cS '{subtotal, promotion_discount, coupon_discount, shipping, tax, total}' <<<"$output" 2>/dev/null || true)"
      if [ "$exit_code" -eq 0 ] && [ "$got" = "$want" ]; then
        passed=$((passed + 1)); echo "ok    $id"
      else
        report_fail "$id" "expected $want with exit code 0; got ${got:-<no valid JSON>} with exit code $exit_code"
      fi
    fi
  done < <(jq -c '.cases[]' "$file")
done

echo
echo "$passed passed, $failed failed"
[ "$failed" -eq 0 ]
