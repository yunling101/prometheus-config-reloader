#!/usr/bin/env bash

set -e
set -o pipefail

CPU_ARCHS="${CPU_ARCHS:-"amd64 arm64 arm"}"
TAG=$(cat version.md | tr -d " \t\n\r")
REPOSITORIES="yunling101/prometheus-config-reloader"

for arch in ${CPU_ARCHS}; do
    make --always-make multi-arch GOARCH="$arch" VERSION="${TAG}-$arch"
done

MANIFEST="${REPOSITORIES}:${TAG}"
IMAGES=()
for arch in $CPU_ARCHS; do
    echo "Pushing image ${MANIFEST}-${arch}"
    docker push "${MANIFEST}-${arch}"
		IMAGES[${#IMAGES[@]}]="${MANIFEST}-${arch}"
done

echo "Creating manifest ${MANIFEST}"
docker manifest create --amend "${MANIFEST}" "${IMAGES[@]}"

for arch in $CPU_ARCHS; do
    docker manifest annotate --arch "$arch" "${MANIFEST}" "${REPOSITORIES}:${TAG}-$arch"
done

echo "Pushing manifest ${MANIFEST}"
docker manifest push --purge "${MANIFEST}"
