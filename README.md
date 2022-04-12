
# Reciper



## GQL

**ID** all ids are UUID v4

## Local setup


Docker compose
```shell
# build and run the stack
docker compose --file docker/docker-compose.yml up --build
```

Kind + helm + ingress-nginx
```shell
# create kind cluster with ingress nginx
./hack/cluster-create.sh

# optionaly: install cert-manger with self signed issuer
./hack/cluster-cert-manager.sh

# install helm chart
# tagged version are listed here https://github.com/encero?tab=packages&repo_name=reciper
helm install reciper --set version=[version] helm/reciper

# or with values file
# copy example and edit the helm/local-values.yaml values file
cp helm/local-values-example.yaml helm/local-values.yaml

# install the chart
helm install reciper --values helm/local-values.yaml helm/reciper
```

## Litestream setup
The litestream setup expect specific k8s secret to be present to pull the AWS credentials from.

```shell
# create k8s secret for litestream containers
kubectl --namespace reciper create secret generic litestream \
    --from-literal=LITESTREAM_ACCESS_KEY_ID="" \
    --from-literal=LITESTREAM_SECRET_ACCESS_KEY=""
```

## Attribution
ios app icon courtesy of <a href="https://www.flaticon.com/free-icons/cooking" title="cooking icons"> by justicon - Flaticon</a>

ios recipe placehoder photo by <a href="https://unsplash.com/@lindsaymoe?utm_source=unsplash&utm_medium=referral&utm_content=creditCopyText">Lindsay Moe</a> on <a href="https://unsplash.com/s/photos/noodles?utm_source=unsplash&utm_medium=referral&utm_content=creditCopyText">Unsplash</a>

