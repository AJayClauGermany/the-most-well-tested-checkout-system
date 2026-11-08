"""Checks that every test case id in spec/test-cases is unique across all files.

JSON Schema can check each file on its own, but not across files, so this tiny
script does that one extra check. Run it with:  python3 spec/check_unique_ids.py
"""

import json
import sys
from pathlib import Path

test_case_folder = Path(__file__).parent / "test-cases"
seen = {}  # id -> file it was first seen in
problems = []

for path in sorted(test_case_folder.glob("*.json")):
    for case in json.loads(path.read_text(encoding="utf-8"))["cases"]:
        case_id = case["id"]
        if case_id in seen:
            problems.append(f"{case_id!r} is in both {seen[case_id]} and {path.name}")
        else:
            seen[case_id] = path.name

if problems:
    print("Duplicate test case ids found:")
    for problem in problems:
        print("  -", problem)
    sys.exit(1)

print(f"OK: {len(seen)} test case ids, all unique.")
