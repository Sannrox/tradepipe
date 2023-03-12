#!/usr/bin/env sh

set -eu


clean_with_patter(){
    pattern="$1"

    for path in $(echo "$pattern"); do
        echo "Cleaning ${path}"
        rm -rf "${path}"
    done
}

