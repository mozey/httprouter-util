#!/usr/bin/env bash
set -eu # exit on error or undefined variable
bash -c 'set -o pipefail' # return code of first cmd to fail in a pipeline

EXPECTED_ARGS=1

if [[ $# -lt ${EXPECTED_ARGS} ]]
then
  echo "Usage:"
  echo "  `basename $0` TARGET"
  echo ""
  echo "Execute the specified target with dev env"
  echo ""
  echo "Examples:"
  echo "  `basename $0` app"
  echo "  `basename $0` test"
  exit 1
fi

TARGET=$1

# This script cannot change ENV of the parent process,
# it can only be set for child processes
export APP_ADDR=":8118"
export APP_EXE="dist/app"
export APP_DEV="true"
export APP_DIR=$(pwd)
export APP_NAME="httprouter-example"
export APP_PROXY="https://petstore.swagger.io/v2"

printenv | sort | grep --color -E "APP_|AWS_"

./make.sh ${TARGET}