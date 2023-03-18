#!/usr/bin/env bash


set -o errexit
set -o nounset
set -o pipefail


unset CDPATH

ROOT_PATH=$(cd "$(dirname "${BASH_SOURCE[0]}")"/../.. && pwd -P)
OUTPUT_SUBPATH=${OUTPUT_SUBPATH:-"_output"}
OUTPUT_PATH="${ROOT_PATH}/${OUTPUT_SUBPATH}"
OUTPUT_BINPATH="${OUTPUT_PATH}/bin"
GO_MODULE_URL=$( grep module < go.mod | cut -d " " -f2)


export OUTPUT_PATH
export OUTPUT_BINPATH
export GO_MODULE_URL





source "${ROOT_PATH}/scripts/lib/golang.sh"
source "${ROOT_PATH}/scripts/lib/version.sh"
source "${ROOT_PATH}/scripts/lib/openapi.sh"
source "${ROOT_PATH}/scripts/lib/protoc.sh"



readlinkdashf (){
  # run in a subshell for simpler 'cd'
  (
    if [[ -d "${1}" ]]; then # This also catch symlinks to dirs.
      cd "${1}"
      pwd -P
    else
      cd "$(dirname "${1}")"
      local f
      f=$(basename "${1}")
      if [[ -L "${f}" ]]; then
        readlink "${f}"
      else
        echo "$(pwd -P)/${f}"
      fi
    fi
  )
}