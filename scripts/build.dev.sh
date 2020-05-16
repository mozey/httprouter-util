#!/usr/bin/env bash
set -eu # exit on error or undefined variable
bash -c 'set -o pipefail' # return code of first cmd to fail in a pipeline

APP_DIR=${APP_DIR}
APP_EXE=${APP_EXE}

cd ${APP_DIR}
go build -o ${APP_DIR}/${APP_EXE} ./

echo "done `basename $0`"
