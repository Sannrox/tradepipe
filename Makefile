# Version: 1.0


# We don't need make's built-in rules.
MAKEFLAGS += --no-builtin-rules
.SUFFIXES:

.PHONY: test lint deps format build coverage coverhtml grpc rest


all:
	./scripts/make-targets/build.sh $(TARGET)

GO ?= $(shell which go)

ifneq ($(GO), )
	GOFMT ?= $(GO)fmt
	GOLIST ?= $(GO) list
endif



deps: ## Get the dependencies
	@echo ">> getting dependencies for core"
	$(GO) get $(GOOPTS) -t ./...
	go install mvdan.cc/gofumpt@latest
	go install github.com/daixiang0/gci@latest

fmt: ## Formates the code
	@echo ">> formatting code"
	gofumpt -l -w .
	gci -w .
	$(GOLIST) -f {{.Dir}} ./... | xargs $(GOFMT) -w -s -d


define CMD_HELP_INFO
# Add rules for all directories in cmd/
#
# Example:
#   make  tradepipe tradegrpc tradehttp
endef
EXCLUDE_TARGET=
CMD_TARGET = $(notdir $(abspath $(wildcard cmd/*/)))
.PHONY: $(CMD_TARGET)
ifeq ($(PRINT_HELP),y)
$(CMD_TARGET): ## $(CMD_TARGET)
	echo "$$CMD_HELP_INFO"
else
$(CMD_TARGET): ## $(CMD_TARGET)
	@echo ">> building $@"
	./scripts/make-targets/build.sh cmd/$@
endif

test: ## Run unittests
	@echo ">> running unit tests"
	./scripts/make-targets/test.sh

validate: ## Validate the code
	@echo ">> validating code"
	scripts/make-targets/validate.sh

clean:
	@echo ">> clean up"
	./scripts/make-targets/clean.sh


help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

