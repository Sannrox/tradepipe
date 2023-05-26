#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

ROOT_PATH="$(dirname "$0")/../.."

ROOT_PATH=$(cd $(dirname "${BASH_SOURCE[0]}")/.. && pwd )
DEPLOYMENTS_PATH="${ROOT_PATH}/deployments"
ENV_VARS_PATH="${DEPLOYMENTS_PATH}/env_vars"

readonly SECRETS=()

function deployments::compose::up() {
    local -r compose_file="${DEPLOYMENTS_PATH}/$1"
    docker-compose -f "${compose_file}" up -d
}

function deployments::compose::down() {
    local -r compose_file="${DEPLOYMENTS_PATH}/$1"
    docker-compose -f "${compose_file}" down
}

function deployments::create_secret(){
      local secret_name="${1}"
      local secret_value="${2}"

      echo "${secret_value}" > "${ENV_VARS}/.${secret_name}"
}

function deployments::create_env_vars_directory(){
    mkdir -p "${ENV_VARS_PATH}"
}

function deployments::create_secrets(){
    for secret in "${SECRETS[@]}"; do
        if [[ -z "${!secret}" ]]; then
            echo "Please set ${secret} in .env file"
            exit 1
        else
            deployments::create_secret "${secret}" "${!secret}"
        fi
    done
}

function deployments::verify_prerequisites(){
    if ! [[ -d "${ENV_VARS_PATH}" ]]; then
        deployments::create_env_vars_directory
    fi
        deployments::create_secrets
}


function deployments::clean(){
    rm -rf "${ROOT_PATH}/deployments/env_vars"
    rm -rf "${ROOT_PATH}/data"
}

function deployments::docker_clean_simple(){
    docker system prune --force
    docker volume prune --force
}


function deployments::docker_load_images(){
    deployments::docker_clean_simple

    for target in $(deployments::get_docker_containers); do
        if [ -f "${ROOT_PATH}/_output/release-images/amd64/${target}.tar" ]; then
           docker load -i "${ROOT_PATH}/_output/release-images/amd64/${target}.tar"
        else
           echo "Image ${target} not available in tar format"
           return 1
        fi
    done
}

function deployments::get_docker_containers(){
    targets=("tradegear")
    echo "${targets[@]}"
}
