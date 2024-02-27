.PHONY: build clean

VERSION=1.0.0
BIN=wstester
DIR_SRC=.

GO_ENV=CGO_ENABLE=1 GO111MODULE=on
GO=$(GO_ENV) $(shell which go)

build: generate
	@$(GO) build -o $(BIN) $(DIR_SRC)

generate:
	@$(GO) get github.com/google/wire/cmd/wire@v0.5.0
	@$(GO) get github.com/golang/mock/mockgen@v1.6.0
	@$(GO) install github.com/google/wire/cmd/wire@latest
	@$(GO) install github.com/golang/mock/mockgen@v1.6.0
	@$(GO) generate ./...
	@$(GO) mod tidy

clean: 
	@$(GO) clean ./...
	@rm -f $(BIN)

all: clean build