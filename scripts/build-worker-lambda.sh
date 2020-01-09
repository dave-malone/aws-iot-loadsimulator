#!/usr/bin/env bash
set -x
PROJECT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )"/.. && pwd )"

rm -f "${PROJECT_DIR}"/build/worker-handler*

GOOS=linux go build -o "${PROJECT_DIR}"/build/worker-handler "${PROJECT_DIR}"/cmd/lambda/worker/main.go
pushd "${PROJECT_DIR}"/build || exit 128
zip worker-handler.zip worker-handler "${PROJECT_DIR}"/certs/*
popd || exit 128
