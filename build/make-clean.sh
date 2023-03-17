#!/usr/bin/env bash

# shellcheck disable=SC2034

set -o errexit
set -o nounset
set -o pipefail

ROOT_PATH=$(dirname "$0")/..
source "${ROOT_PATH}/build/general.sh"

build::clean
