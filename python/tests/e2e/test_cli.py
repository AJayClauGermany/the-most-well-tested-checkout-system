"""End-to-end tests: run the real command line program (rules §8) as a separate process,
exactly as a user (or the shared runner in e2e/) would, and look only at what comes out.

`python -m checkout` is started with the standard library's subprocess module. Because the
test only sees stdout, stderr and the exit code, it would pass unchanged for a program
written in any language - which is the idea behind the shared runner e2e/run.sh.
"""

import json
import subprocess
import sys
from dataclasses import dataclass

import pytest
from spec_loader import SpecCase, load_spec_cases
from support.schemas import schema_errors

SPEC_CASES = load_spec_cases()


@dataclass(frozen=True)
class CliRun:
    stdout: str
    stderr: str
    exit_code: int


def run_cli(stdin_text: str) -> CliRun:
    """Runs `python -m checkout` once, with `stdin_text` on standard input.

    sys.executable is the Python running these tests, so the program uses the same
    environment (and the same installed `checkout` package) as the tests do.
    """
    completed = subprocess.run(
        [sys.executable, "-m", "checkout"],
        input=stdin_text,
        capture_output=True,
        text=True,
        check=False,  # the exit code is part of what we test, so do not raise on it
    )
    return CliRun(stdout=completed.stdout, stderr=completed.stderr, exit_code=completed.returncode)


# Every shared case, through the real program.
@pytest.mark.parametrize("case", SPEC_CASES, ids=[case.id for case in SPEC_CASES])
def test_spec_case_through_the_real_program(case: SpecCase) -> None:
    result = run_cli(json.dumps(case.basket_json) + "\n")

    if case.expected_error is not None:
        assert result.exit_code == 1
        assert json.loads(result.stdout) == {"error": case.expected_error.value}
    else:
        assert result.exit_code == 0
        receipt = json.loads(result.stdout)
        assert case.expected is not None
        for name, expected in vars(case.expected).items():
            assert receipt[name] == expected, name
        assert schema_errors("receipt.schema.json", receipt) == []


def test_prints_exactly_one_line_of_json() -> None:
    result = run_cli('{"items":[{"sku":"APPLE","quantity":1}]}')

    assert result.stdout.endswith("\n")
    assert len(result.stdout.strip().split("\n")) == 1
    assert result.stderr == ""


def test_exits_with_2_and_explains_on_stderr_when_the_input_is_not_json() -> None:
    result = run_cli("this is not json")

    assert result.exit_code == 2
    assert result.stdout == ""
    assert result.stderr.startswith("checkout: ")


def test_exits_with_2_when_the_json_is_not_a_basket() -> None:
    result = run_cli('{"hello": "world"}')

    assert result.exit_code == 2
    assert result.stdout == ""
    assert "not a basket" in result.stderr


def test_exits_with_2_on_empty_input() -> None:
    result = run_cli("")

    assert result.exit_code == 2
    assert result.stdout == ""


def test_the_error_output_follows_cli_error_schema_json() -> None:
    result = run_cli('{"items":[{"sku":"APPLE","quantity":1}],"coupon":"NOPE"}')

    assert result.exit_code == 1
    assert schema_errors("cli-error.schema.json", json.loads(result.stdout)) == []
