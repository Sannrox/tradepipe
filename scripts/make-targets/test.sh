#!/usr/bin/env sh

set -eux

CI=${CI:-}
GO_TEST=${GO_TEST:-}
GO_TEST_DIR=${GO_TEST_DIR:-}
GO_TEST_FLAGS=${GO_TEST_FLAGS:-}
GO_TESTSUM=${GO_TESTSUM:-}
GO_OPTS=${GO_OPTS:-}
GOHOSTARCH=${GOHOSTARCH:-}
PKGS="./..."



try_ci_test_run(){
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

check_if_gotestsum_is_in_path(){
    if type "gotestsum" > /dev/null; then
        GOTESTSUM="gotestsum"
    fi
}

install_gotestsum() {
	go get gotest.tools/gotestsum
}

check_test_flags(){
    if [ -n "${GOHOSTARCH}" ]; then
        if [ "${GOHOSTARCH}" = "amd64" ]; then 
            case "${GOHOSTOS}" in 
                linux)
                GO_TEST_FLAGS="-race"
                ;;
                freebsd)
                GO_TEST_FLAGS="-race"
                ;;
                darwin)
                GO_TEST_FLAGS="-race"
                ;;
                windows)
                GO_TEST_FLAGS="-race"
                ;;
            esac
        fi
    fi
}


create_test_dir(){
    if [ -n "${GO_TEST_DIR}" ]; then
        mkdir -p "${GO_TEST_DIR}"
    fi
}

run_test(){
	CGO_ENABLED=1 ${GO_TEST} ${GO_TEST_FLAGS} ${GO_OPTS} ${PKGS}
}




try_ci_test_run
check_test_flags
create_test_dir
run_test


