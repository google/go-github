#!/bin/bash -e
#/ script/run-check-structfield-settings.sh runs tools/check-structfield-settings, which
#/ checks the `structfield` linter exceptions in `.golangci.yml` against the code.
#/
#/ With no arguments it reports the exceptions that are no longer needed and the ones that
#/ are listed twice, and exits non-zero. `-fix` removes and sorts them instead.

CDPATH="" cd -- "$(dirname -- "$0")/.."

go run -C tools/check-structfield-settings . "$@"
