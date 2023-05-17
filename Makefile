# Version: 1.0


SHELL := /usr/bin/env bash -o errexit -o pipefail -o nounset
# We don't need make's built-in rules.
MAKEFLAGS += --no-builtin-rules
.SUFFIXES:





all:
	./scripts/make-targets/build.sh $(TARGETS)

test:
	./scripts/make-targets/test.sh

validate:
	scripts/make-targets/validate.sh

update:
	./scripts/make-targets/update.sh

clean:
	./scripts/make-targets/clean.sh

define CMD_HELP_INFO
# Add rules for all directories in cmd/
#
# Example:
#   make  tradepipe tradegear tradeapi
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


