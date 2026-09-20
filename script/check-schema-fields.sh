#!/bin/sh
#/ script/check-schema-fields.sh runs ./tools/schemafields in the repository root with the
#/ given arguments.
#/
#/ It reports every request body struct field whose optionality disagrees with the OpenAPI
#/ request body schema of the operation the method is annotated with. `-fix` repairs the
#/ findings that can be repaired automatically.

set -e

CDPATH="" cd -- "$(dirname -- "$0")/.."
REPO_DIR="$(pwd)"

(
  cd tools/schemafields
  go build -o "$REPO_DIR"/bin/schemafields
)

exec bin/schemafields "$@"
