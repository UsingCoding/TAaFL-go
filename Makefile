export APP_FACTORIZER_CMD_NAME = factorizer
export DOCKER_IMAGE_NAME = vadimmakerov/$(APP_CMD_NAME):master

all: build check test

.PHONY: build
build: modules
	bin/go-build.sh "cmd/$(APP_FACTORIZER_CMD_NAME)" "bin/$(APP_FACTORIZER_CMD_NAME)" $(APP_FACTORIZER_CMD_NAME)

.PHONY: modules
modules:
	go mod tidy

.PHONY: test
test:
	go test ./...

.PHONY: check
check:
	golangci-lint run

.PHONY: run
run-factorizer: build
	bin/$(APP_FACTORIZER_CMD_NAME)

.PHONY: publish
publish:
	docker build . --tag=$(DOCKER_IMAGE_NAME)