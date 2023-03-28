MAKEFLAGS+="k"

GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
RESET  := $(shell tput -Txterm sgr0)

include .env
export $(shell sed 's/=.*//' .env)

.PHONY: help
help:
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@echo "  ${YELLOW}help            ${RESET} Show this help message"
	@echo "  ${YELLOW}build           ${RESET} Build application binary"
	@echo "  ${YELLOW}setup           ${RESET} Setup local environment"
	@echo "  ${YELLOW}check           ${RESET} Run tests, linters and tidy of the project"
	@echo "  ${YELLOW}test            ${RESET} Run tests only"
	@echo "  ${YELLOW}lint            ${RESET} Run linters via golangci-lint"
	@echo "  ${YELLOW}tidy            ${RESET} Run tidy for go module to remove unused dependencies"
	@echo "  ${YELLOW}run             ${RESET} Run application"

.PHONY: build
build:
	OS="$(OS)" APP="$(APP)" ./hacks/build.sh

.PHONY: setup
setup:
	bash -c ./hacks/setup.sh

.PHONY: check
check: %: tidy lint test

.PHONY: test
test:
	TEST_RUN_ARGS="$(TEST_RUN_ARGS)" TEST_DIR="$(TEST_DIR)" ./hacks/run-tests.sh

.PHONY: lint
lint:
	golangci-lint run

.PHONY: tidy
tidy:
	go mod tidy -v

.PHONY: run
run:
	cd ./cmd/app/ && go run .
