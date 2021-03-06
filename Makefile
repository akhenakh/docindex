.EXPORT_ALL_VARIABLES:

ifndef VERSION
VERSION := $(shell git describe --always --tags)
endif

DATE := $(shell date -u +%Y%m%d.%H%M%S)

LDFLAGS = -trimpath 


targets = query index

.PHONY: all lint test clean query index

all: lint test $(targets)

test: testnolint

testnolint:
	go test -race ./...

lint:
	golangci-lint run ./...

lint-drone:
	golangci-lint run -c .golangci-drone.yml ./...

query:
	go build -o query $(LDFLAGS) github.com/akhenakh/docindex/cmd/query

index:
	go build -o index $(LDFLAGS) github.com/akhenakh/docindex/cmd/index

clean:
	rm -rf cmd/index.db
	rm -f query
	rm -r index
