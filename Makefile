SHELL = /bin/bash
MAKEFLAGS += --silent
name = merchant
module = github.com/samherrmann/$(name)
version = $(shell git rev-parse --short HEAD)$(shell [[ -z $$(git status -s) ]] || echo "-dirty")
target = $(shell go env GOOS)-$(shell go env GOARCH)
dist = dist/$(target)

build:
	mkdir -p $(dist) && \
	go build \
		-race \
		-ldflags "-s -w -X $(module)/config.Version=$(version) -X $(module)/config.AppName=$(name)" \
		-o $(dist) .

build.all:
	export GOOS=linux && export GOARCH=amd64 && make build && make tar
	export GOOS=windows && export GOARCH=amd64 && make build && make zip

test:
	go test ./... -race -cover

clean:
	rm -rf dist

tar:
	cd $(dist) && tar -czvf ../$(name)-$(target).tar.gz *

zip:
	cd $(dist) && zip -r ../$(name)-$(target).zip *

# Resources:
# List of available target OSs and architectures:
# https://golang.org/doc/install/source#environment
