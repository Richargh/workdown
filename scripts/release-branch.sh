#!/usr/bin/env sh
set -eu

version="${1:-}"
if [ -z "$version" ]; then
  echo "usage: mise run release:branch -- <version>" >&2
  exit 2
fi

branch="release/$version"
if [ -n "$(git status --porcelain)" ]; then
  echo "working tree must be clean before preparing a release" >&2
  exit 1
fi

git switch -c "$branch"
go run ./internal/tools/releaseprep --version "$version"
git add CHANGELOG.md VERSION
git commit -m "^ E(release): release $version"

echo "Created $branch. Push it with: git push -u origin $branch"
