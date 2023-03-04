#!/usr/bin/env sh

set -eu

GOOS="$(go env GOOS)"

GOARCH="$(go env GOARCH)"

if [ "${GOARCH}" = "arm" ]; then
	GOARM="$(go env GOARM)"
fi

PLATFORM=${PLATFORM:-}
PLATFORM_LDFLAGS=
if test -n "${PLATFORM}"; then
	PLATFORM_LDFLAGS="-X \"main.PlatformName=${PLATFORM}\""
fi


PLATFORM=${PLATFORM:-}
VERSION=${VERSION:-$(git describe --tags --abbrev=0 || echo "1.0.0")}
GITCOMMIT=${GITCOMMIT:-$(git rev-parse --short HEAD 2> /dev/null || true)}
BUILDTIME=${BUILDTIME:-$(date -u +"%Y-%m-%dT%H:%M:%SZ")}

LDFLAGS="\
    -w \
    ${PLATFORM_LDFLAGS} \
    -X \"main.GitCommit=${GITCOMMIT}\" \
    -X \"main.BuildTime=${BUILDTIME}\" \
    -X \"main.Version=${VERSION}\" \
    -X \"main.BuildArch=${GOARCH}\" \
    -X \"main.BuildOs=${GOOS}\" \
    ${LDFLAGS:-} \
"


golang_build() {
    TARGET_PATH=$@
    NAME=$(basename "$@")
    TARGET="${OUTPUT_BINPATH}/${NAME}-${GOOS}-${GOARCH}"
    SOURCE="${GO_MODULE_URL}/${TARGET_PATH}"

    : "${CGO_ENABLED=}"
    : "${GO_LINKMODE=static}"
    : "${GO_BUILDMODE=}"
    : "${GO_BUILDTAGS=}"
    : "${GO_STRIP=}"

    echo "Building $GO_LINKMODE $(basename "${TARGET}")"

    export GO111MODULE=auto

    go build -o "${TARGET}" -tags "${GO_BUILDTAGS}" --ldflags "${LDFLAGS}" ${GO_BUILDMODE} "${SOURCE}"

    echo ">> build ${TARGET}"
}



