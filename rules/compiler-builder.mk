export APP_CMD_NAME = compiler

APP_EXECUTABLE_OUT?=bin
BINARY?="${APP_EXECUTABLE_OUT}/${APP_CMD_NAME}"

STATIC_FLAGS=CGO_ENABLED=0

GO_BUILD=$(STATIC_FLAGS) go build -trimpath -v

ifdef DEBUG
	GO_BUILD+="--gcflags=\"all=-N -l\""
	BINARY := "${BINARY}-debug"
endif

.PHONY: build
build:
	${GO_BUILD} -o ${BINARY} ./cmd/${APP_CMD_NAME}

.PHONY: modules
modules:
	go mod tidy

.PHONY: check
check:
	golangci-lint run
