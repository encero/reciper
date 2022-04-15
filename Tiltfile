

docker_build("ghcr.io/encero/reciper-api", ".", dockerfile="docker/go.Dockerfile", build_args={"build_target": "cmd/api/api.go"})
docker_build("ghcr.io/encero/reciper-gql", ".", dockerfile="docker/go.Dockerfile", build_args={"build_target": "gql/server.go"})

k8s_yaml(helm("helm/reciper", name="reciper", values="helm/local-values.yaml"))

k8s_resource("reciper-api", resource_deps=["reciper-nats"], labels=["services"])
k8s_resource("reciper-gql", resource_deps=["reciper-api"], labels=["services"], port_forwards=["8080:8080"])


k8s_resource("reciper-nats", labels=["infrastructure"])


