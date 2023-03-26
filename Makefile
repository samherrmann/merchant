SHELL = /bin/bash
MAKEFLAGS += --silent
name = merchant
module = github.com/samherrmann/$(name)
version = $(shell git rev-parse --short HEAD)$(shell [[ -z $$(git status -s) ]] || echo "-dirty")
target = $(shell go env GOOS)-$(shell go env GOARCH)
dist = dist/$(target)

.PHONY: build
build:
	mkdir -p $(dist) && \
	go build \
		-ldflags "-s -w -X $(module)/config.Version=$(version) -X $(module)/config.AppName=$(name)" \
		-o $(dist) .

.PHONY: build.all
build.all:
	export GOOS=linux && export GOARCH=amd64 && make build && make tar
	export GOOS=windows && export GOARCH=amd64 && make build && make zip

.PHONY: test
test:
	go test ./... -race -cover

.PHONY: clean
clean:
	rm -rf dist

.PHONY: tar
tar:
	cd $(dist) && tar -czvf ../$(name)-$(target).tar.gz *

.PHONY: zip
zip:
	cd $(dist) && zip -r ../$(name)-$(target).zip *

.PHONY: lint
lint:
	staticcheck -checks=all ./...

# Resources:
# List of available target OSs and architectures:
# https://golang.org/doc/install/source#environment
