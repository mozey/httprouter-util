#!/usr/bin/env bash
set -eu # exit on error or undefined variable
bash -c 'set -o pipefail' # return code of first cmd to fail in a pipeline

export APP_ADDR=":8118"
export APP_DEV="true"
export APP_DIR=$(pwd)
export APP_PROXY="https://petstore.swagger.io/v2"

printenv | grep APP_

make dev