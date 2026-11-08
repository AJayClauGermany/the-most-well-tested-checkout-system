"""Money helpers (rules §1 and §7). All amounts are whole cents."""


def round_half_up(amount: int, percent: int) -> int:
    """``percent``% of ``amount``, rounded half up to a whole cent (rules §7).

    Not implemented yet - see tests/unit/test_money.py for what it must do.
    """
    raise NotImplementedError(f"round_half_up({amount}, {percent}) is not implemented yet")
