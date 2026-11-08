"""Test helper: turns "$12.39" into 1239 cents - without ever touching a floating point
number (spec/features/README.md insists, and rules §1 explains why).
"""

import re

# A dollar sign (optional), digits, and optionally a dot with one or two more digits.
_DOLLARS = re.compile(r"^\$?(?P<dollars>\d+)(?:\.(?P<cents>\d{1,2}))?$")


def dollars_to_cents(text: str) -> int:
    """``"$12.39"`` or ``"12.39"`` -> ``1239``. ``"$5"`` -> ``500``. ``"$0.5"`` -> ``50``.

    >>> dollars_to_cents("$12.39")
    1239
    >>> dollars_to_cents("$5")
    500
    """
    match = _DOLLARS.match(text.strip())
    if match is None:
        raise ValueError(f"not a dollar amount: {text!r}")
    dollars = int(match.group("dollars"))
    # "5" after the dot means 50 cents, so pad a single digit on the right with a zero.
    cents = int((match.group("cents") or "0").ljust(2, "0"))
    return dollars * 100 + cents


class Money:
    """A custom type for pytest-bdd's ``parsers.parse``: ``{amount:Money}`` matches ``$12.39``
    in a step sentence and hands the step function ``1239``.

    The ``parse`` library looks for a ``pattern`` attribute on a custom type, and calls the
    type with the matched text.
    """

    pattern = r"\$\d+(?:\.\d{1,2})?"

    def __call__(self, text: str) -> int:
        return dollars_to_cents(text)
