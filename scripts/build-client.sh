#!/usr/bin/env bash
set -eu                   # exit on error or undefined variable
bash -c 'set -o pipefail' # return code of first cmd to fail in a pipeline

APP_DIR=${APP_DIR}

# Passing in version per illustration,
# best practise would be to use a version control tag,
# or something unique and readable like commit hash plus timestamp
VERSION=${VERSION}

cd "${APP_DIR}"

# Generate base64 config
./make.sh depends configu
configu -env "dev" -key "APP_VERSION" -value "${VERSION}"
CONFIG_BASE64=$(configu -base64)

# Build client
mkdir -p dist
rm -f "${APP_DIR}"/dist/client
go build -ldflags "-X main.configBase64=${CONFIG_BASE64}" \
  -o "${APP_DIR}"/dist/client ./cmd/client

# Instead of reloading server with new config, 
# write version and checksum to file
CHECKSUM=$("${APP_DIR}"/dist/client -checksum "${APP_DIR}"/dist/client)
echo "{\"version\": \"${VERSION}\", \"checksum\": \"${CHECKSUM}\"}" > \
  "${APP_DIR}"/dist/client.json

echo "done $(basename "$0")"
