FROM golang:1.24 AS builder

WORKDIR /app

COPY . .

ARG LDFLAGS

RUN export GOPATH=$PWD/.cache && \
  export PATH=$GOPATH/bin:$PATH && \
  mkdir -p $GOPATH && \
  go generate ./... && \
  go mod tidy && \
  go mod download -x && \
  CGO_ENABLED=0 GOOS=linux go build -o server -ldflags "$LDFLAGS" ./cmd/server/main.go

FROM alpine:3.21

COPY --from=builder /app/server /app/server
COPY --from=builder /app/.env /app/.env
COPY --from=builder /app/tmp /app/tmp

WORKDIR /app

CMD ["./server"]