#!/usr/bin/env bash
set -eu                   # exit on error or undefined variable
bash -c 'set -o pipefail' # return code of first cmd to fail in a pipeline

# Env
APP_DIR=${APP_DIR}
GOPATH=${GOPATH}

# Safety check for working on multiple projects
if [[ "${APP_DIR}" != $(pwd) ]]; then
  echo "To avoid clobbering files in other projects run this in APP_DIR, e.g."
  echo "  APP_DIR=\$(pwd) ./script/config.sh"
  exit 1
fi
cd ${APP_DIR}

if ! test -f "${GOPATH}/bin/configu"; then
  echo "Install https://github.com/mozey/config"
  exit 1
fi

# Create config files if they don't exist
if [[ ! -f ${APP_DIR}/config.dev.json ]]; then
  echo "create dev config..."
  cp ${APP_DIR}/sample.config.dev.json ${APP_DIR}/config.dev.json
fi

echo "generate config helper..."
cd ${APP_DIR}
${GOPATH}/bin/configu -generate ./pkg/config
go fmt ./pkg/config/config.go

echo "done `basename $0`"