#!/usr/bin/env sh
set -e

VERSION_TYPE=$1
OLD_VERSION=$(git ls-remote --tags origin | cut -d'/' -f3 | sort --version-sort | tail -n1)
NEW_VERSION=$(echo "${OLD_VERSION}" | semv increment "${VERSION_TYPE}")
echo "updating from ${OLD_VERSION} to ${NEW_VERSION}"

git tag "${NEW_VERSION}"
git push origin "${NEW_VERSION}"
