#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail



clean_with_patter(){
    pattern="$1"

    for path in $(echo "$pattern"); do
        echo "Cleaning ${path}"
        rm -rf "${path}"
    done
}

