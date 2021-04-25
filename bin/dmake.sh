#!/bin/bash

set -o errexit

REPO_DIR="$(dirname "$(dirname "$(readlink -fm "$0")")")"
IMAGE_NAME="vadimmakerov/builder-docker-proxy"

pushd "$REPO_DIR" > /dev/null || exit 1

docker run \
  --rm \
  --name "TAaFL-go-builder" \
  --interactive \
  --tty \
  -v "$PWD":"$PWD" \
  -w "$PWD" \
  $IMAGE_NAME \
  "$@"

popd > /dev/null || exit 1
