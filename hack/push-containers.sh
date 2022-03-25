#!/usr/bin/env bash
set -euo pipefail

__dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd ${__dir}/..

source ./hack/_functions.sh

echo "Pushing images to ghcr.io"

image="${DOCKER_BASE}reciper-api:${git_revision}"
echo ">> pushing ${image}"
docker push "${image}" 

for tag in $(git tag --points-at=HEAD); do
    image="${DOCKER_BASE}reciper-api:${tag}"
    echo ">> pushing ${image}"

    docker push "${image}"
done
