#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail




openapi_gen(){
    oapi-codegen --generate client,server,types,spec \
    --package api \
    -o "${ROOT_PATH}/rest/api/rest.gen.go" \
     "${ROOT_PATH}/api/openapi/openapi.yaml"
}


openapi_install(){
    go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.12.4
}
