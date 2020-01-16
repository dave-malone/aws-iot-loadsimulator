#!/usr/bin/env bash
set -x
PROJECT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )"/.. && pwd )"

rm -f "${PROJECT_DIR}"/build/worker-handler*

GOOS=linux go build -o "${PROJECT_DIR}"/build/worker-handler "${PROJECT_DIR}"/cmd/lambda/worker/main.go
pushd "${PROJECT_DIR}"/build || exit 128
zip -r worker-handler.zip worker-handler
pushd "${PROJECT_DIR}" || exit 128
zip -ur "${PROJECT_DIR}"/build/worker-handler.zip certs
popd || exit 128
popd || exit 128
