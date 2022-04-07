package dagger

import (
	"dagger.io/dagger"
	"universe.dagger.io/go"
	"universe.dagger.io/docker"
)

dagger.#Plan & {
	client: {
		filesystem: ".": read: {
			contents: dagger.#FS
		}
		filesystem: "./go.*": read: {
			contents: dagger.#FS
		}
	}
	actions: {
		goimage: go.#Image & {
			version: "1.18"
		}

		lint: {
			lintimage: docker.#Build & {
				steps: [
					docker.#Pull & {
						source: "golangci/golangci-lint"
					},
					docker.#Copy & {
						dest:     "/source"
						contents: client.filesystem.".".read.contents
					},
					docker.#Run & {
						command: {
							name: "go"
							args: ["mod", "download"]
						}
						workdir: "/source"
					},
				]
			}

			lint: docker.#Run & {
				command: {
					name: "golangci-lint"
					args: ["run"]
				}
				input:   lintimage.output
				workdir: "/source"
				mounts: {
					"source": {
						dest:     "/source"
						contents: client.filesystem.".".read.contents
					}
				}
			}
		}
		build: go.#Build & {
			source:    client.filesystem.".".read.contents
			package:   "./cmd/api/api.go"
			container: go.#Container & {
				input: goimage.output
			}
		}
	}
}
