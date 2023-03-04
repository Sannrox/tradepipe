#!/usr/bin/env sh

set -eu

PREFIX="github.com/Sannrox/tradepipe"

ROOT_PATH=$(cd "$(dirname "$0")"/.. && pwd -P)
PROTOC_VERSION="3.7.1"


protoc_gen(){
find "${ROOT_PATH}/api/proto" -type f  -iname "*.proto" > "${ROOT_PATH}/protos.txt" 
while IFS= read -r file; do 
    protoc --go_out=.  --go_opt=module=$PREFIX  \
    --go-grpc_out=.  --go-grpc_opt=module=$PREFIX \
    -I "${ROOT_PATH}" \
    "$file"; 
done < "${ROOT_PATH}/protos.txt"
rm "${ROOT_PATH}/protos.txt"
}


protoc_install(){
	if ! type "protoc" > /dev/null; then 
        wget https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/protoc-${PROTOC_VERSION}-linux-x86_64.zip; 
        unzip protoc-${PROTOC_VERSION}-linux-x86_64.zip -d protoc3; 
        sudo mv protoc3/bin/* /usr/local/bin/; \
        sudo mv protoc3/include/* /usr/local/include/; 
        go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
        go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	fi
}