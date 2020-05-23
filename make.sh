#!/usr/bin/env bash
set -eu # exit on error or undefined variable
bash -c 'set -o pipefail' # return code of first cmd to fail in a pipeline

EXPECTED_ARGS=1

if [[ $# -ne ${EXPECTED_ARGS} ]]
then
  echo "Usage:"
  echo "  `basename $0` TARGET"
  echo ""
  echo "Execute the specified target"
  echo ""
  echo "Examples:"
  echo "  `basename $0` app"
  echo "  `basename $0` test"
  exit 1
fi

TARGET=${1}

# Executable to kill/restart
APP_EXE=${APP_EXE}

depends() {
    go version >/dev/null 2>&1 || \
	{ printf >&2 "Install https://golang.org\n"; exit 1; }
    fswatch --version >/dev/null 2>&1 || \
	{ printf >&2 "Install https://github.com/emcrisostomo/fswatch\n"; exit 1; }
}

# Build dev server
app_build_dev() {
    echo ${FUNCNAME}
    scripts/build.dev.sh
}

# Attempt to kill running server
app_kill() {
    echo ${FUNCNAME}
    # TODO On Windows use taskkill
    # https://docs.microsoft.com/en-us/windows-server/administration/windows-commands/taskkill
    killall -9 ${APP_EXE} 2>/dev/null || true
}

# Just run the server, no live reload
app_run() {
    echo ${FUNCNAME}
    app_kill
    app_build_dev; (if [[ "${?}" -eq 0 ]]; then (./${APP_EXE} ); fi)
}

# Restart server, for use with fswatch
app_restart() {
    echo ${FUNCNAME}
    app_kill
    app_build_dev; (if [[ "${?}" -eq 0 ]]; then (./${APP_EXE}& ); fi)
}

# Run app server with live reload
# Watch .go files for changes then recompile & try to start server
# will also kill server on ctrl+c
# fswatch includes everything unless an exclusion filter says otherwise
# https://stackoverflow.com/a/37237681/639133
app() {
    echo ${FUNCNAME}
    app_restart
    fswatch -or --exclude ".*" \
    --include "^.*pkg.*go$" \
    --include "./main.go$" \
    --include "./middleware.go$" ./ | \
	xargs -n1 bash -c "./make.sh app_restart" || bash -c "./make.sh app_kill"
}

TYPE=$(type -t ${TARGET} || echo "undefined")
if [[ ${TYPE} == "function" ]]; then
    ${TARGET}
else
    echo "${TARGET} not implemented"
    exit 1
fi
