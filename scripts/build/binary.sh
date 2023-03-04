#!/usr/bin/env sh
#
# Build a static binary for the host OS/ARCH
#


set -eu

ROOT_PATH="$(dirname "$0")/../.."

. "${ROOT_PATH}/scripts/lib/init"

golang_build "$@"






