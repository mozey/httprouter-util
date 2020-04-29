#!/usr/bin/env bash

# Set (e) exit on error
# Set (u) no-unset to exit on undefined variable
set -eu
# If any command in a pipeline fails,
# that return code will be used as the
# return code of the whole pipeline.
bash -c 'set -o pipefail'

# Env
APP_DIR=${APP_DIR}

# Build config util
cd ${APP_DIR}
go build \
-ldflags "-X main.AppDir=${APP_DIR}" \
-o ${APP_DIR}/config ./cmd/config

# Create config files if they don't exist
if [[ ! -f ${APP_DIR}/config.dev.json ]]; then
    echo "Create dev config"
    cp ${APP_DIR}/config.dev.sample.json ${APP_DIR}/config.dev.json
    ${APP_DIR}/config -key APP_DIR -value ${APP_DIR}
fi

# Generate config helper
cd ${APP_DIR}
./config -generate ./pkg/config
go fmt ./pkg/config/config.go

