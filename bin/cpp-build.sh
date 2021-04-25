#!/bin/bash

set -o errexit

if [ ! -d "$(dirname "$1")" ] || [ ! -d "$(dirname "$2")" ] ; then
    echo "usage: $(basename "$0") <path-to-main> <path-to-output-file>" 1>&2
    exit 1
fi

MAIN_FILE=$1
EXECUTABLE_PATH=$2

echo_call() {
    echo "$@"
    "$@"
}

# touch file cause linker cannot create file
touch "$EXECUTABLE_PATH"
echo_call clang++-6.0 -std=c++17 -o "$EXECUTABLE_PATH" "$MAIN_FILE"