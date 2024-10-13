SHELL := /bin/bash

VERSION := 1.0.0

.PHONY: default
default: fmt build install

.PHONY: fmt
fmt:
	@go fmt ./...
	@goimports -w -local github.com/cmgsj/blob $$(find . -type f -name "*.go" ! -path "./vendor/*")

.PHONY: gen
gen:
	@buf format --write && buf generate

.PHONY: test
test:
	@go test -v ./...

.PHONY: build
build:
	@$(MAKE) binary cmd=build version=$(VERSION)

.PHONY: install
install:
	@$(MAKE) binary cmd=install version=$(VERSION)

.PRONY: binary
binary: gen
	@if [[ -z "$${cmd}" ]]; then \
		echo "must set cmd env var"; \
		exit 1; \
	fi; \
	if [[ "$${cmd}" != "build" && "$${cmd}" != "install" ]]; then \
		echo "unknown cmd '$${cmd}'"; \
		exit 1; \
	fi; \
	if [[ -z "$${version}" ]]; then \
		version="$$(git describe --tags --abbrev=0 2>/dev/null | sed 's/^v//')"; \
	fi; \
	ldflags="-s -w -extldflags='-static'"; \
	if [[ -n "$${version}" ]]; then \
		ldflags+=" -X 'github.com/cmgsj/blob/pkg/cmd/blob.version=$${version}'"; \
	fi; \
	flags=(-trimpath -ldflags="$${ldflags}"); \
	if [[ "$${cmd}" == "build" ]]; then \
		flags+=(-o "bin/blob"); \
	fi; \
	echo "$${cmd}ing blob@$${version} $$(go env GOOS)/$$(go env GOARCH) cgo=$$(go env CGO_ENABLED)"; \
	go "$${cmd}" "$${flags[@]}" .
