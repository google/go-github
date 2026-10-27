#!/bin/sh
#/ script/setup-custom-gcl.sh ensures custom golangci-lint is installed.
#/ It returns the path to the custom-gcl binary.

set -e

ROOT_DIR="$(cd -- "$(dirname -- "$0")/.." >/dev/null 2>&1 && pwd -P)"
BIN_DIR="${ROOT_DIR}/bin"
CUSTOM_GCL_CONFIG="${ROOT_DIR}/.custom-gcl.yml"
CUSTOM_GCL="${BIN_DIR}/custom-gcl"

GOLANGCI_LINT_VERSION="$(
  sed -n 's/^[[:space:]]*version:[[:space:]]*//p' "${CUSTOM_GCL_CONFIG}" |
    sed '1{s/[[:space:]]*#.*$//;s/^"//;s/"$//;s/[[:space:]]*$//;p;};d'
)"

if [ -z "${GOLANGCI_LINT_VERSION}" ]; then
  echo "Error: could not determine golangci-lint version from ${CUSTOM_GCL_CONFIG}" >&2
  exit 1
fi

mkdir -p "${BIN_DIR}"

# check if the custom-gcl binary exists and has the correct version
if ! "${CUSTOM_GCL}" version --short 2>/dev/null | grep -q "${GOLANGCI_LINT_VERSION}"; then
  GCL="golangci-lint"
  if ! command -v "${GCL}" >/dev/null 2>&1 || ! "${GCL}" version --short 2>/dev/null | grep -q "^v*2\."; then
    curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b "${BIN_DIR}" "${GOLANGCI_LINT_VERSION}"
    GCL="${BIN_DIR}/golangci-lint"
  fi

  "${GCL}" custom --name custom-gcl --destination "${BIN_DIR}"
fi

echo "${CUSTOM_GCL}"
