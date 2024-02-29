.PHONY: build clean

VERSION=1.5.1
GIT_VER=$(shell git rev-parse --short HEAD)
BIN=wstester
DIR_SRC=.

GO_ENV=CGO_ENABLE=1 GO111MODULE=on
GO=$(GO_ENV) $(shell which go)
DOCKER=$(shell which docker)


build:
	@$(GO) build -o $(BIN) $(DIR_SRC)

generate:
	@$(GO) get github.com/google/wire/cmd/wire@v0.5.0
	@$(GO) get github.com/golang/mock/mockgen@v1.6.0
	@$(GO) install github.com/google/wire/cmd/wire@latest
	@$(GO) install github.com/golang/mock/mockgen@v1.6.0
	@$(GO) generate ./...
	@$(GO) mod tidy

docker: build
	@$(DOCKER) build -t hub.hitry.io/test/wstester:v$(GIT_VER) .

clean: 
	@$(GO) clean ./...
	@rm -f $(BIN)

all: clean build