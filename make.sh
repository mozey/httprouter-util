#!/usr/bin/env bash
set -eu                   # exit on error or undefined variable
bash -c 'set -o pipefail' # return code of first cmd to fail in a pipeline

EXPECTED_ARGS=1

if [[ $# -lt ${EXPECTED_ARGS} ]]; then
  echo "Usage:"
  echo "  $(basename $0) TARGET"
  echo ""
  echo "Execute the specified target"
  echo ""
  echo "Examples:"
  echo "  $(basename $0) app"
  echo "  $(basename $0) test"
  exit 1
fi

GOPATH=${GOPATH}
TARGET=${1}

# Depends lists, and can be used to check for, programs this script depends on
depends() {
  if [[ ${1} == "go" ]]; then
    go version >/dev/null 2>&1 ||
      {
        printf >&2 \
          "Install https://golang.org\n"
        exit 1
      }

  elif [[ ${1} == "watcher" ]]; then
    ${GOPATH}/bin/watcher -version >/dev/null 2>&1 ||
      {
        printf >&2 \
          "Install https://github.com/mozey/watcher\n"
        exit 1
      }

  elif [[ ${1} == "jq" ]]; then
    jq --version >/dev/null 2>&1 ||
      {
        printf >&2 \
          "Install https://stedolan.github.io/jq\n"
        exit 1
      }

  elif [[ ${1} == "gotest" ]]; then
    # TODO How to make gotest print version?
    TYPE=$(type -t gotest || echo "undefined")
    if [[ ${TYPE} != "file" ]]; then
      {
        printf >&2 \
          "Install https://github.com/rakyll/gotest\n"
        exit 1
      }
    fi

  elif [[ ${1} == "APP_EXE_PATH" ]]; then
    export APP_EXE_PATH=${APP_EXE_PATH:-undefined}
    if [[ ${APP_EXE_PATH} == "undefined" ]]; then
      # Error if APP_EXE is not set
      export APP_EXE=${APP_EXE}
      # Use full path to avoid conflicts
      export APP_EXE_PATH="$(pwd)/${APP_EXE}"
    fi

  else
    echo "unknown dependency ${1}"
    exit 1
  fi
}

detect_os() {
  case "$(uname -s)" in
  Darwin)
    echo 'macOS'
    ;;
  Linux)
    echo 'linux'
    ;;
  CYGWIN* | MINGW32* | MSYS* | MINGW*)
    echo 'windows'
    ;;
  # Detect additional OS's here...
  # See correspondence table at the bottom of this link
  # https://stackoverflow.com/a/27776822/639133
  *)
    echo 'other'
    ;;
  esac
}

# Kill process by matching full path to bin
kill_path() {
  OS=$(detect_os)
  if [[ ${OS} == "macOS" ]] || [[ ${OS} == "linux" ]]; then
    PID=$(pgrep -fx "${1}" || echo "")
    if [[ -n "${PID}" ]]; then
      kill ${PID}
    fi
  else
    echo "OS ${OS} not implemented"
    exit 1
  fi
}

app_build_dev() {
  echo ${FUNCNAME}
  scripts/build.dev.sh
}

app_kill() {
  echo ${FUNCNAME}
  depends APP_EXE_PATH
  kill_path ${APP_EXE_PATH}
}

# Run the binary, no live reload.
# Use full path to avoid conflicts
app_run() {
  echo ${FUNCNAME}
  depends go
  depends APP_EXE_PATH
  app_kill
  app_build_dev
  (if [[ "${?}" -eq 0 ]]; then (${APP_EXE_PATH}); fi)
}

# Restart the binary.
# Use full path to avoid conflicts
app_restart() {
  echo ${FUNCNAME}
  depends go
  depends APP_EXE_PATH
  app_kill
  app_build_dev
  (if [[ "${?}" -eq 0 ]]; then (${APP_EXE_PATH} &); fi)
}

# Run app bin with live reload
# Watch .go files for changes then recompile & try to start bin
# will also kill bin on ctrl+c
app() {
  echo ${FUNCNAME}
  depends watcher
  app_restart
  APP_DIR=${APP_DIR}
  ${GOPATH}/bin/watcher -d 1500 -r -dir "" \
    --include ".*.go$" \
    --excludeDir "${APP_DIR}/cmd.*" \
    --excludeDir "${APP_DIR}/dist.*" \
    --excludeDir "${APP_DIR}/scripts.*" \
    --excludeDir "${APP_DIR}/vendor.*" \
    --excludeDir "${APP_DIR}/www.*" |
    xargs -n1 bash -c "./make.sh app_restart" || bash -c "./make.sh app_kill"
}

# Generate dev.sh from config.dev.json
dev_sh() {
  echo ${FUNCNAME}
  depends jq

  read -r -p "Generate new dev.sh? [y/N] " response
  if [[ "$response" =~ ^([yY][eE][sS]|[yY])$ ]]; then
    :
  else
    echo "abort"
    exit 0
  fi

  SAMPLE_CONFIG="sample.config.dev.json"
  CONFIG="config.dev.json"
  DEV_SH="dev.sh"

  if [[ -f "${CONFIG}" ]]; then
    echo "using existing config file"
  else
    echo "creating new ${CONFIG}"
    cp ${SAMPLE_CONFIG} ${CONFIG}
  fi

  echo "#!/usr/bin/env bash" >${DEV_SH}
  echo "# Generated from ${CONFIG} with make.sh, DO NOT EDIT!" >>${DEV_SH}
  echo "" >>${DEV_SH}
  echo "export APP_DIR=$(pwd)" >>${DEV_SH}
  echo "" >>${DEV_SH}

  # Create array of key values
  # https://stackoverflow.com/a/23118607/639133
  KEY_VALUES=()
  while IFS='' read -r line; do
    KEY_VALUES+=("$line")
  done < <(jq -r 'to_entries|map("\(.key)\n\(.value|tostring)")|.[]' ./${CONFIG})
  for ((i = 0; i < ${#KEY_VALUES[@]}; i = i + 2)); do
    KEY="${KEY_VALUES[$i]}"
    VALUE="${KEY_VALUES[$i + 1]}"
    echo "export ${KEY}=\"${VALUE}\"" >>${DEV_SH}
  done

  echo "" >>${DEV_SH}
  echo 'printenv | sort | grep --color -E "APP_|AWS_"' >>${DEV_SH}

  chmod u+x ${DEV_SH}
}

# Execute target if it's a func defined in this script.
TYPE=$(type -t ${TARGET} || echo "undefined")
if [[ ${TYPE} == "function" ]]; then
  # Additional arguments, after the target, are passed through.
  # For example, `./make.sh depends something`
  ${TARGET} ${@:2}
else
  echo "TARGET ${TARGET} not implemented"
  exit 1
fi
