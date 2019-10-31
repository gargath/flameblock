include .env

COLLECTOR := collector
VISUALIZER := visualizer

VERSION := $(shell git describe --always --dirty --tags 2>/dev/null || echo "undefined")

all: $(COLLECTOR) $(VISUALIZER)

fmt:
	$(GO) fmt ./cmd/... ./pkg/...

vet:
	$(GO) vet ./cmd/... ./pkg/...

lint:
	@ echo "\033[36mLinting code\033[0m"
	$(LINTER) run --disable-all \
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
	$(GO) build -ldflags "-X main.VERSION=${VERSION}" github.com/gargath/flameblock/cmd/collector

$(VISUALIZER): fmt vet lint
	$(GO) build -ldflags "-X main.VERSION=${VERSION}" github.com/gargath/flameblock/cmd/visualizer

clean:
	rm -f $(COLLECTOR) $(VISUALIZER)

.PHONY: all clean fmt vet lint
