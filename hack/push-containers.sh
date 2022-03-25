#!/usr/bin/env bash
set -euo pipefail

__dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd ${__dir}/..

source ./hack/_functions.sh

echo "Pushing images to gcr.io"

docker push "${DOCKER_BASE}reciper-api:${git_revision}"

for tag in $(git tag --points-at=HEAD); do
    docker push "${DOCKER_BASE}reciper-api:${tag}"
done
