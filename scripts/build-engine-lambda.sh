#!/usr/bin/env bash
set -x
PROJECT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )"/.. && pwd )"

rm -f "${PROJECT_DIR}"/build/engine-handler*

GOOS=linux go build -o "${PROJECT_DIR}"/build/engine-handler "${PROJECT_DIR}"/cmd/lambda/engine/main.go
pushd "${PROJECT_DIR}"/build || exit 128
zip engine-handler.zip engine-handler
popd || exit 128
