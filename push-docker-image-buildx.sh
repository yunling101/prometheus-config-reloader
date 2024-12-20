#!/usr/bin/env bash

set -e
set -o pipefail

CPU_ARCHS="${CPU_ARCHS:-"amd64 arm64 arm/v7"}"
TAG=$(cat version.md | tr -d " \t\n\r")
REPOSITORIES="yunling101/prometheus-config-reloader"

MANIFEST="${REPOSITORIES}:${TAG}"
PLATFORM=

for arch in $CPU_ARCHS; do
    if [[ "x${PLATFORM}" == "x" ]]; then
      PLATFORM="linux/${arch}"
    else
      PLATFORM="${PLATFORM},linux/${arch}"
    fi
done

docker buildx build --platform=${PLATFORM} -t "${MANIFEST}" . --push
