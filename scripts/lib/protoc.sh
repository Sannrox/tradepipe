#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail


PREFIX=$(grep module go.mod | cut -d " " -f2 )

PROTOC_VERSION="3.7.1"


function protoc::gen(){
find "${ROOT_PATH}/api/protobuf" -type f  -iname "*.proto" > "${ROOT_PATH}/protos.txt"
while IFS= read -r file; do
    protoc_command=(protoc --go_out=.  --go_opt=module="$PREFIX"  --go-grpc_out=.
        --go-grpc_opt=module="$PREFIX" -I "${ROOT_PATH}" "$file")
    output=$("${protoc_command[@]}" 2>&1) || {
    cat <<EOF >&2
    Error: Protoc failed to generate code.
    ${output}

    to retry manually, run:
    ${protoc_command[@]}
EOF
    exit 1
    }

done < "${ROOT_PATH}/protos.txt"
rm "${ROOT_PATH}/protos.txt"
}


function protoc::install(){
    (
    wget https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/protoc-${PROTOC_VERSION}-linux-x86_64.zip;
    unzip protoc-${PROTOC_VERSION}-linux-x86_64.zip -d protoc3;
    sudo mv protoc3/bin/* /usr/local/bin/; \
    sudo mv protoc3/include/* /usr/local/include/;
    go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
    )
}



function protoc::check(){
    if [[ -z $(which protoc) ]]; then
        echo "protoc not found"
        return 1
    fi
}
