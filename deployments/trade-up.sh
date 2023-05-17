#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

ROOT_PATH="$(dirname "$0")/.."

if [ -f "${ROOT_PATH}/.env" ]; then
    source "${ROOT_PATH}/.env"
fi

source "${ROOT_PATH}/deployments/general.sh"

# deployments::verify_prerequisites
deployments::compose::up "docker-compose.yml"