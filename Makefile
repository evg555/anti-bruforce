APP := "./bin/app"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

GOOSE=goose
MIGRATIONS_DIR=migrations

ifneq (,$(wildcard ./deployments/.env))
    include ./deployments/.env
    export $(shell sed 's/=.*//' ./deployments/.env)
endif

build:
	go build -v -o $(APP) -ldflags "$(LDFLAGS)" ./cmd/main.go

test:
	go test -v -count=1 ./internal/...

test_integrate: up
	go test -v -count=1 ./tests/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.64.5

lint: install-lint-deps
	golangci-lint run ./...

up:
	docker compose -f deployments/docker-compose.yaml up -d

down:
	docker compose -f deployments/docker-compose.yaml down

install-protoc-deps:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

generate: install-protoc-deps
	go generate ./...

migrate-up:
	$(GOOSE) -dir $(MIGRATIONS_DIR) postgres "postgres://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_DATABASE)?sslmode=disable" up

migrate-down:
	$(GOOSE) -dir $(MIGRATIONS_DIR) postgres "postgres://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_DATABASE)?sslmode=disable" down


.PHONY: build test lint up down generate test_integrate migrate-up migrate-down
