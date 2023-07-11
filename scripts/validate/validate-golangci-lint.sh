#!/usr/bin/env bash

# shellcheck disable=SC2034
set -o errexit
set -o nounset
set -o pipefail



ROOT_PATH=$(dirname "${BASH_SOURCE[0]}")/../..
source "${ROOT_PATH}/scripts/lib/init.sh"
# GOLANG CI LINT
GOLANGCI_LINT=golangci-lint
GOLANGCI_LINT_OPTS=${GOLANGCI_LINT_OPTS:-}
GOLANGCI_LINT_VERSION="v1.52.0"

golang::setup_environment

function check_if_golangci_lint_is_in_path(){
    if ! type "${GOLANGCI_LINT}" > /dev/null; then
        install_golangci_lint
    fi

}

function check_golangci_lint_version(){
    if [ "$("${GOLANGCI_LINT}" version | grep -o "${GOLANGCI_LINT_VERSION/v/}")" != "${GOLANGCI_LINT_VERSION/v/}" ]; then
        echo "Install new version ${GOLANGCI_LINT_VERSION/v/} of golangci-lint "
        install_golangci_lint 
    fi
}

function run_golangci_lint(){
    golangci-lint run 
}

function install_golangci_lint(){
	mkdir -p "${GOPATH}/bin"
	curl -sfL "https://raw.githubusercontent.com/golangci/golangci-lint/${GOLANGCI_LINT_VERSION}/install.sh" \
		| sed -e '/install -d/d' \
		| sh -s -- -b "${GOPATH}/bin" "${GOLANGCI_LINT_VERSION}"
}


check_if_golangci_lint_is_in_path
check_golangci_lint_version
run_golangci_lint
