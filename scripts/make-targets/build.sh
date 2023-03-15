#!/usr/bin/env bash
#
# Build a static binary for the host OS/ARCH
#


set -eu

ROOT_PATH="$(dirname "$0")/../.."

. "${ROOT_PATH}/scripts/lib/init.sh"

golang_build "$@"

