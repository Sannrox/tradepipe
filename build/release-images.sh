#!/usr/bin/env bash

# shellcheck disable=SC2034

set -o errexit
set -o nounset
set -o pipefail
set -x

ROOT_PATH=$(dirname "${BASH_SOURCE[0]}")/..


source "${ROOT_PATH}/build/general.sh"
source "${ROOT_PATH}/build/lib/release.sh"


CMD_TARGETST=${SERVER_TARGETS[*]}

build::build_image
build::run_build_command make all TARGETS="${CMD_TARGETST}"

#release::build_server_images