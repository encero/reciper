#!/usr/bin/env bash
set -euo pipefail

__dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd ${__dir}/..

source ./hack/_functions.sh

echo "Building reciper-api"

docker build $(tags_for_docker) --file docker/api.Dockerfile --build-arg build_target="cmd/api/api.go" .


