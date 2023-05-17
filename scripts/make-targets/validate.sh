#!/usr/bin/env bash

# shellcheck disable=SC2034

set -o errexit
set -o nounset
set -o pipefail


ROOT_PATH=$(dirname "$0")/../..

run_cmd(){
    filname="${##*/validate-}"
    
    "$@"
}

run_validate(){
    for script in $(echo "$1"); do
            echo  "Validating $(basename "${script}")"
            run_cmd "${script}"
    done
}


run_validate "${ROOT_PATH}/scripts/validate/*.sh"