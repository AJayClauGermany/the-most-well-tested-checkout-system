#!/usr/bin/env bash
# A FAKE checkout CLI that is wrong on purpose: whatever the basket, it reports EMPTY_BASKET.
# Used to check that e2e/run.sh really fails cases that are wrong.
cat >/dev/null   # swallow the basket
echo '{"error": "EMPTY_BASKET"}'
exit 1
