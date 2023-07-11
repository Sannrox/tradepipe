#!/usr/bin/env bash
#
# Build a static binary for the host OS/ARCH
#


set -o errexit
set -o nounset
set -o pipefail


ROOT_PATH="$(dirname "$0")/../.."

. "${ROOT_PATH}/scripts/lib/init.sh"

golang::build_binaries "$@"

