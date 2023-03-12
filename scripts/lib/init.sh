#!/usr/bin/env sh


set -eu

ROOT_PATH=$(cd "$(dirname "$0")"/../.. && pwd -P)
OUTPUT_SUBPATH=${OUTPUT_SUBPATH:-"_output"}
OUTPUT_PATH="${ROOT_PATH}/${OUTPUT_SUBPATH}"
OUTPUT_BINPATH="${OUTPUT_PATH}/bin"
GO_MODULE_URL=$( grep module < go.mod | cut -d " " -f2)


export OUTPUT_PATH
export OUTPUT_BINPATH
export GO_MODULE_URL




. "${ROOT_PATH}/scripts/lib/golang.sh"



