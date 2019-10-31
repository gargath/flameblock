include .env

COLLECTOR := collector
VISUALIZER := visualizer

VERSION := $(shell git describe --always --dirty --tags 2>/dev/null || echo "undefined")

IMG_PREFIX := flameblock

all: $(COLLECTOR) $(VISUALIZER)

fmt:
	GO111MODULE=on $(GO) fmt ./cmd/... ./pkg/...
	GO111MODULE=on $(GOIMPORTS) -w ./cmd ./pkg
vet:
	GO111MODULE=on $(GO) vet ./cmd/... ./pkg/...

lint:
	@ echo "\033[36mLinting code\033[0m"
	GO111MODULE=on $(LINTER) run --disable-all \
                --exclude-use-default=false \
                --enable=govet \
                --enable=ineffassign \
                --enable=deadcode \
                --enable=golint \
                --enable=goconst \
                --enable=gofmt \
                --enable=goimports \
                --skip-dirs=pkg/client/ \
                --deadline=120s \
                --tests ./...
	@ echo


$(COLLECTOR): fmt vet lint
	GO111MODULE=on $(GO) build -ldflags "-X main.VERSION=${VERSION}" github.com/gargath/flameblock/cmd/collector

$(VISUALIZER): fmt vet lint
	GO111MODULE=on $(GO) build -ldflags "-X main.VERSION=${VERSION}" github.com/gargath/flameblock/cmd/visualizer

docker-build: fmt vet lint
	docker build --build-arg VERSION=${VERSION} -t ${IMG_PREFIX}-collector:${VERSION} -f deploy/docker/collector/Dockerfile .
	@echo "\033[36mBuilt ${IMG_PREFIX}-collector:${VERSION}\033[0m"
	docker build --build-arg VERSION=${VERSION} -t ${IMG_PREFIX}-visualizer:${VERSION} -f deploy/docker/visualizer/Dockerfile .
	@echo "\033[36mBuilt ${IMG_PREFIX}-visualizer:${VERSION}\033[0m"

clean:
	rm -f $(COLLECTOR) $(VISUALIZER)

distclean: clean
	rm -f ./env

.PHONY: all clean distclean fmt vet lint docker-build
