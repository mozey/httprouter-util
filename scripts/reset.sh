#!/usr/bin/env bash
set -eu # exit on error or undefined variable
bash -c 'set -o pipefail' # return code of first cmd to fail in a pipeline

APP_DIR=${APP_DIR}
echo "APP_DIR=${APP_DIR}"

read -r -p "Reset dev config in APP_DIR? [y/N] " response
if [[ "$response" =~ ^([yY][eE][sS]|[yY])$ ]]
then
    rm -f ${APP_DIR}/*.out
    rm -f ${APP_DIR}/dist/*
    rm -f ${APP_DIR}/config
    rm -f ${APP_DIR}/config.dev.json
    rm -f ${APP_DIR}/dev.sh
    rm -f ${APP_DIR}/up.sh
    rm -f ${APP_DIR}/down.sh
    echo "done `basename $0`"
else
    echo "abort"
fi


