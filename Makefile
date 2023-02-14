# Version: 1.0


.PHONY: test lint deps format build coverage coverhtml grpc rest 


GO ?= $(shell which go)

ifneq ($(GO), )
	GOFMT ?= $(GO)fmt
	GOLIST ?= $(GO) list
	GOPATH := $(firstword $(subst :, ,$(shell $(GO) env GOPATH)))
	GOHOSTOS     ?= $(shell $(GO) env GOHOSTOS)
	GOHOSTARCH   ?= $(shell $(GO) env GOHOSTARCH)
	GOTESTSUM    ?= $(shell which gotestsum)
	PKGS          = ./...


# GO TEST
GOTEST := $(GO) test
endif
GOTEST_DIR :=
ifneq (${CI},)
ifneq ($(GOTESTSUM),)
	GOTEST_DIR := test-results
	GOTEST := gotestsum --junitfile $(CURRENT_DIR)/$(GOTEST_DIR)/unit-tests.xml --
endif
endif


ifeq ($(GOHOSTARCH),amd64)
        ifeq ($(GOHOSTOS),$(filter $(GOHOSTOS),linux freebsd darwin windows))
                # Only supported on amd64
                test-flags := -race
        endif
endif


ifneq ($(CI),)
ifdef ($(GOTESTSUM),)
$(GOTESTSUM):
	$(GO) get gotest.tools/gotestsum
endif
endif



$(GOTEST_DIR):
	@mkdir -p $@

test: $(GOTESTSUM) $(GOTEST_DIR) ## Run unit-tests
	@echo ">> running test for core"
	./scripts/test/pre-test-steps
	CGO_ENABLED=1 $(GOTEST) $(test-flags) $(GOOPTS) $(PKGS)

# GOLANG CI LINT
GOLANGCI_LINT := $(shell which golangci-lint)
GOLANGCI_LINT_OPTS ?=
GOLANGCI_LINT_VERSION ?= v1.39.0


ifeq ($(GOHOSTOS),$(filter $(GOHOSTOS),linux darwin))
	ifeq ($(GOHOSTARCH),$(filter $(GOHOSTARCH),amd64 i386))
		GOLANGCI_LINT := $(GOPATH)/bin/golangci-lint
	endif
endif

ifdef GOLANGCI_LINT
$(GOLANGCI_LINT):
	mkdir -p $(GOPATH)/bin
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/$(GOLANGCI_LINT_VERSION)/install.sh \
		| sed -e '/install -d/d' \
		| sh -s -- -b $(GOPATH)/bin $(GOLANGCI_LINT_VERSION)
endif


deps: ## Get the dependencies
	@echo ">> getting dependencies for core"
	$(GO) get $(GOOPTS) -t ./...
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	go install mvdan.cc/gofumpt@latest
	go install github.com/daixiang0/gci@latest
	@if ! type "protoc" > /dev/null; then \
        wget https://github.com/protocolbuffers/protobuf/releases/download/v3.7.1/protoc-3.7.1-linux-x86_64.zip; \
        unzip protoc-3.7.1-linux-x86_64.zip -d protoc3; \
        sudo mv protoc3/bin/* /usr/local/bin/; \
        sudo mv protoc3/include/* /usr/local/include/; \
	fi
fmt: ## Formates the code
	@echo ">> formatting code"
	gofumpt -l -w .
	gci -w .
	$(GOLIST) -f {{.Dir}} ./... | xargs $(GOFMT) -w -s -d


build: ## Build the binary
	$(eval TARGET ?=  )
	@echo ">> building binaries"
	@if [ -z "$(TARGET)" ] ; then \
		./scripts/build/binary ;\
	else \
		TARGET=$(TARGET) ./scripts/build/binary ;\
	fi


coverage: ## Generate global code coverage report
	./scripts/test/coverage;

coverhtml: ## Generate global code coverage report in HTML
	./scripts/test/coverage html;

changelog: ## Generates changelog for last updates
	./scripts/utils/log

lint: ## run all the lint tools
	$(GOLANGCI_LINT) run
rest: ## Generate go code from openapi spec
	./scripts/generate/rest
grpc: ## Generate go code from protobuf files
	./scripts/generate/grpc
help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

clean:
	rm -rf "./tmp" "./timelineEventsWithDocs.json" "./timelineEventsWithoutDocs.json" "./tradepip.txt"

