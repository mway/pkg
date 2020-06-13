SHELL = /bin/bash
Q = @
.DEFAULT_GOAL := test
PKG ?= ./...

.PHONY: test
test:
	go test -v -race -count 1 $(PKG)

.PHONY: lint
lint:
	golangci-lint run $(PKG)

.PHONY: cover
cover:
	go test -race -coverprofile=cover.out $(PKG)
	go tool cover -html=cover.out -o cover.html
