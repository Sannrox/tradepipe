#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail
set -x 


ROOT_PATH=$(dirname "${BASH_SOURCE[0]}")/..

source "${ROOT_PATH}/build/general.sh"


build::build_image