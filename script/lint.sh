#!/bin/sh
#/ [ CHECK_GITHUB_OPENAPI=1 ] script/lint.sh runs linters and validates generated files.
#/ When CHECK_GITHUB is set, it validates that openapi_operations.yaml is consistent with the
#/ descriptions from github.com/github/rest-api-description.

set -e

GOLANGCI_LINT_VERSION="1.63.4"

CDPATH="" cd -- "$(dirname -- "$0")/.."
BIN="$(pwd -P)"/bin

mkdir -p "$BIN"

EXIT_CODE=0

fail() {
  echo "$@"
  EXIT_CODE=1
}

# install golangci-lint and custom-gcl in ./bin if they don't exist with the correct version
if ! "$BIN"/custom-gcl --version 2> /dev/null | grep -q "$GOLANGCI_LINT_VERSION"; then
  GOBIN="$BIN" go install "github.com/golangci/golangci-lint/cmd/golangci-lint@v$GOLANGCI_LINT_VERSION"
  "$BIN"/golangci-lint custom && mv ./custom-gcl "$BIN"
fi

MOD_DIRS="$(git ls-files '*go.mod' | xargs dirname | sort)"

for dir in $MOD_DIRS; do
  [ "$dir" = "example/newreposecretwithlibsodium" ] && continue
  echo linting "$dir"
  (
    cd "$dir"
    # github actions output when running in an action
    if [ -n "$GITHUB_ACTIONS" ]; then
      "$BIN"/custom-gcl run --path-prefix "$dir" --out-format colored-line-number
    else
      "$BIN"/custom-gcl run --path-prefix "$dir"
    fi
  ) || fail "failed linting $dir"
done

if [ -n "$CHECK_GITHUB_OPENAPI" ]; then
  echo validating openapi_operations.yaml
  script/metadata.sh update-openapi --validate || fail "failed validating openapi_operations.yaml"
fi

echo validating generated files
script/generate.sh --check || fail "failed validating generated files"

[ -z "$FAILED" ] || exit 1

exit "$EXIT_CODE"
