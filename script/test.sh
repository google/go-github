#!/bin/sh
#/ script/test.sh runs tests on each go module in go-github. Arguments are passed to each go test invocation.
#/ "-race -covermode atomic ./..." is used when no arguments are given.

set -e

CDPATH="" cd -- "$(dirname -- "$0")/.."

if [ "$#" = "0" ]; then
  set -- -race -covermode atomic ./...
fi

MOD_DIRS="$(git ls-files '*go.mod' | xargs dirname | sort)"

for dir in $MOD_DIRS; do
  [ "$dir" = "example/newreposecretwithlibsodium" ] && continue
  echo "testing $dir"
  (
    cd "$dir"
    go test "$@"
  ) || FAILED=1
done

if [ -n "$FAILED" ]; then
  exit 1
fi
