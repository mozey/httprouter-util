#!/usr/bin/env bash
set -eu                   # exit on error or undefined variable
bash -c 'set -o pipefail' # return code of first cmd to fail in a pipeline

EXPECTED_ARGS=1

if [[ $# -lt ${EXPECTED_ARGS} ]]; then
  echo "Usage:"
  echo "  $(basename "$0") FUNC [ARGS...]"
  echo ""
  echo "Execute the specified func"
  echo ""
  echo "Examples:"
  echo "  $(basename "$0") depends go"
  exit 1
fi

FUNC=${1}

# depends checks for programs this script depends on
depends() {
  if [[ ${1} == "go" ]]; then
    go version >/dev/null 2>&1 ||
      {
        echo "Install https://golang.org"
        exit 1
      }

  elif [[ ${1} == "watcher" ]]; then
    "${GOPATH}"/bin/watcher -version >/dev/null 2>&1 ||
      {
        echo "Install https://github.com/mozey/watcher"
        exit 1
      }

  elif [[ ${1} == "gojq" ]]; then
    gojq --version >/dev/null 2>&1 ||
      {
        echo "Install https://github.com/itchyny/gojq"
        exit 1
      }

  elif [[ ${1} == "gotest" ]]; then
    # TODO How to make gotest print version?
    TYPE=$(type -t gotest || echo "undefined")
    if [[ ${TYPE} != "file" ]]; then
      {
        echo "Install https://github.com/rakyll/gotest"
        exit 1
      }
    fi

  elif [[ ${1} == "APP_EXE_PATH" ]]; then
    export APP_EXE_PATH=${APP_EXE_PATH:-undefined}
    if [[ ${APP_EXE_PATH} == "undefined" ]]; then
      # Error if APP_EXE is not set
      export APP_EXE=${APP_EXE}
      # Use full path to avoid conflicts
      # shellcheck disable=2155
      export APP_EXE_PATH="$(pwd)/${APP_EXE}"
    fi

  else
    echo "unknown dependency ${1}"
    exit 1
  fi
}

# detect_os is useful when writing cross platform functions.
# Return value must corrrespond to GOOS listed here
# https://go.dev/doc/install/source#environment
# shellcheck disable=2120
detect_os() {
  # Check if func arg is set
  # https://stackoverflow.com/a/13864829/639133
  if [ -z ${1+x} ]; then
    OUTPUT=$(uname -s)
  else
    # This is useful when parsing output from a remote host
    OUTPUT="$1"
  fi
  case "$OUTPUT" in
  Darwin)
    echo 'darwin'
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

# detect_arch is useful when writing cross platform functions
# Return value must corrrespond to GOARCH listed here
# https://go.dev/doc/install/source#environment
# shellcheck disable=2120
detect_arch() {
  # Check if func arg is set
  # https://stackoverflow.com/a/13864829/639133
  if [ -z ${1+x} ]; then
    OUTPUT=$(uname -m)
  else
    # This is useful when parsing output from a remote host
    OUTPUT="$1"
  fi
  case "$OUTPUT" in
  amd64)
    echo 'amd64'
    ;;
  x86_64)
    echo 'amd64'
    ;;
  arm64)
    echo 'arm64'
    ;;
  *)
    echo 'other'
    ;;
  esac
}

# Kill process by matching full path to bin
kill_path() {
  OS=$(detect_os)
  if [[ ${OS} == "darwin" ]] || [[ ${OS} == "linux" ]]; then
    PID=$(pgrep -fx "${1}" || echo "")
    if [[ -n "${PID}" ]]; then
      kill "${PID}"
    fi
  else
    echo "OS ${OS} not implemented"
    exit 1
  fi
}

app_build_dev() {
  scripts/build-dev.sh
}

app_kill() {
  depends APP_EXE_PATH
  kill_path "${APP_EXE_PATH}"
}

# Run the binary, no live reload.
# Use full path to avoid conflicts
app_run() {
  depends go
  depends APP_EXE_PATH
  app_kill
  app_build_dev
  # shellcheck disable=2181
  (if [[ "${?}" -eq 0 ]]; then (${APP_EXE_PATH}); fi)
}

# Restart the binary.
# Use full path to avoid conflicts
app_restart() {
  depends go
  depends APP_EXE_PATH
  app_kill
  app_build_dev
  # shellcheck disable=2181
  (if [[ "${?}" -eq 0 ]]; then (${APP_EXE_PATH} &) fi)
}

# Run app bin with live reload
# Watch .go files for changes then recompile & try to start bin
# will also kill bin on ctrl+c
app() {
  depends watcher
  app_restart
  APP_DIR=${APP_DIR}
  "${GOPATH}"/bin/watcher -d 1500 -r -dir "" \
    --include ".*.go$" \
    --excludeDir "${APP_DIR}/cmd.*" \
    --excludeDir "${APP_DIR}/dist.*" \
    --excludeDir "${APP_DIR}/scripts.*" \
    --excludeDir "${APP_DIR}/vendor.*" \
    --excludeDir "${APP_DIR}/www.*" |
    xargs -n1 bash -c "./make.sh app_restart" || bash -c "./make.sh app_kill"
}

# env_sh generates ENV.sh from config.ENV.json
env_sh() {
  depends gojq
  ENV=${1}
  PROMPT=${2:-undefined}

  if [[ ! -d "$APP_DIR" ]]; then
    echo "Dir $APP_DIR does not exist"
    exit 1
  fi

  if [[ ${PROMPT} != "PROMPT_DISABLED" ]]; then
    read -r -p "Generate new ${ENV}.sh? [y/N] " response
    if [[ "$response" =~ ^([yY][eE][sS]|[yY])$ ]]; then
      :
    else
      echo "abort"
      exit 1
    fi
  fi

  SAMPLE_CONFIG="$APP_DIR/sample.config.${ENV}.json"
  CONFIG="$APP_DIR/config.${ENV}.json"
  ENV_SH="$APP_DIR/${ENV}.sh"

  if [[ -f "${CONFIG}" ]]; then
    echo "Using existing ${CONFIG}"
  else
    echo "Creating new ${CONFIG}"
    cp "${SAMPLE_CONFIG}" "${CONFIG}"
  fi

  echo "#!/usr/bin/env bash" >"${ENV_SH}"
  {
    echo "# ------------------------ DO NOT EDIT!"
    echo "# Generated from $CONFIG"
    echo "# Edit the JSON file to change your local config,"
    echo "# then refresh $APP_DIR/$ENV.sh like this"
    echo "#     ./make.sh env_sh $APP_DIR $ENV"
    echo ""
    echo "if [ -z \${APP_DIR+x} ]; then"
    echo "    # Caller can override \$APP_DIR"
    echo "    export APP_DIR=$APP_DIR"
    echo "fi"
    echo ""
  } >>"${ENV_SH}"

  # Create array of key values
  # https://stackoverflow.com/a/23118607/639133
  KEY_VALUES=()
  while IFS='' read -r line; do
    KEY_VALUES+=("$line")
  done < <(gojq -r 'to_entries|map("\(.key)\n\(.value|tostring)")|.[]' "${CONFIG}")
  for ((i = 0; i < ${#KEY_VALUES[@]}; i = i + 2)); do
    KEY="${KEY_VALUES[$i]}"
    VALUE="${KEY_VALUES[$i + 1]}"
    echo "export ${KEY}=\"${VALUE}\"" >>"${ENV_SH}"
  done

  echo "" >>"${ENV_SH}"
  echo 'printenv | sort | grep --color -E "APP_|AWS_"' >>"${ENV_SH}"

  chmod u+x "${ENV_SH}"
}

# ..............................................................................

# Required env vars
APP_DIR=${APP_DIR}
GOPATH=${GOPATH}

# Execute FUNC if it's a func defined in this script
TYPE=$(type -t "${FUNC}" || echo "undefined")
if [[ ${TYPE} == "function" ]]; then
  # Additional arguments, after the func, are passed through.
  # For example, `./make.sh FUNC ARG1`
  ${FUNC} "${@:2}"
else
  echo "FUNC ${FUNC} not implemented"
  exit 1
fi
