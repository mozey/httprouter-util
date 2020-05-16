#!/usr/bin/env bash
set -eu # exit on error or undefined variable
bash -c 'set -o pipefail' # return code of first cmd to fail in a pipeline

EXPECTED_ARGS=1

if [[ $# -lt ${EXPECTED_ARGS} ]]
then
  echo "Usage:"
  echo "  `basename $0` MODE"
  echo ""
  echo "Run dev server"
  echo ""
  echo "Examples:"
  echo "  `basename $0` run"
  echo "  `basename $0` reload"
  exit 1
fi

MODE=$1

export APP_ADDR=":8118"
export APP_EXE="app.out"
export APP_DEV="true"
export APP_DIR=$(pwd)
export APP_NAME="httprouter-example"
export APP_PROXY="https://petstore.swagger.io/v2"

printenv | grep APP_

make ${MODE}