export APP_CMD_NAME = lexer

APP_EXECUTABLE_OUT?=bin
BINARY?="${APP_EXECUTABLE_OUT}/${APP_CMD_NAME}"

CLANG=clang++-6.0 -std=c++17

.PHONY: build
build: init-binary-file
	${CLANG} -o ${BINARY} ./data/lexer/main.cpp

.PHONY: init-binary-file
init-binary-file:
	mkdir -p ${APP_EXECUTABLE_OUT}
	touch ${BINARY}