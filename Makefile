LDFLAGS :=

build:
	docker compose build --build-arg LDFLAGS="$(LDFLAGS)"

run:
	docker compose up -d

down:
	docker compose down

test:
	go test -race -count 100 -v ./...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v2.1.5

lint: install-lint-deps
	golangci-lint run ./...