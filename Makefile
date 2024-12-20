SHELL=/usr/bin/env sh -o pipefail

GOOS?=$(shell go env GOOS)
GOARCH?=$(shell go env GOARCH)

ifeq ($(GOARCH),arm)
	ARCH=armv7
else
	ARCH=$(GOARCH)
endif

TAG?=$(shell git rev-parse --short HEAD)
VERSION?=$(shell cat version.md | tr -d " \t\n\r")
BUILD_DATE=$(shell date +"%Y%m%d-%T")

PROJECT_PKG=github.com/yunling101/prometheus-config-reloader

GO_BUILD_LDFLAGS=\
	-s \
	-X $(PROJECT_PKG)/version.BuildDate=$(BUILD_DATE) \
	-X $(PROJECT_PKG)/version.Version=$(VERSION) \
	-X $(PROJECT_PKG)/version.Revision=$(TAG)

GO_BUILD_RECIPE=\
	GOOS=$(GOOS) \
	GOARCH=$(GOARCH) \
	CGO_ENABLED=0 \
	go build -ldflags="$(GO_BUILD_LDFLAGS)"

.PHONY: prometheus-config-reloader
prometheus-config-reloader:
	$(GO_BUILD_RECIPE) -o $@ main.go

.PHONY: prometheus-config-reloader-image
prometheus-config-reloader-image:
	@docker build \
		--build-arg GOARCH=$(GOARCH) \
		--build-arg OS=$(GOOS) \
		-t yunling101/prometheus-config-reloader:$(VERSION) -f Dockerfile .

.PHONY: multi-arch
multi-arch: GOOS := linux
multi-arch:
	@docker build \
		--build-arg ARCH=$(ARCH) \
		--build-arg GOARCH=$(GOARCH) \
		--build-arg OS=$(GOOS) \
		-t yunling101/prometheus-config-reloader:$(VERSION) -f Dockerfile.multi-arch .
