#!/usr/bin/env bash
set -eu                   # exit on error or undefined variable
bash -c 'set -o pipefail' # return code of first cmd to fail in a pipeline

APP_DIR=${APP_DIR}
APP_CLIENT_VERSION=${APP_CLIENT_VERSION}

cd ${APP_DIR}
mkdir -p dist
rm -f ${APP_DIR}/dist/client
go build -ldflags "-X main.version=${APP_CLIENT_VERSION}" \
  -o ${APP_DIR}/dist/client ./cmd/client
echo ${APP_CLIENT_VERSION} > ${APP_DIR}/dist/client.version

echo "done $(basename $0)"
