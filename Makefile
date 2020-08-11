.PHONY: build test get install docker compose decompose clean

GOPATH=$(shell pwd)/vendor:$(shell pwd)
GOBIN=$(shell pwd)/bin
GONAME=$(shell basename "$(PWD)")
SCALE=1
COMPOSE_FILE_PATH=-f ./docker-compose.yml
TESTFILES=$(dir $(shell find . -name "*_test.go"))

all: build test install

build:
	GOPATH=$(GOPATH) GOBIN=$(GOBIN) go build .

run:
	GOPATH=$(GOPATH) GOBIN=$(GOBIN) go run .

test:
	GOPATH=$(GOPATH) GOBIN=$(GOBIN) go test ./test

get:
	GOPATH=$(GOPATH) GOBIN=$(GOBIN) go mod tidy

install:
	GOPATH=$(GOPATH) GOBIN=$(GOBIN) go install

clean:
	rm -f bin/*
	GOPATH=$(GOPATH) GOBIN=$(GOBIN) go clean -testcache

docker: build
	GONAME=$(GONAME) GOBIN=$(GOBIN) docker build -t $(GONAME) .

compose:
	GONAME=$(GONAME) TAG=$(TAG) docker-compose $(COMPOSE_FILE_PATH) up -d --scale service=$(SCALE)

decompose:
	GONAME=$(GONAME) TAG=$(TAG) docker-compose $(COMPOSE_FILE_PATH) down

