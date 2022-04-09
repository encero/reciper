#!/usr/bin/env bash
set -euo pipefail

__dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd ${__dir}/..

echo "===> Creating cert-manager CRDs"
kubectl apply -f hack/k8s/cert-manager-1.8.0.crds.yaml

echo "===> Setting up helm repo"
helm repo add jetstack https://charts.jetstack.io

echo "===> Create cert-manager namespace"
kubectl create namespace cert-manager

echo "===> Installing cert-manager"
helm install cert-manager jetstack/cert-manager --namespace cert-manager --version v1.8.0 --wait

echo "===> Creating self signed cert issuer"
kubectl apply -f hack/k8s/self-signed-issuer.yaml
