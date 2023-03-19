#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail


ROOT_PATH=$(dirname "${BASH_SOURCE[0]}")/../..
source "${ROOT_PATH}/scripts/lib/init.sh"

CI=${CI:-}
GO_TEST_TIMEOUT=${GO_TEST_TIMEOUT:-"-timeout=180s"}
GO_TEST=${GO_TEST:-}
GO_TEST_DIR=${GO_TEST_DIR:-}
GO_TESTSUM=${GO_TESTSUM:-}
GO_TEST_FLAGS=${GO_TEST_FLAGS:-}
GO_OPTS=${GO_OPTS:-}
GOHOSTARCH=${GOHOSTARCH:-}
PKGS="./..."

golang::setup_environment

function try_ci_test_run(){
    if [ -n "${CI}" ]; then
        if [ -n "${GOTESTSUM}" ]; then
            check_if_golangci_lint_is_in_path 
            if [ -z "${GOTESTSUM}" ]; then
                echo "gotestsum not found in PATH"
                install_gotestsum
            fi
            GOTEST_DIR := test-results
            GOTEST := gotestsum --junitfile "${GO_TEST_DIR}/unit-tests.xml" --
        fi
    else
        GO_TEST="go test"
    fi
}

function check_if_gotestsum_is_in_path(){
    if type "gotestsum" > /dev/null; then
        GOTESTSUM="gotestsum"
    fi
}

function install_gotestsum() {
	go get gotest.tools/gotestsum
}


function check_if_race(){
if [ -n "${GOHOSTARCH}" ]; then
    if [ "${GOHOSTARCH}" = "amd64" ]; then 
        case "${GOHOSTOS}" in 
            linux)
            goflags+=(-race)
            ;;
            freebsd)
            goflags=(-race)
            ;;
            darwin)
            goflags=(-race)
            ;;
            windows)
            goflags=(-race)
            ;;
        esac
    fi
fi

}


function create_test_dir(){
    if [ -n "${GO_TEST_DIR}" ]; then
        mkdir -p "${GO_TEST_DIR}"
    fi
}

function run_tests(){
    check_if_race
    if [ -n "${GO_TEST}" ]; then
       GO_ENABLED=1  ${GO_TEST} ${GO_TEST_FLAGS} ${GO_TEST_TIMEOUT} ${GO_OPTS} ${PKGS}
    fi
}

try_ci_test_run
create_test_dir
run_tests