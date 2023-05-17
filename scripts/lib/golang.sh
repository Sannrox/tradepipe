#!/usr/bin/env bash

# shellcheck disable=SC2034

set -o errexit
set -o nounset
set -o pipefail

readonly GO_PACKAGE=github.com/Sannrox/tradepipe

readonly SERVER_PLATFORMS=(
    linux/amd64
    # linux/arm
    linux/arm64
)

GOOS="$(go env GOOS)"

GOARCH="$(go env GOARCH)"

if [ "${GOARCH}" = "arm" ]; then
	GOARM="$(go env GOARM)"
fi



PLATFORM=${PLATFORM:-}
VERSION=${VERSION:-$(git describe --tags --abbrev=0 || echo "1.0.0")}
GITCOMMIT=${GITCOMMIT:-$(git rev-parse --short HEAD 2> /dev/null || true)}
BUILDTIME=${BUILDTIME:-$(date -u +"%Y-%m-%dT%H:%M:%SZ")}



function golang::setup_environment(){
    export GOPATH="${GOPATH:-${HOME}/go}"

    export GOCACHE="${GOCACHE:-${HOME}/.cache/go-build}"
    export GOMODCACHE="${GOMODCACHE:-${HOME}/cache/go/mod}"

    export PATH="${GOPATH}/bin:${PATH}"

    GOROOT=$(go env GOROOT)

    unset GOBIN

}



function golang::server_targets(){
    local targets=(
        cmd/tradegear
        cmd/tradepipe
        cmd/tradeapi
        )

    echo "${targets[@]}"
}

IFS=" " read -r -a SERVER_TARGETS <<< "$(golang::server_targets)"
readonly SERVER_TARGETS
readonly SERVER_BINARIES=("${SERVER_TARGETS[@]##*/}")


function golang::build_binaries(){
    (
        local host_platform
        host_platform="$(go env GOOS)/$(go env GOARCH)"

        local -a platform
        IFS=" " read -r -a platform <<< "${BUILD_PLATFORMS:-}"
        if [[ ${#platform[@]} -eq 0 ]]; then
            platform=("${host_platform}")

        fi

        local -a targets=()
        for arg; do
            if [[ "${arg}" == -* ]]; then
                continue
            else
                targets+=("${arg}")
            fi
        done


        if [[ ${#targets[@]} -eq 0 ]]; then
            targets=("${SERVER_TARGETS[@]}")
        fi

        for platform in "${platform[@]}"; do
            golang::build_binaries_for_plattform "${platform}"
        done

    )

}


function golang::build_binaries_for_plattform() {
    local -a ldflags=()
    local platform_ldflags
    local -r platform="$1"
    local arch="${platform##*/}"
    local os="${platform%%/*}"

    platform_ldflags="-X \"main.PlatformName=${platform}\""
    ldflags="\
    -w \
    ${platform_ldflags} \
    -X \"main.GitCommit=${GITCOMMIT}\" \
    -X \"main.BuildTime=${BUILDTIME}\" \
    -X \"main.Version=${VERSION}\" \
    -X \"main.BuildArch=${arch}\" \
    -X \"main.BuildOs=${os}\" \
    ${LD_FLAGS:-} \
"

    for target in ${targets[@]}; do
        golang::build_binary "${target}"
    done
}


function golang::build_binary() {
    local -r target_path="$1"
    local -r target_name="${target_path##*/}"
    local -r target="${OUTPUT_BINPATH}/${platform}/${target_name}"
    local -r source="${GO_MODULE_URL}/${target_path}"

    : "${CGO_ENABLED=}"
    : "${GO_LINKMODE=static}"
    : "${GO_BUILDMODE=}"
    : "${GO_BUILDTAGS=}"
    : "${GO_STRIP=}"

    echo "Building $GO_LINKMODE ${target_name}"

    export GO111MODULE=auto

    build_cmd=(go build -o "${target}" -tags "${GO_BUILDTAGS}" --ldflags "${ldflags}" ${GO_BUILDMODE} "${source}" )

    build_cmd_output=$("${build_cmd[@]}" 2>&1) || {
        cat <<EOF >&2
Error building ${target_name}:
${build_cmd_output}
EOF
        exit 1
    }
    echo "Built ${target}"
}





