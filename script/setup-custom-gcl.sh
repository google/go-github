#!/bin/sh
#/ script/setup-custom-gcl.sh ensures custom golangci-lint is installed.
#/ It returns the path to the custom-gcl binary.

set -e

GOLANGCI_LINT_VERSION="2.9.0"

# should in sync with fmt.sh and lint.sh
BIN="$(pwd -P)"/bin

mkdir -p "$BIN"

# install golangci-lint and custom-gcl in ./bin if they don't exist with the correct version
if ! "$BIN"/custom-gcl --version 2> /dev/null | grep -q "$GOLANGCI_LINT_VERSION"; then
  GOBIN="$BIN" go install "github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v$GOLANGCI_LINT_VERSION"
  "$BIN"/golangci-lint custom --name custom-gcl --destination "$BIN"
fi

echo "$BIN/custom-gcl"
