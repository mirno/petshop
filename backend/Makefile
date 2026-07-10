.PHONY: build lint generate test

build: bin/petshop

bin/petshop: go.sum
	go build -trimpath -ldflags="-s -w" -o ./$@ ./cmd/petshop

lint:
	golangci-lint run ./...

generate:
	go generate ./...

test:
	go test -covermode atomic -race ./...
