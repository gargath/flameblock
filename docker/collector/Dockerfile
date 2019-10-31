FROM golang:1.12 as builder

ARG VERSION=undefined

WORKDIR /go/src/github.com/gargath/flameblock

COPY go.mod go.sum ./
RUN GO111MODULE=on go mod download

COPY cmd/ cmd/
COPY pkg/ pkg/

RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o collector -ldflags="-X main.VERSION=${VERSION}" github.com/gargath/flameblock/cmd/collector

FROM alpine:3.9
RUN apk --no-cache add ca-certificates
WORKDIR /bin
COPY --from=builder /go/src/github.com/gargath/flameblock/collector .
ENTRYPOINT ["/bin/collector"]
