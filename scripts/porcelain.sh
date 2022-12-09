#!/usr/bin/env sh
set -e
if [ -n "$(git status --porcelain)" ]; then
    echo "aborting: non-commited files"
    exit 1
fi
