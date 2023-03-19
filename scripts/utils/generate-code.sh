#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail
set -x



ROOT_PATH=$(dirname "$0")/../..

source "${ROOT_PATH}/scripts/lib/init.sh"


codegen(){
    openapi_gen
    protoc_gen
}

codegen_install(){
    openapi_install
    protoc_install
}


case "$1" in
    "install")
        codegen_install
        ;;
    "codegen")
        codegen
        ;;
    *)
        echo "Usage: $0 [install|codegen]"
        exit 1
        ;;
esac