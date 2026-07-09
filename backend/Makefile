.PHONY: build lint generate

build: bin/petshop

bin/petshop: go.sum
	go build -trimpath -ldflags="-s -w" -o ./$@ ./cmd/petshop

lint:
	golangci-lint run ./...

generate:
	go generate ./...