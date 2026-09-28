#!/usr/bin/env sh
set -eu

version="${1:-}"
if [ -z "$version" ]; then
  echo "usage: mise run release:pr -- <version>" >&2
  exit 2
fi

if ! command -v gh >/dev/null 2>&1; then
  echo "gh must be installed and authenticated before creating a pull request" >&2
  exit 1
fi

current_branch=$(git branch --show-current)
branch="release/$version"
title="^ E(release): release $version"
body=$(cat <<EOF
Prepare release $version.

This pull request moves the manual changelog notes from Unreleased to $version and updates VERSION.
EOF
)

./scripts/release-branch.sh "$version"
git push -u origin "$branch"
gh pr create --base "$current_branch" --head "$branch" --title "$title" --body "$body" --label release
git switch "$current_branch"
