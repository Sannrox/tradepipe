#!/usr/bin/env bash 

# shellcheck disable=SC2034

set -o errexit
set -o nounset
set -o pipefail



USER_ID=$(id -u)
GROUP_ID=$(id -g)


DOCKER_OPTS=${DOCKER_OPTS:-}
IFS=" " read -r -a DOCKER <<< "docker ${DOCKER_OPTS}"
DOCKER_HOST=${DOCKER_HOST:-}

ROOT_PATH=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)

source "${ROOT_PATH}/scripts/lib/init.sh"


readonly BASE_IMAGE="golang:1.20.2-bullseye"
readonly BUILD_IMAGE_REPO="tradepipe"

readonly LOCAL_OUTPUT_ROOT="${ROOT_PATH}/_output"
readonly LOCAL_OUTPUT_BINPATH="${LOCAL_OUTPUT_ROOT}/bin"
readonly LOCAL_OUTPUT_IMAGES="${LOCAL_OUTPUT_ROOT}/images"


readonly REMOTE_ROOT="/go/src/${GO_PACKAGE}"
readonly REMOTE_OUTPUT_ROOT="${REMOTE_ROOT}/_output"
readonly REMOTE_OUTPUT_BINPATH="${REMOTE_OUTPUT_ROOT}/bin"


readonly GRPC_SERVER_BASE_IMAGE="${GRPC_SERVER_BASE_IMAGE:-$BASE_IMAGE}"
readonly HTTP_SERVER_BASE_IMAGE="${HTTP_SERVER_BASE_IMAGE:-$BASE_IMAGE}"


GIT_BRANCH=$(git symbolic-ref --short -q HEAD > /dev/null 2>&1 || true)
GIT_COMMIT_HASH=$(git rev-parse --short HEAD)
BUILD_IMAGE_BASE_TAG="build-${GIT_COMMIT_HASH}"
BUILD_IMAGE="${BUILD_IMAGE_REPO}:${BUILD_IMAGE_BASE_TAG}"
BUILD_CONTAINER_BASE_NAME="tradepipe-build"
BUILD_CONTAINER_NAME="${BUILD_CONTAINER_BASE_NAME}-${GIT_COMMIT_HASH}"
DATA_CONTAINER_BASE_NAME="build-data"
DATA_CONTAINER_NAME="${DATA_CONTAINER_BASE_NAME}-${GIT_COMMIT_HASH}"

DOCKER_MOUNT_AGRS=(--volumes-from "${DATA_CONTAINER_NAME}" --volume "${ROOT_PATH}:${REMOTE_ROOT}")
LOCAL_OUTPUT_BUILD_CONTEXT="${LOCAL_OUTPUT_ROOT}/${BUILD_IMAGE}"



function build::available_on_osx(){
if [ -z "${DOCKER_HOST}" ]; then
    if [ -S "/var/run/docker.sock" ] || [ -S "$(docker context inspect --format  '{{.Endpoints.docker.Host}}' | awk -F 'unix://' '{echo $2}')" ]; then
        echof "Docker is available on OSX"
      return 0
    fi
fi

    echof "Docker is not available on OSX"
    return 1
}

function build::check_docker_if_in_path(){
  if ! command -v docker  > /dev/null 2>&1; then
    echof "Docker is not available in PATH"
    return 1
  fi
}

function build::docker_image_exits(){
    [[ -n "$1" && -n "$2" ]] || {
        echo "Internal error. Image not specified"
        exit 2
    }

    [[ $("${DOCKER[@]}" images -q "${1}:${2}") ]]
}

function build::docker_delete_old_images() {
    for tag in $("${DOCKER[@]}" images "${1}" | tail -n +2 | awk '{echo $2}'); do 
    if  [[ "${tag}" != "${2}*" ]] ; then 
        echo "Keeping image ${1}:${tag}"
        continue
    fi 
    if [[ -z "${2:-}"  ||  "${tag}" != "${3}" ]]; then 
        echo "Deleting image ${1}:${tag}"
        "${DOCKER[@]}" rmi "${1}:${tag}"  > /dev/null
    else 
        echo "Keeping image ${1}:${tag}" 
    fi
    done
}

function build::docker_delete_old_containers(){
    for container in $("${DOCKER[@]}" ps -a | tail -n +2 | awk '{echo $NF}'); do 
    if  [[ "${container}" != "${1}*" ]] ; then 
        echo "Keeping container ${container}"
        continue
    fi
    if [[ -z "${2:-}"  ||  "${container}" != "${2}" ]]; then
        echo "Deleting container ${container}"
        build::docker_destroy_container "${container}"
    else 
        echo "Keeping container ${container}"
    fi
    done
}

function build::docker_destroy_container(){
    "${DOCKER}" kill "${1}" > /dev/null 2>&1 || true
    "${DOCKER}" rm -f -v "${1}" > /dev/null 2>&1 || true
}


function build::clean(){
    build::docker_delete_old_containers "${BUILD_CONTAINER_BASE_NAME}" 
    build::docker_delete_old_containers "${DATA_CONTAINER_BASE_NAME}"
    build::docker_delete_old_images "${BUILD_IMAGE_REPO}"

    echo "Cleaning all untaged images"
    "${DOCKER[@]}" rmi $("${DOCKER[@]}" images -q --filter "dangling=true") > /dev/null 2>&1 || true


    if [[ -d "${LOCAL_OUTPUT_ROOT}" ]]; then
        echo "Cleaning ${LOCAL_OUTPUT_ROOT}"
        rm -rf "${LOCAL_OUTPUT_ROOT}"
    fi
}

function build::docker_clean(){
    docker_delete_old_images "${BUILD_IMAGE_REPO}/${BUILD_CONTAINER_BASE_NAME}" "${GIT_COMMIT_HASH}" "${GIT_COMMIT_HASH}"
    docker_delete_old_containers "${BUILD_IMAGE_REPO}/${BUILD_CONTAINER_BASE_NAME}" "${GIT_COMMIT_HASH}" "${GIT_COMMIT_HASH}"
    rm -rf "${LOCAL_OUTPUT_BUILD_CONTEXT}"
    echo "Cleaned up docker build context and images"
}

function build::load_data_container(){
    # If the data container already exists, we don't need to do anything
    local ret=0
    local code=0

    code=$(docker inspect --format='{{.State.Running}}' "${DATA_CONTAINER_NAME}" 2>/dev/null || ret=$?)

    if [[ "${ret}" -eq 0 && "${code}" != 0 ]]; then
        build::docker_destroy_container "${DATA_CONTAINER_NAME}"
        ret=1 
    fi

    if [[ "${ret}" -ne 0 ]]; then 
        echo "Creating data container ${DATA_CONTAINER_NAME}"


        local -ra docker_run_cmd=(
            "${DOCKER[@]}" run
            --volume "${REMOTE_ROOT}"   # white-out the whole output dir
            --volume /usr/local/go/pkg/linux_386_cgo
            --volume /usr/local/go/pkg/linux_amd64_cgo
            --volume /usr/local/go/pkg/linux_arm_cgo
            --volume /usr/local/go/pkg/linux_arm64_cgo
            --volume /usr/local/go/pkg/linux_ppc64le_cgo
            --volume /usr/local/go/pkg/darwin_amd64_cgo
            --volume /usr/local/go/pkg/darwin_386_cgo
            --volume /usr/local/go/pkg/windows_amd64_cgo
            --volume /usr/local/go/pkg/windows_386_cgo
            --name "${DATA_CONTAINER_NAME}"
            --hostname "${HOSTNAME}"
            "${BUILD_IMAGE}"
            chown -R "$(id -u):$(id -g)" "${REMOTE_ROOT}" /usr/local/go/pkg
        )

        "${docker_run_cmd[@]}" 
    fi
}

function build::build_image(){

    mkdir -p "${LOCAL_OUTPUT_BUILD_CONTEXT}"

    cp "${ROOT_PATH}/build/build-image/Dockerfile" "${LOCAL_OUTPUT_BUILD_CONTEXT}/Dockerfile"


    build::docker_build "${BUILD_IMAGE}" "${LOCAL_OUTPUT_BUILD_CONTEXT}" "false" "--build-arg BASE_IMAGE=${BASE_IMAGE}"
    #Clean up old images
    build::docker_delete_old_images "${BUILD_IMAGE_REPO}/${BUILD_CONTAINER_BASE_NAME}" "${GIT_COMMIT_HASH}" "${GIT_COMMIT_HASH}"
    build::docker_delete_old_containers "${BUILD_IMAGE_REPO}/${BUILD_CONTAINER_BASE_NAME}" "${GIT_COMMIT_HASH}" "${GIT_COMMIT_HASH}"

    build::load_data_container

}

function build::run_build_command(){
    echo "Running build command: $@"
    build::run_build_command_in_build_image "${BUILD_CONTAINER_NAME}" -- "$@"
}

function build::run_build_command_in_build_image(){
    [[ $# != 0  ]] || { echo "Internal error. No image specified" >&2; return 4; }
    local -r container_name="${1}"
    shift

    local -a docker_run_opts=(
        "--name=${container_name}"
        "--user=$(id -u):$(id -g)"
        "--hostname=${HOSTNAME}"
        "${DOCKER_MOUNT_AGRS[@]}"
    )

    local detach=false

    [[ $# != 0  ]] || { echo "Invalid input - docker args followed by --" >&2; return 4; }
    # Everything before the -- is passed to the docker run command
    until [[ -z "${1-}" ]]; do
        if [[ "${1}" == "--" ]]; then
            shift
            break
        fi
        docker_run_opts+=("${1}")
        if [[ "${1}" == "-d" || "${1}" == "--detach" ]]; then
            detach=true
        fi
        shift
    done

    # Everything after the -- is passed to the docker command
    [[ $# != 0  ]] || { echo "Invalid input - command to run not given" >&2; return 4; }
    local -a docker_cmd=()
    until [[ -z "${1-}" ]]; do
        docker_cmd+=("${1}")
        shift
    done

    # To add more options to the docker run command, add them to the array above
    docker_run_opts+=(
        --env "BUILD_PLATFORMS=${BUILD_PLATFORMS:-}"
    )


    if [[ -t 0 ]]; then 
        docker_run_opts+=(--interactive --tty)
    elif [[ "${detach}" == "false" ]]; then
        docker_run_opts+=(--attach=stdout --attach=stderr)
    fi

    local -ra docker_run_cmd=(${DOCKER[@]} run "${docker_run_opts[@]}" "${BUILD_IMAGE}")
    #Clean up old containers
    build::docker_destroy_container "${container_name}"
    "${docker_run_cmd[@]}" "${docker_cmd[@]}"
    if [[ "${detach}" == "false" ]]; then
        build::docker_destroy_container "${container_name}"
    fi
}

function build::docker_build() {
    local -r image="${1}"
    local -r context_dir="${2}"
    local -r pull="${3:-true}"
    local build_args 
    IFS=" " read -r -a build_args <<< "$4"
    
    local -ra build_cmd=("${DOCKER[@]}" buildx build "${build_args[@]}" -t "${image}" "--pull=${pull}" "${context_dir}")

    echo "Building docker image ${image} with context ${context_dir}"
    docker_ouput=$(DOCKER_CLI_EXPERIMENTAL=enabled "${build_cmd[@]}" 2>&1) || {
        cat <<EOF >&2
    Error: Docker build failed:
    ${docker_ouput}

    to retry manually, run:

    DOCKER_CLI_EXPERIMENTAL=enabled ${build_cmd[@]}

EOF
        return  1

}
}

function build::get_docker_wrapped_binaries() {
    targets="tradegrpc,${GRPC_SERVER_BASE_IMAGE}\
        tradehttp,${HTTP_SERVER_BASE_IMAGE}"

    echo "${targets}"
}

