#!/usr/bin/env bash

# shellcheck disable=SC2034

set -o errexit
set -o nounset
set -o pipefail
set -x

readonly RELEASE_STAGE_PATH="${LOCAL_OUTPUT_ROOT}/release-stage"
readonly RELEASE_IMAGES_PATH="${LOCAL_OUTPUT_ROOT}/release-images"

BUILD_PULL_LATEST_IMAGES=${BUILD_PULL_LATEST_IMAGES:-0}


function release::build_server_images(){
    rm -rf "${RELEASE_IMAGES_PATH}"

    local platform 
    for platform in "${SERVER_PLATFORMS[@]}"; do
        local platform_tag arch
        platform_tag=${platform/\//-} 
        arch=$(basename "${platform}")

        echo "Building release images for ${platform_tag} (${arch})"

        local releas_stage 
        releas_stage="${RELEASE_STAGE_PATH}/servers/${platform_tag}/tradepipe"
        rm -rf "${releas_stage}"
        mkdir -p "${releas_stage}/server/bin"

        cp "${SERVER_BINARIES[@]/#/${LOCAL_OUTPUT_BINPATH}/${platform}/}" "${releas_stage}/server/bin/"

        release::create_server_images "${releas_stage}/server/bin" "${arch}"


    done
}

function release::create_server_images(){
    binary_dir="$1"
    arch="$2"
    binaries="$(build::get_docker_wrapped_binaries)"
    images_dir="${RELEASE_IMAGES_PATH}/${arch}"

    mkdir -p "${images_dir}"

    docker_registry="${DOCKER_REGISTRY:-}"

     docker_tag="${GIT_VERSION}"
     if [ -z "${docker_tag}" ]; then
        echo "GIT_VERSION is not set; cannot create docker images"
        return 1
     fi

     docker_build_opts=
     if [ "${BUILD_PULL_LATEST_IMAGES}" -eq 1 ]; then
        docker_build_opts="${docker_build_opts} --pull"
     fi

     for wrapped in $binaries; do 
        binary_name="${wrapped%%,*}"
        base_image="${wrapped##*,}"
        binary_path="${binary_dir}/${binary_name}"
        docker_build_path="${binary_dir}.dockerbuild"
        docker_image_tag="${docker_registry}${binary_name}:${docker_tag}-${arch}"
        docker_file_path="${ROOT_PATH}/build/server-image/Dockerfile"

        if [ -f "${ROOT_PATH}/build/image/${binary_name}/Dockerfile" ]; then
            docker_file_path="${ROOT_PATH}/build/image/${binary_name}/Dockerfile"
        fi

        echo "Building docker image ${docker_image_tag} from ${docker_file_path} with context ${docker_build_path}"
        (
        rm -rf "${docker_build_path}"
        mkdir -p "${docker_build_path}"

        ln  "${binary_path}" "${docker_build_path}/${binary_name}"

        local build_log="${docker_build_path}/build.log"

        if ! DOCKER_CLI_EXPERIMENTAL=enabled "${DOCKER[@]}" buildx build  \
            --load ${docker_build_opts:+"${docker_build_opts}"} \
            -t "${docker_image_tag}" \
            -f "${docker_file_path}" \
            --platform linux/"${arch}" \
            --build-arg BINARY_NAME="${binary_name}" \
            --build-arg BASE_IMAGE="${base_image}" \
            "${docker_build_path}" > "${build_log}" 2>&1; then
            cat "${build_log}"
            exit 1
        fi
        rm "${build_log}"

        echo "Created docker image ${docker_image_tag}"


        "${DOCKER[@]}" save -o "${binary_path}.tar" "${docker_image_tag}"
        echo "${docker_tag}" >> "${binary_path}.docker_tag"
        rm -rf "${docker_build_path}" 
        ln "${binary_path}.tar" "${images_dir}/"
        "${DOCKER[@]}" rmi "${docker_image_tag}" &> /dev/null || true
        ) &
    
    done



}