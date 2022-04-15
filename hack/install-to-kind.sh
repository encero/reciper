#!/usr/bin/env bash
set -euo pipefail

if [ ! -e helm/local-values-example.yaml ]; then 
    echo "create local-values.yaml:"
    echo 
    echo "cp helm/local-values-example.yaml helm/local-values.yaml"
    exit 1
fi

helm upgrade --install --create-namespace --namespace reciper -f helm/local-values.yaml reciper helm/reciper
