#!/bin/sh

set -e

# shellcheck disable=SC2028
echo "Starting build ksb agent...\n"

echo "Build for WINDOWS..."
GOOS=windows GOARCH=amd64 CXX=x86_64-w64-mingw32-g++ CC=x86_64-w64-mingw32-gcc CGO_ENABLED=1 \
  go build -ldflags "-s -w -X main.version=$VERSION -X main.name=KSBScanner" \
  -o /opt/project/builded/ksbagent/KSBAgent.exe cmd/KSBAgent/*.go
echo "End building for WINDOWS..."

echo "Build for LINUX..."
# собираем сборку для Linux
GOOS=linux GOARCH=amd64 CXX=g++ CC=gcc CGO_ENABLED=1 \
  go build -ldflags "-s -w -X main.version=$VERSION -X main.name=KSBAgent" \
  -o /opt/project/builded/ksbagent/KSBAgent cmd/KSBAgent/*.go
echo "End building for LINUX..."

# shellcheck disable=SC2028
echo "\n"
