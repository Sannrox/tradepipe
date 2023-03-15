#!/usr/bin/env sh

# shellcheck disable=SC2034

set -eu

ROOT_PATH=$(dirname "$0")/../..

run_cmd(){
    filname="${##*/validate-}"
    
    "$@"
}

run_validate(){
    for script in $(echo "$1"); do
            echo  "Validating $(basename "${script}")"
            run_cmd "${script}" || true
    done
}


run_validate "${ROOT_PATH}/scripts/validate/*.sh"