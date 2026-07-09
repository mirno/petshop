.PHONY: build

build: bin/petshop

bin/petshop: go.sum
	go build -trimpath -ldflags="-s -w" -o ./$@ ./cmd/petshop
