#!/usr/bin/env bash
set -eu # exit on error or undefined variable
bash -c 'set -o pipefail' # return code of first cmd to fail in a pipeline

# Env
APP_DIR=${APP_DIR}

cd ${APP_DIR}
echo "build config cmd..."
go build -o ${APP_DIR}/config ./cmd/config

# Create config files if they don't exist
if [[ ! -f ${APP_DIR}/config.dev.json ]]; then
    echo "create dev config..."
    cp ${APP_DIR}/config.dev.sample.json ${APP_DIR}/config.dev.json
fi

echo "generate config helper..."
cd ${APP_DIR}
./config -generate ./pkg/config
go fmt ./pkg/config/config.go

