export APP_FACTORIZER_CMD_NAME = factorizer
export APP_GENERATOR_CMD_NAME = generator
export APP_RUNNER_CMD_NAME = runner
export APP_LEXER_CMD_NAME = lexer
export DOCKER_IMAGE_NAME = vadimmakerov/$(APP_CMD_NAME):master

all: build test

.PHONY: build
build: modules
	bin/go-build.sh "cmd/$(APP_FACTORIZER_CMD_NAME)" "bin/$(APP_FACTORIZER_CMD_NAME)" $(APP_FACTORIZER_CMD_NAME)
	bin/go-build.sh "cmd/$(APP_GENERATOR_CMD_NAME)" "bin/$(APP_GENERATOR_CMD_NAME)" $(APP_GENERATOR_CMD_NAME)
	bin/go-build.sh "cmd/$(APP_RUNNER_CMD_NAME)" "bin/$(APP_RUNNER_CMD_NAME)" $(APP_RUNNER_CMD_NAME)
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

.PHONY: run-factorizer
run-factorizer: build
	bin/$(APP_FACTORIZER_CMD_NAME)

.PHONY: run-generator
run-generator:
	bin/$(APP_GENERATOR_CMD_NAME)

.PHONY: run-runner
run-runner:
	bin/$(APP_RUNNER_CMD_NAME) "-l" "bin/$(APP_LEXER_CMD_NAME)" "-f" "data/LL(1)/grammar" "-g" "data/LL(1)/program"

.PHONY: publish
publish:
	docker build . --tag=$(DOCKER_IMAGE_NAME)

.PHONY: clear
clear:
	rm -rf bin/$(APP_FACTORIZER_CMD_NAME)
	rm -rf bin/$(APP_GENERATOR_CMD_NAME)
	rm -rf bin/$(APP_RUNNER_CMD_NAME)
	rm -rf bin/$(APP_LEXER_CMD_NAME)

.PHONY: build-dproxy
build-dproxy:
	docker build . -f Dockerfile.proxy --tag=vadimmakerov/builder-docker-proxy