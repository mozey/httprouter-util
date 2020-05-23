#!/usr/bin/env bash
set -eu # exit on error or undefined variable
bash -c 'set -o pipefail' # return code of first cmd to fail in a pipeline

APP_DIR=${APP_DIR}
echo "APP_DIR=${APP_DIR}"

read -r -p "Reset dev configuration in APP_DIR? [y/N] " response
if [[ "$response" =~ ^([yY][eE][sS]|[yY])$ ]]
then
    rm -f ${APP_DIR}/*.out
    rm -f ${APP_DIR}/build/*
    rmdir build
    rm -f ${APP_DIR}/config
    rm -f ${APP_DIR}/config.dev.json
    rm -f ${APP_DIR}/dev.sh
    rm -f ${APP_DIR}/up.sh
    echo "done"
else
    echo "abort"
fi


