#!/bin/sh
#/ script/fmt.sh runs go fmt on all go files in the project.

set -e

CDPATH="" cd -- "$(dirname -- "$0")/.."

MOD_DIRS="$(git ls-files '*go.mod' | xargs dirname | sort)"

for dir in $MOD_DIRS; do
  (
    cd "$dir"
    go fmt ./...
  )
done
