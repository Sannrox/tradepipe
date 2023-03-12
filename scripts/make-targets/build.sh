#!/usr/bin/env sh
#
# Build a static binary for the host OS/ARCH
#


set -eux

ROOT_PATH="$(dirname "$0")/../.."

. "${ROOT_PATH}/scripts/lib/init.sh"

golang_build "$@"

