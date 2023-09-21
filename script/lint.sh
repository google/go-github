#!/bin/sh
#/ script/lint.sh runs linters and validates generated files.

set -e

GOLANGCI_LINT_VERSION="1.54.2"

CDPATH="" cd -- "$(dirname -- "$0")/.."
BIN="$(pwd -P)"/bin

mkdir -p "$BIN"

# install golangci-lint bin/golangci-lint doesn't exist with the correct version
if ! "$BIN"/golangci-lint --version 2> /dev/null | grep -q "$GOLANGCI_LINT_VERSION"; then
  GOBIN="$BIN" go install "github.com/golangci/golangci-lint/cmd/golangci-lint@v$GOLANGCI_LINT_VERSION"
fi

MOD_DIRS="$(git ls-files '*go.mod' | xargs dirname | sort)"

for dir in $MOD_DIRS; do
  [ "$dir" = "example/newreposecretwithlibsodium" ] && continue
  echo linting "$dir"
  (
    cd "$dir"
    # github actions output when running in an action
    if [ -n "$GITHUB_ACTIONS" ]; then
      "$BIN"/golangci-lint run --path-prefix "$dir" --out-format github-actions
    else
      "$BIN"/golangci-lint run --path-prefix "$dir"
    fi
  ) || FAILED=1
done

script/generate.sh --check || FAILED=1

if [ -n "$FAILED" ]; then
  exit 1
fi
