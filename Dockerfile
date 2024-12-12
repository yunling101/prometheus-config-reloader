ARG GOLANG_BUILDER="1.23-alpine"

FROM golang:${GOLANG_BUILDER} as go_builder

ENV GOPROXY=https://goproxy.cn
ENV GO111MODULE=on
ENV GOPATH=/go

WORKDIR /workspace
COPY . .

RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build go mod download -x && go mod verify
RUN apk add make git
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build make prometheus-config-reloader

FROM busybox:latest

COPY --from=go_builder workspace/prometheus-config-reloader /bin/prometheus-config-reloader
USER nobody
ENTRYPOINT ["/bin/prometheus-config-reloader"]
