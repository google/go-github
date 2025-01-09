#!/bin/sh
#/ script/generate.sh runs go generate on all modules in this repo.
#/ `script/generate.sh --check` checks that the generated files are up to date.

set -e

CDPATH="" cd -- "$(dirname -- "$0")/.."

if [ "$1" = "--check" ]; then
  GENTEMP="$(mktemp -d)"
  git worktree add -q --detach "$GENTEMP"
  trap 'git worktree remove -f "$GENTEMP"; rm -rf "$GENTEMP"' EXIT
  for f in $(git ls-files -com --exclude-standard); do
    target="$GENTEMP/$f"
    mkdir -p "$(dirname -- "$target")"
    cp "$f" "$target"
  done
  if [ -f "$(pwd)"/bin ]; then
    ln -s "$(pwd)"/bin "$GENTEMP"/bin
  fi
  (
    cd "$GENTEMP"
    git add .
    git -c user.name='bot' -c user.email='bot@localhost' commit -m "generate" -q --allow-empty
    script/generate.sh
    [ -z "$(git status --porcelain)" ] || {
      msg="Generated files are out of date. Please run script/generate.sh and commit the results"
      if [ -n "$GITHUB_ACTIONS" ]; then
        echo "::error ::$msg"
      else
        echo "$msg" 1>&2
      fi
      git diff
      exit 1
    }
  )
  exit 0
fi

MOD_DIRS="$(git ls-files '*go.mod' | xargs dirname | sort)"

for dir in $MOD_DIRS; do
  (
    cd "$dir"
    go generate ./...
    GOTOOLCHAIN="go1.22+auto" go mod tidy
  )
done
