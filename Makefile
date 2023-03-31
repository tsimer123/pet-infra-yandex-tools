MAKEFLAGS+="k"

GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
RESET  := $(shell tput -Txterm sgr0)

MAKEFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
CURRENT_DIR := $(dir $(MAKEFILE_PATH))
PROJECT_NAME = pet-infra-yandex-tools
REPO_NAME = tsimer123
DOCKER_IMAGE_COMMIT_SHA=$(shell git show -s --format=%h)
DOCKER_IMAGE_REPO = ghcr.io/${REPO_NAME}/${PROJECT_NAME}

BUILD_OS := $(shell uname | sed 's/./\L&/g')
BUILD_ARCH := $(shell uname -m)
ifeq ($(BUILD_ARCH),x86_64)
	BUILD_ARCH = amd64
endif
ifeq ($(BUILD_ARCH),aarch64)
	BUILD_ARCH = arm64
endif

include .env
export $(shell sed 's/=.*//' .env)

.PHONY: help
help:
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@echo "  ${YELLOW}help                      ${RESET} Show this help message"
	@echo "  ${YELLOW}test                      ${RESET} Run tests only"
	@echo "  ${YELLOW}lint                      ${RESET} Run linters via golangci-lint"
	@echo "  ${YELLOW}tidy                      ${RESET} Run tidy for go module to remove unused dependencies"
	@echo "  ${YELLOW}job-list                  ${RESET} View all implemented jobs"
	@echo "  ${YELLOW}job-build                 ${RESET} Build specific job locally. Requires argument JOB"
	@echo "  ${YELLOW}job-build-all             ${RESET} Build specific job locally. Requires argument JOB"
	@echo "  ${YELLOW}job-build-docker          ${RESET} Build specific job locally and create docker image. Requires argument JOB"
	@echo "  ${YELLOW}job-build-docker-all      ${RESET} Build all jobs locally and create docker images"
	@echo "  ${YELLOW}job-push-docker           ${RESET} Push specific job to docker registry. Requires argument JOB"
	@echo "  ${YELLOW}job-push-docker-all       ${RESET} Push all jobs to docker registry"
	@echo "  ${YELLOW}job-run                   ${RESET} Run specific job locally. Requires argument JOB"
	@echo "  ${YELLOW}job-run-docker            ${RESET} Run specific job in docker. Requires argument JOB"

.PHONY: test
test:
	TEST_RUN_ARGS="$(TEST_RUN_ARGS)" TEST_DIR="$(TEST_DIR)" ./hacks/run-tests.sh

.PHONY: lint
lint:
	golangci-lint run

.PHONY: tidy
tidy:
	go mod tidy -v

.PHONY: job-list
job-list:
	@find cmd -mindepth 1 -maxdepth 1 -exec basename {} \; | sort

.PHONY: job-build
job-build:
	mkdir -p ./cmd/${JOB}/bin/config && \
	cp ./cmd/${JOB}/config/config.json ./cmd/${JOB}/bin/config && \
	GOOS=$(BUILD_OS) GOARCH=$(BUILD_ARCH) CGO_ENABLED=0 go build -v -o ./cmd/${JOB}/bin/${JOB} ./cmd/${JOB}

.PHONY: job-build-all
job-build-all: BUILD_OS = linux
job-build-all:
	make -s job-list | xargs -I % make -s JOB=% job-build 2>/dev/null

.PHONY: job-build-docker
job-build-docker: BUILD_OS = linux
job-build-docker: job-build
	docker build \
	--progress plain \
	--platform linux/${BUILD_ARCH} \
	--build-arg JOB_NAME=${JOB} \
	--tag "${DOCKER_IMAGE_REPO}/${JOB}:${DOCKER_IMAGE_COMMIT_SHA}" \
	--file Dockerfile \
	.

.PHONY: job-build-docker-all
job-build-docker-all: BUILD_OS = linux
job-build-docker-all:
	make -s job-list | xargs -I % make -s JOB=% job-build-docker 2>/dev/null

.PHONY: job-push-docker
job-push-docker: BUILD_OS = linux
job-push-docker: job-build-docker
	docker push "${DOCKER_IMAGE_REPO}/${JOB}:${DOCKER_IMAGE_COMMIT_SHA}"

.PHONY: job-push-docker-all
job-push-docker-all: BUILD_OS = linux
job-push-docker-all:
	make -s job-list | xargs -I % make -s JOB=% job-push-docker 2>/dev/null

.PHONY: job-run
job-run: job-build
	@pushd >/dev/null ./cmd/${JOB}/bin && \
 	./${JOB} && \
 	popd >/dev/null

.PHONY: job-run-docker
job-run-docker: job-build-docker
	@docker run \
	--entrypoint "./${JOB}" \
	-it \
	--rm "${DOCKER_IMAGE_REPO}/${JOB}:${DOCKER_IMAGE_COMMIT_SHA}"
