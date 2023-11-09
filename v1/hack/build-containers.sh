#!/usr/bin/env bash
set -euo pipefail

__dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd ${__dir}/..

source ./hack/_functions.sh

echo "Building reciper container"

docker build --tag reciper-api:latest --file docker/go.Dockerfile --build-arg build_target="cmd/api/api.go" .
docker build --tag reciper-gql:latest --file docker/go.Dockerfile --build-arg build_target="gql/server.go" .


