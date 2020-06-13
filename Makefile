SHELL = /bin/bash
Q = @

.DEFAULT_GOAL := test

.PHONY: test
test:
	go test -v -race -count 1 ./...

.PHONY: cover
cover:
	go test -race -coverprofile=cover.out ./...
	go tool cover -html=cover.out -o cover.html
