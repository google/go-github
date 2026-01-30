#!/bin/sh
#/ [ CHECK_GITHUB_OPENAPI=1 ] script/lint.sh runs linters and validates generated files.
#/ When CHECK_GITHUB is set, it validates that openapi_operations.yaml is consistent with the
#/ descriptions from github.com/github/rest-api-description.

set -e

CUSTOM_GCL="$(script/setup-custom-gcl.sh)"

CDPATH="" cd -- "$(dirname -- "$0")/.."
BIN="$(pwd -P)"/bin

EXIT_CODE=0

fail() {
  echo "$@"
  EXIT_CODE=1
}

MOD_DIRS="$(git ls-files '*go.mod' | xargs dirname | sort)"

for dir in $MOD_DIRS; do
  [ "$dir" = "example/newreposecretwithlibsodium" ] && continue
  echo linting "$dir"
  (
    cd "$dir"
    "$CUSTOM_GCL" run
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
