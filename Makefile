SHELL = /bin/bash
module = github.com/samherrmann/shopctl
version = $(shell git rev-parse --short HEAD)$(shell [[ -z $$(git status -s) ]] || echo "-dirty")
target = $(shell go env GOOS)-$(shell go env GOARCH)
dist = dist/$(target)

build:
	mkdir -p $(dist) && go build -ldflags "-X $(module)/cmd.Version=$(version)" -o $(dist) .

build.all:
	@export GOOS=linux && export GOARCH=amd64 && make build && make tar
	@export GOOS=windows && export GOARCH=amd64 && make build && make zip

test:
	@go test ./... -cover

clean:
	@rm -rf dist

tar:
	@cd dist && tar -czvf $(target).tar.gz $(target)/*

zip:
	@cd dist && zip -r $(target).zip $(target)/*

# Resources:
# List of available target OSs and architectures:
# https://golang.org/doc/install/source#environment
