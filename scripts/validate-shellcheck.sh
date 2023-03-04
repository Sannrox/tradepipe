#!/usr/bin/env sh

set -eux

ROOT_PATH=$(cd "$(dirname "$0")"/.. && pwd -P)

. "${ROOT_PATH}/scripts/lib/init.sh"


SHELLCHECK_VERSION="0.8.0"
SHELLCHECK_IMAGE="docker.io/koalaman/shellcheck-alpine:v0.8.0@sha256:f42fde76d2d14a645a848826e54a4d650150e151d9c81057c898da89a82c8a56"


all_files() {
  find . -type f -name '*.sh' -not -path './vendor/*' -not -path './_output/*'
}


SHELLCHECK=0
if type "shellcheck" > /dev/null; then 
  detected_version="$(shellcheck --version | grep 'version: .*')"
  if [ "${detected_version}" = "version: ${SHELLCHECK_VERSION}" ]; then
    SHELLCHECK=1
  fi
fi

if [ "${SHELLCHECK}" == 1 ]; then
  echo "Using host shellcheck ${SHELLCHECK_VERSION} binary."
  shellcheck $(all_files)
fi