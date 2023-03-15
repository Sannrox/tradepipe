#!/usr/bin/env sh

# shellcheck disable=SC2034
set -eu

# GOLANG CI LINT
GOLANGCI_LINT=${GOLANGCI_LINT:-}
GOLANGCI_LINT_OPTS=${GOLANGCI_LINT_OPTS:-}
GOLANGCI_LINT_VERSION="1.39.0"


check_if_golangci_lint_is_in_path(){
    if type "golangci-lint" > /dev/null; then
        GOLANGCI_LINT="golangci-lint"
    fi

}

check_golangci_lint_version(){
    if [ "$("${GOLANGCI_LINT}" version | grep -o "${GOLANGCI_LINT_VERSION}")" != "${GOLANGCI_LINT_VERSION}" ]; then
        print "Install new version of golangci-lint"
        install_golangci_lint 
    fi
}

run_golangci_lint(){
    "${GOLANGCI_LINT}" run 
}

install_golangci_lint(){
	mkdir -p "${GOPATH}/bin"
	curl -sfL "https://raw.githubusercontent.com/golangci/golangci-lint/v${GOLANGCI_LINT_VERSION}/install.sh" \
		| sed -e '/install -d/d' \
		| sh -s -- -b "${GOPATH}/bin" "${GOLANGCI_LINT_VERSION}"
}


check_if_golangci_lint_is_in_path
if [ -z "${GOLANGCI_LINT}" ]; then
    echo "golangci-lint not found in PATH"
    install_golangci_lint
fi

check_golangci_lint_version
run_golangci_lint
