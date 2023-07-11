#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

ROOT_PATH="$(dirname "$0")/.."

source "${ROOT_PATH}/deployments/general.sh"

deployments::docker_load_images
