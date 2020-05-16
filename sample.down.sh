#!/usr/bin/env bash
set -eu # exit on error or undefined variable
bash -c 'set -o pipefail' # return code of first cmd to fail in a pipeline

APP_NAME=${APP_NAME}

SESSIONS=(
    "app"
)
for (( i = 0 ; i < ${#SESSIONS[@]} ; i = i + 1 )) do
    SESSION="${APP_NAME}-${SESSIONS[$i]}"
    if tmux has-session -t ${SESSION} 2>/dev/null; then
        echo ${SESSION}
        tmux send-keys -t ${SESSION} C-c
        tmux send -t ${SESSION} "exit" ENTER
    fi
done

echo "done `basename $0`"