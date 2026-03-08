#!/bin/sh
#/ [ CHECK_GITHUB_OPENAPI=1 ] script/lint.sh runs linters and validates generated files.
#/ When CHECK_GITHUB is set, it validates that openapi_operations.yaml is consistent with the
#/ descriptions from github.com/github/rest-api-description.

set -e

CUSTOM_GCL="$(script/setup-custom-gcl.sh)"

CDPATH="" cd -- "$(dirname -- "$0")/.."

EXIT_CODE=0

fail() {
  echo "$@"
  EXIT_CODE=1
}

MOD_DIRS="$(git ls-files '*go.mod' | xargs dirname | sort -u)"

# Number of module lint jobs to run concurrently.
# Override with LINT_JOBS, otherwise use detected CPU count.
: "${LINT_JOBS:=$(getconf _NPROCESSORS_ONLN 2>/dev/null || sysctl -n hw.ncpu 2>/dev/null || echo 4)}"

LINT_DIRS="$(printf '%s\n' "$MOD_DIRS" | grep -v '^example/newreposecretwithlibsodium$')"

LINT_FAILED=0
RUNNING=0
PIDS=""

wait_pids() {
  for pid in $PIDS; do
    if ! wait "$pid"; then
      LINT_FAILED=1
    fi
  done
  PIDS=""
  RUNNING=0
}

for dir in $LINT_DIRS; do
  (
    echo linting "$dir"
    cd "$dir"
    "$CUSTOM_GCL" run
  ) &

  PIDS="$PIDS $!"
  RUNNING=$((RUNNING + 1))

  if [ "$RUNNING" -ge "$LINT_JOBS" ]; then
    wait_pids
  fi
done

wait_pids

if [ "$LINT_FAILED" -ne 0 ]; then
  fail "failed linting one or more module directories"
fi

if [ -n "$CHECK_GITHUB_OPENAPI" ]; then
  echo validating openapi_operations.yaml
  script/metadata.sh update-openapi --validate || fail "failed validating openapi_operations.yaml"
fi

echo validating generated files
script/generate.sh --check || fail "failed validating generated files"

exit "$EXIT_CODE"
