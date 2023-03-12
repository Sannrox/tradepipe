#!/usr/bin/env sh

set -eu


ROOT_PATH=$(dirname "$0")/..

. "${ROOT_PATH}/scripts/lib/openapi.sh"
.  "${ROOT_PATH}/scripts/lib/protoc.sh"


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