#!/bin/sh
#/ script/setup-custom-gcl.sh ensures custom golangci-lint is installed.
#/ It returns the path to the custom-gcl binary.

set -e

# should be in sync with .custom-gcl.yml
GOLANGCI_LINT_VERSION="v2.10.1"

# should in sync with fmt.sh and lint.sh
BIN="$(pwd -P)"/bin

mkdir -p "$BIN"

# install golangci-lint and custom-gcl in ./bin if they don't exist with the correct version
if ! "$BIN"/custom-gcl version --short 2> /dev/null | grep -q "$GOLANGCI_LINT_VERSION"; then
  curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b "$BIN" "$GOLANGCI_LINT_VERSION"
  "$BIN"/golangci-lint custom --name custom-gcl --destination "$BIN"
fi

echo "$BIN/custom-gcl"
