export APP_CMD_NAME = astbackend

APP_EXECUTABLE_OUT?=bin
BINARY?="${APP_EXECUTABLE_OUT}/${APP_CMD_NAME}"

DENO=deno compile

.PHONY: build
build: init-binary-file
	${DENO} -o ${BINARY} data/astbackend/app.js

.PHONY: init-binary-file
init-binary-file:
	mkdir -p ${APP_EXECUTABLE_OUT}