#!/usr/bin/env bash


# shellcheck disable=SC2034
set -o errexit
set -o nounset
set -o pipefail



function version::get_version() {
    if [[ -n "${GIT_COMMIT-}" ]] || GIT_COMMIT=$(git rev-parse --short "HEAD^{commit}" 2>/dev/null); then
        if [[ -z ${GIT_TREE_STATE-} ]]; then
            if git_status=$(git status --porcelain 2>/dev/null) && [[ -z "${git_status}" ]]; then
                GIT_TREE_STATE="clean"
            else
                GIT_TREE_STATE="dirty"
            fi
        fi

    if [[ -n "${GIT_VERSION-}" ]] || GIT_VERSION=$(git describe --tags --match 'v[0-9]*' --abbrev=14 "${GIT_COMMIT}^{commit}" 2>/dev/null); then
        if ! echo "${GIT_VERSION}" | grep -q -E "^[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9]+(\.[0-9]+)?)?$"; then
            echo "Version ${GIT_VERSION} is not a semantic version"
            return 1
        fi

    fi
   fi
}

version::get_version
