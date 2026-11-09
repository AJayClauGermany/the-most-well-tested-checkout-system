"""Money helpers (rules §1 and §7). All amounts are whole cents."""


def round_half_up(amount: int, percent: int) -> int:
    """``percent``% of ``amount``, rounded half up to a whole cent (rules §7).

    Done with whole numbers only, as the hint in rules §7 suggests: ``amount * percent`` is
    the exact answer in hundredths of a cent; adding 50 hundredths and dividing by 100 with
    ``//`` (which rounds down) is the same as rounding half up. No floating point, so no
    ``0.1 + 0.2`` surprises - and no banker's rounding, which Python's ``round()`` does.

    The examples below are *doctests*: ``pytest --doctest-modules`` runs them as tests.

    >>> round_half_up(224, 10)   # 22.4 rounds down
    22
    >>> round_half_up(225, 10)   # 22.5 rounds UP (round(22.5) would give 22!)
    23
    >>> round_half_up(226, 10)   # 22.6 rounds up
    23
    >>> round_half_up(120, 8)    # 9.6 cents of tax on $1.20
    10
    """
    return (amount * percent + 50) // 100
