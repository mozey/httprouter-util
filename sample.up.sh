#!/usr/bin/env bash
set -eu # exit on error or undefined variable
bash -c 'set -o pipefail' # return code of first cmd to fail in a pipeline

APP_DIR=${APP_DIR}
APP_NAME=${APP_NAME}

echo "creating tmux sessions..."
tmux new -d -s ${APP_NAME}-app

echo "set env in tmux..."
tmux send -t ${APP_NAME}-app "source ~/.bashrc && cd ${APP_DIR} && conf dev" ENTER

echo "start services..."
tmux send -t ${APP_NAME}-app "make reload" ENTER

echo "done `basename $0`"
