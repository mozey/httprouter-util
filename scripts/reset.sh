#!/usr/bin/env bash

# Set (e) exit on error
# Set (u) no-unset to exit on undefined variable
set -eu
# If any command in a pipeline fails,
# that return code will be used as the
# return code of the whole pipeline.
bash -c 'set -o pipefail'

APP_DIR=${APP_DIR}
echo "APP_DIR=${APP_DIR}"

read -r -p "Reset dev configuration in APP_DIR? [y/N] " response
if [[ "$response" =~ ^([yY][eE][sS]|[yY])$ ]]
then
    rm -f ${APP_DIR}/config
    rm -f ${APP_DIR}/config.dev.json
    rm -f ${APP_DIR}/*.out
    rm -f ${APP_DIR}/dev.sh
    rm -f ${APP_DIR}/up.sh
    echo "done"
else
    echo "abort"
fi


