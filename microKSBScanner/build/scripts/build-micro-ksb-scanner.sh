#!/bin/sh

set -e

# shellcheck disable=SC2028
echo "Starting build microKSBScanner...\n"

echo "Build for WINDOWS..."
GOOS=windows GOARCH=amd64 CXX=x86_64-w64-mingw32-g++ CC=x86_64-w64-mingw32-gcc CGO_ENABLED=1 \
  go build -ldflags "-s -w -X main.version=$VERSION -X main.name=microKSBScanner" \
  -o /opt/project/builded/microksbscanner/microKSBScanner.exe cmd/microKSBScanner/*.go
echo "End building for WINDOWS..."

echo "Build for LINUX..."
# собираем сборку для Linux
GOOS=linux GOARCH=amd64 CXX=g++ CC=gcc CGO_ENABLED=1 \
  go build -ldflags "-s -w -X main.version=$VERSION -X main.name=microKSBScanner" \
  -o /opt/project/builded/microksbscanner/microKSBScanner cmd/microKSBScanner/*.go
echo "End buildingBuild for LINUX..."

# shellcheck disable=SC2028
echo "\n"
