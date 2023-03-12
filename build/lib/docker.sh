#!/usr/bin/env sh 

# shellcheck disable=SC2034

set -eu


USER_ID=$(id -u)
GROUP_ID=$(id -g)

DOCKER_OPTS=${DOCKER_OPTS:-}
DOCKER_CMD=${DOCKER:-docker}
DOCKER="${DOCKER_CMD} ${DOCKER_OPTS}"
DOCKER_HOST=${DOCKER_HOST:-}

DOCKER_IMAGE_TAG=${DOCKER_IMAGE_TAG:-"$(git rev-parse --short HEAD)"}

LOCAL_OUTPUT_ROOT="${ROOT_PATH}/_output"
LOCAL_OUTPUT_SUBPATH="${LOCAL_OUTPUT_ROOT}/dockerized"
LOCAL_OUTPUT_BINPATH="${LOCAL_OUTPUT_SUBPATH}/bin"
LOCAL_OUTPUT_GOPATH="${LOCAL_OUTPUT_SUBPATH}/go"
LOCAL_OUTPUT_IMAGES="${LOCAL_OUTPUT_SUBPATH}/images"

docker_available_on_osx(){
if [ -z "${DOCKER_HOST}" ]; then
    if [ -S "/var/run/docker.sock" ] || [ -S "$(docker context inspect --format  '{{.Endpoints.docker.Host}}' | awk -F 'unix://' '{print $2}')" ]; then
        printf "Docker is available on OSX"
      return 0
    fi
fi

    printf "Docker is not available on OSX"
    return 1
}

docker_check_if_in_path(){
  if ! command -v docker  > /dev/null 2>&1; then
    printf "Docker is not available in PATH"
    return 1
  fi
}

docker_image_exits(){
    if [ -z "$(docker images -q "${DOCKER_IMAGE_TAG}")" ]; then
        printf "Docker image ${DOCKER_IMAGE_TAG} does not exist"
        return 1
    fi
}

docker_delete_old_images() {
    for tag in $("$DOCKER" images "${1}" | tail -n +2 | awk '{print $2}'); do 
    if  echo "${tag}" | grep  -q "${2}*" ; then 
        print "Keeping image ${1}:${tag}"
    fi 
    if [ -z "${3:-}" ] || [ "${tag}" != "${3}" ]; then 
        print "Deleting image ${1}:${tag}"
        "$DOCKER" rmi "${1}:${tag}"  > /dev/null
    else 
        print "Keeping image ${1}:${tag}" 
    fi
    done
}

docker_delete_old_containers(){
    for container in $("$DOCKER" ps -a | tail -n +2 | awk '{print $1}'); do 
    if  echo "${container}" | grep  -q "${2}*" ; then 
        print "Keeping container ${container}"
        continue
    fi
    if [ -z "${3:-}" ] || [ "${container}" != "${3}" ]; then
        print "Deleting container ${container}"
        "$DOCKER" rm "${container}" > /dev/null
    else 
        print "Keeping container ${container}"
    fi
    done
}

docker_destroy_container(){
    "${DOCKER}" kill "${1}" > /dev/null 2>&1 || true
    "${DOCKER}" rm -f -v "${1}" > /dev/null 2>&1 || true
}

BINARY_PATH=${BINARY_PATH:-}
BASE_IMAGE="docker.io/library/golang:1.17.2-alpine3.14"
GIT_BRANCH=$(git symbolic-ref --short -q HEAD > /dev/null 2>&1 || true)
GIT_COMMIT_HASH=$(git rev-parse --short HEAD)
BUILD_IMAGE_REPO=${BUILD_IMAGE_REPO:-}
BUILD_CONTAINER_BASE_NAME=${BUILD_CONTAINER_BASE_NAME:-}
BUILD_IMAGE="${BUILD_IMAGE_REPO}/${BUILD_CONTAINER_BASE_NAME}:${GIT_COMMIT_HASH}"
LOCAL_OUTPUT_BUILD_CONTEXT="${LOCAL_OUTPUT_SUBPATH}/${BUILD_IMAGE}"


docker_clean(){
    docker_delete_old_images "${BUILD_IMAGE_REPO}/${BUILD_CONTAINER_BASE_NAME}" "${GIT_COMMIT_HASH}" "${GIT_COMMIT_HASH}"
    docker_delete_old_containers "${BUILD_IMAGE_REPO}/${BUILD_CONTAINER_BASE_NAME}" "${GIT_COMMIT_HASH}" "${GIT_COMMIT_HASH}"
    rm -rf "${LOCAL_OUTPUT_BUILD_CONTEXT}"
    print "Cleaned up docker build context and images"
}

docker_build_image(){
    mkdir -p "${LOCAL_OUTPUT_BUILD_CONTEXT}"

    cp "${ROOT_PATH}/Dockerfile" "${LOCAL_OUTPUT_BUILD_CONTEXT}/Dockerfile"

    docker_build "${BUILD_IMAGE}" \
     "${LOCAL_OUTPUT_BUILD_CONTEXT}" \
     "--build-arg USER_ID=${USER_ID} \
      --build-arg GROUP_ID=${GROUP_ID} \
      --build-arg BASE_IMAGE=${BASE_IMAGE} \
      --build-arg BINARY_PATH=${BINARY_PATH}"
    #Clean up old images
    docker_delete_old_images "${BUILD_IMAGE_REPO}/${BUILD_CONTAINER_BASE_NAME}" "${GIT_COMMIT_HASH}" "${GIT_COMMIT_HASH}"
    docker_delete_old_containers "${BUILD_IMAGE_REPO}/${BUILD_CONTAINER_BASE_NAME}" "${GIT_COMMIT_HASH}" "${GIT_COMMIT_HASH}"

}

docker_build() {
    image="${1}"
    context="${2}"
    docker_build_args="${3:-}"
    build_cmd="${DOCKER} build ${docker_build_args} -t "${image}" "${context}""

    echo "Building docker image ${image} with context ${context}"
    docker_ouput=$( eval $build_cmd) || {
        cat <<EOF >&2
    Error: Docker build failed:
    ${docker_ouput}

    to retry manually, run:

    ${build_cmd}

EOF
        return  1

}
}

build_get_docker_wrapped_binaries() {
    targets="tradegrpc,${GRPC_SERVER_IMAGE}\
        tradehttp,${HTTP_SERVER_IMAGE}"

    echo "${targets}"
}
