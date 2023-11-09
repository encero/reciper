DOCKER_BASE="ghcr.io/encero/"

git_revision=$(git rev-parse --short HEAD)
git_tags=$(git tag --points-at=HEAD)


function tags_for_docker () {

    local tags=("--tag ${DOCKER_BASE}reciper-api:${git_revision}")

    for tag in $git_tags; do
        tags+=("--tag ${DOCKER_BASE}reciper-api:${tag}")
    done

    echo "${tags[@]}"
}

