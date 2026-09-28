#!/usr/bin/env sh
set -eu

find ./scripts \( -name '*.sh' -o -name 'configure' \) -type f -exec shellcheck {} +
