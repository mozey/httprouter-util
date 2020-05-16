#!/usr/bin/env bash
set -eu # exit on error or undefined variable
bash -c 'set -o pipefail' # return code of first cmd to fail in a pipeline

EXPECTED_ARGS=1

if [[ $# -lt ${EXPECTED_ARGS} ]]
then
  echo "Usage:"
  echo "  `basename $0` FN"
  echo ""
  echo "Attach to tmux for specified fn"
  echo ""
  echo "Examples:"
  echo "  `basename $0` app"
  exit 1
fi

FN=${1}
APP_NAME=${APP_NAME}

tmux attach -t "${APP_NAME}-${FN}"
