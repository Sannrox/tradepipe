#!/usr/bin/env bash

# shellcheck disable=SC2034

set -o errexit
set -o nounset
set -o pipefail

ROOT_PATH=$(dirname "$0")/../..



run_update(){
    local -r pattern="$1"
    for script in ${pattern}; do
            echo  "Updating $(basename "${script}")"
            ${script}
    done
}

run_update "${ROOT_PATH}/scripts/update/*.sh"
