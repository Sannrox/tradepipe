#!/usr/bin/env sh

# shellcheck disable=SC2034

set -eu

ROOT_PATH=$(dirname "$0")/..

. "${ROOT_PATH}/build/lib/docker.sh"

GRPC_SERVER_IMAGE="${GRPC_SERVER_IMAGE:-"grpc-server:1.0.0"}"
export GRPC_SERVER_IMAGE

HTTP_SERVER_IMAGE="${HTTP_SERVER_IMAGE:-"http-server:1.0.0"}"
export HTTP_SERVER_IMAGE


targets=$(build_get_docker_wrapped_binaries)

for target in ${targets}; do 
    docker_build "${target}" "${ROOT_PATH}/build/Dockerfile"
done

