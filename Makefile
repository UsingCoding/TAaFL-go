export APP_LL1_RUNNER_CMD_NAME = ll1-runner
export APP_LEXER_CMD_NAME = lexer
export DOCKER_IMAGE_NAME = vadimmakerov/$(APP_CMD_NAME):master

all: build test

.PHONY: build
build: modules
	bin/go-build.sh "cmd/$(APP_LL1_RUNNER_CMD_NAME)" "bin/$(APP_LL1_RUNNER_CMD_NAME)" $(APP_LL1_RUNNER_CMD_NAME)
	bin/cpp-build.sh "data/$(APP_LEXER_CMD_NAME)/main.cpp" "bin/$(APP_LEXER_CMD_NAME)"

.PHONY: modules
modules:
	go mod tidy

.PHONY: test
test:
	go test ./...

.PHONY: check
check:
	golangci-lint run

.PHONY: run-runner
run-runner:
	bin/$(APP_LL1_RUNNER_CMD_NAME) "-l" "bin/$(APP_LEXER_CMD_NAME)" "-g" "data/LL_1/grammar" "-i" "data/LL_1/program"

.PHONY: publish
publish:
	docker build . --tag=$(DOCKER_IMAGE_NAME)

.PHONY: clear
clear:
	rm -rf bin/$(APP_LL1_RUNNER_CMD_NAME)
	rm -rf bin/$(APP_LEXER_CMD_NAME)

.PHONY: build-dproxy
build-dproxy:
	docker build . -f data/docker/Dockerfile.proxy --tag=vadimmakerov/builder-docker-proxy