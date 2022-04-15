FROM  golang:1.18-alpine AS build

ENV CGO_ENABLED=0
ARG build_target
ARG version

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build \
go build -o /out/app -ldflags="-X 'main.Version=${version}'" ${build_target}

FROM scratch AS bin-unix

LABEL org.opencontainers.image.source https://github.com/encero/reciper

COPY --from=build /etc/ssl/certs /etc/ssl/certs
COPY --from=build /out/app /app
CMD ["/app"]



