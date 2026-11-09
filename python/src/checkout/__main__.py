"""The command line program (rules §8): basket JSON in on stdin, receipt JSON out on stdout.

    echo '{"items":[{"sku":"APPLE","quantity":1}]}' | python -m checkout

All the logic is in command.py; this file only connects it to the terminal, which is why it
is left out of the in-process coverage and mutation reports and tested end-to-end instead
(tests/e2e/). Python runs this file when you type ``python -m checkout``.
"""

import sys

from checkout.command import main

if __name__ == "__main__":
    sys.exit(main())
