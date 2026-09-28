#!/usr/bin/env sh
set -eu

version="${1:-}"
if [ -z "$version" ]; then
  echo "usage: mise run release:branch -- <version>" >&2
  exit 2
fi

base_branch=trunk
current_branch=$(git branch --show-current)
if [ "$current_branch" != "$base_branch" ]; then
  echo "release branch must be created from $base_branch; current branch is ${current_branch:-detached HEAD}" >&2
  exit 1
fi

branch="release/$version"
changed_files=$(git status --short --untracked-files=no)
if [ -n "$changed_files" ]; then
  echo "Tracked files must be clean before preparing a release." >&2
  echo "Changed tracked files:" >&2
  printf '%s\n' "$changed_files" >&2
  exit 1
fi

untracked_files=$(git ls-files --others --exclude-standard)
if [ -n "$untracked_files" ]; then
  echo "Untracked files are present:"
  printf '%s\n' "$untracked_files"
  printf 'Continue and leave these files untracked? [y/N] '
  read -r answer
  case "$answer" in
    y | Y | yes | YES)
      ;;
    *)
      echo "release branch creation canceled" >&2
      exit 1
      ;;
  esac
fi

git switch -c "$branch"
go run ./internal/tools/releaseprep --version "$version"
git add CHANGELOG.md VERSION
git commit -m "^ E(release): release $version"

echo "Created $branch. Push it with: git push -u origin $branch"
